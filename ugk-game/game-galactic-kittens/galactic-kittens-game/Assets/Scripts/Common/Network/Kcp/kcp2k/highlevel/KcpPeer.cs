// Kcp Peer, similar to UDP Peer but wrapped with reliability, channels,
// timeouts, authentication, state, etc.
//
// still IO agnostic to work with udp, nonalloc, relays, native, etc.

using System;
using System.Diagnostics;
using System.Net.Sockets;

namespace kcp2k
{
    enum KcpState
    {
        Connected,
        Authenticated,
        Disconnected
    }

    /// <summary>
    /// <para>连接会话</para>
    /// <para>协议格式：消息长度4+消息id4+序列号4+时间戳8+protobuf消息体</para>
    /// Mirror 有可靠、不可靠消息，验证、心跳等逻辑，这里自定义简化
    /// </summary>
    public class KcpPeer
    {
        // kcp reliability algorithm
        internal Kcp kcp;


        // IO agnostic
        readonly Action<ArraySegment<byte>> RawSend;

        // state: connected as soon as we create the peer.
        // leftover from KcpConnection. remove it after refactoring later.
        KcpState state = KcpState.Connected;

        readonly Action OnAuthenticated;
        readonly Action<ArraySegment<byte>> OnData;
        //发送心跳
        private readonly Action SendPing; 

        readonly Action OnDisconnected;

        // error callback instead of logging.
        // allows libraries to show popups etc.
        // (string instead of Exception for ease of use and to avoid user panic)
        readonly Action<ErrorCode, string> OnError;

        // If we don't receive anything these many milliseconds
        // then consider us disconnected
        public const int DEFAULT_TIMEOUT = 10000;
        public int timeout;
        //最后收到消息时间，用于判断超时
        uint lastReceiveTime;

        // internal time.
        // StopWatch offers ElapsedMilliSeconds and should be more precise than
        // Unity's time.deltaTime over long periods.
        readonly Stopwatch watch = new Stopwatch();


        // reliable channel (= kcp) MaxMessageSize so the outside knows largest
        // allowed message to send. the calculation in Send() is not obvious at
        // all, so let's provide the helper here.
        //
        // kcp does fragmentation, so max message is way larger than MTU.
        //
        // -> runtime MTU changes are disabled: mss is always MTU_DEF-OVERHEAD
        // -> Send() checks if fragment count < rcv_wnd, so we use rcv_wnd - 1.
        //    NOTE that original kcp has a bug where WND_RCV default is used
        //    instead of configured rcv_wnd, limiting max message size to 144 KB
        //    https://github.com/skywind3000/kcp/pull/291
        //    we fixed this in kcp2k.
        // -> we add 1 byte KcpHeader enum to each message, so -1
        //
        // IMPORTANT: max message is MTU * rcv_wnd, in other words it completely
        //            fills the receive window! due to head of line blocking,
        //            all other messages have to wait while a maxed size message
        //            is being delivered.
        //            => in other words, DO NOT use max size all the time like
        //               for batching.
        //            => sending UNRELIABLE max message size most of the time is
        //               best for performance (use that one for batching!)
        static int ReliableMaxMessageSize_Unconstrained(int mtu, uint rcv_wnd) =>
            (mtu - Kcp.OVERHEAD) * ((int)rcv_wnd - 1) - 1;

        // kcp encodes 'frg' as 1 byte.
        // max message size can only ever allow up to 255 fragments.
        //   WND_RCV gives 127 fragments.
        //   WND_RCV * 2 gives 255 fragments.
        // so we can limit max message size by limiting rcv_wnd parameter.
        public static int ReliableMaxMessageSize(int mtu, uint rcv_wnd) =>
            ReliableMaxMessageSize_Unconstrained(mtu, Math.Min(rcv_wnd, Kcp.FRG_MAX));

        // buffer to receive kcp's processed messages (avoids allocations).
        // IMPORTANT: this is for KCP messages. so it needs to be of size:
        //            1 byte header + MaxMessageSize content
        readonly byte[] kcpMessageBuffer; // = new byte[1 + ReliableMaxMessageSize];

        // send buffer for handing user messages to kcp for processing.
        // (avoids allocations).
        // IMPORTANT: needs to be of size:
        //            1 byte header + MaxMessageSize content
        readonly byte[] kcpSendBuffer; // = new byte[1 + ReliableMaxMessageSize];

        // raw send buffer is exactly MTU.
        readonly byte[] rawSendBuffer;

        // send a ping occasionally so we don't time out on the other end.
        // for example, creating a character in an MMO could easily take a
        // minute of no data being sent. which doesn't mean we want to time out.
        // same goes for slow paced card games etc.
        public const int PING_INTERVAL = 2000;
        uint lastPingTime;

        // if we send more than kcp can handle, we will get ever growing
        // send/recv buffers and queues and minutes of latency.
        // => if a connection can't keep up, it should be disconnected instead
        //    to protect the server under heavy load, and because there is no
        //    point in growing to gigabytes of memory or minutes of latency!
        // => 2k isn't enough. we reach 2k when spawning 4k monsters at once
        //    easily, but it does recover over time.
        // => 10k seems safe.
        //
        // note: we have a ChokeConnectionAutoDisconnects test for this too!
        internal const int QueueDisconnectThreshold = 10000;

        // getters for queue and buffer counts, used for debug info
        public int SendQueueCount => kcp.snd_queue.Count;
        public int ReceiveQueueCount => kcp.rcv_queue.Count;
        public int SendBufferCount => kcp.snd_buf.Count;
        public int ReceiveBufferCount => kcp.rcv_buf.Count;

        // maximum send rate per second can be calculated from kcp parameters
        // source: https://translate.google.com/translate?sl=auto&tl=en&u=https://wetest.qq.com/lab/view/391.html
        //
        // KCP can send/receive a maximum of WND*MTU per interval.
        // multiple by 1000ms / interval to get the per-second rate.
        //
        // example:
        //   WND(32) * MTU(1400) = 43.75KB
        //   => 43.75KB * 1000 / INTERVAL(10) = 4375KB/s
        //
        // returns bytes/second! @
        public uint MaxSendRate => kcp.snd_wnd * kcp.mtu * 1000 / kcp.interval;
        public uint MaxReceiveRate => kcp.rcv_wnd * kcp.mtu * 1000 / kcp.interval;

        // calculate max message sizes based on mtu and wnd only once
        public readonly int reliableMax;

        // SetupKcp creates and configures a new KCP instance.
        // => useful to start from a fresh state every time the client connects
        // => NoDelay, interval, wnd size are the most important configurations.
        //    let's force require the parameters so we don't forget it anywhere.
        public KcpPeer(
            Action<ArraySegment<byte>> output,
            Action OnAuthenticated,
            Action<ArraySegment<byte>> OnData,
            Action OnDisconnected,
            Action<ErrorCode, string> OnError,
            Action sendPing,
            KcpConfig config,
            uint cookie)
        {
            // initialize callbacks first to ensure they can be used safely.
            this.OnAuthenticated = OnAuthenticated;
            this.OnData = OnData;
            this.OnDisconnected = OnDisconnected;
            this.OnError = OnError;
            this.RawSend = output;
            this.SendPing = sendPing;

            // set up kcp over reliable channel (that's what kcp is for)
            kcp = new Kcp(0, RawSendReliable);

            // set nodelay.
            // note that kcp uses 'nocwnd' internally so we negate the parameter
            kcp.SetNoDelay(config.NoDelay ? 1u : 0u, config.Interval, config.FastResend, !config.CongestionWindow);
            kcp.SetWindowSize(config.SendWindowSize, config.ReceiveWindowSize);

            // IMPORTANT: high level needs to add 1 channel byte to each raw
            // message. so while Kcp.MTU_DEF is perfect, we actually need to
            // tell kcp to use MTU-1 so we can still put the header into the
            // message afterwards.
            kcp.SetMtu((uint)config.Mtu);

            // create mtu sized send buffer
            rawSendBuffer = new byte[config.Mtu];

            reliableMax = ReliableMaxMessageSize(config.Mtu, config.ReceiveWindowSize);

            // set maximum retransmits (aka dead_link)
            kcp.dead_link = config.MaxRetransmits;

            // create message buffers AFTER window size is set
            // see comments on buffer definition for the "+1" part
            kcpMessageBuffer = new byte[1 + reliableMax];
            kcpSendBuffer = new byte[1 + reliableMax];

            timeout = config.Timeout;

            watch.Start();
        }

        /// <summary>
        /// 超时处理 
        /// </summary>
        /// <param name="time"></param>
        void HandleTimeout(uint time)
        {
            // note: we are also sending a ping regularly, so timeout should
            //       only ever happen if the connection is truly gone.
            if (time >= lastReceiveTime + timeout)
            {
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.Timeout,
                    $"KcpPeer: Connection timed out after not receiving any message for {timeout}ms. Disconnecting.");
                Disconnect();
            }
        }

        void HandleDeadLink()
        {
            // kcp has 'dead_link' detection. might as well use it.
            if (kcp.state == -1)
            {
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.Timeout,
                    $"KcpPeer: dead_link detected: a message was retransmitted {kcp.dead_link} times without ack. Disconnecting.");
                Disconnect();
            }
        }

        // send a ping occasionally in order to not time out on the other end. @
        void HandlePing(uint time)
        {
            
            // enough time elapsed since last ping?
            if (time >= lastPingTime + PING_INTERVAL)
            {
                // 发送心跳，使用回调函数
                SendPing();
                lastPingTime = time;
            }
        }

        /// <summary>
        /// 阻塞检测
        /// </summary>
        void HandleChoked()
        {
            // disconnect connections that can't process the load.
            // see QueueSizeDisconnect comments.
            // => include all of kcp's buffers and the unreliable queue!
            int total = kcp.rcv_queue.Count + kcp.snd_queue.Count +
                        kcp.rcv_buf.Count + kcp.snd_buf.Count;
            if (total >= QueueDisconnectThreshold)
            {
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.Congestion,
                    $"KcpPeer: disconnecting connection because it can't process data fast enough.\n" +
                    $"Queue total {total}>{QueueDisconnectThreshold}. rcv_queue={kcp.rcv_queue.Count} snd_queue={kcp.snd_queue.Count} rcv_buf={kcp.rcv_buf.Count} snd_buf={kcp.snd_buf.Count}\n" +
                    $"* Try to Enable NoDelay, decrease INTERVAL, disable Congestion Window (= enable NOCWND!), increase SEND/RECV WINDOW or compress data.\n" +
                    $"* Or perhaps the network is simply too slow on our end, or on the other end.");

                // let's clear all pending sends before disconnting with 'Bye'.
                // otherwise a single Flush in Disconnect() won't be enough to
                // flush thousands of messages to finally deliver 'Bye'.
                // this is just faster and more robust.
                kcp.snd_queue.Clear();

                Disconnect();
            }
        }

        // reads the next reliable message type & content from kcp.
        // -> to avoid buffering, unreliable messages call OnData directly. 
        bool ReceiveNextReliable(out ArraySegment<byte> message)
        {
            message = default;

            int msgSize = kcp.PeekSize();
            if (msgSize <= 0) return false;

            // only allow receiving up to buffer sized messages.
            // otherwise we would get BlockCopy ArgumentException anyway.
            if (msgSize > kcpMessageBuffer.Length)
            {
                // we don't allow sending messages > Max, so this must be an
                // attacker. let's disconnect to avoid allocation attacks etc.
                // pass error to user callback. no need to log it manually.
                OnError(ErrorCode.InvalidReceive,
                    $"KcpPeer: possible allocation attack for msgSize {msgSize} > buffer {kcpMessageBuffer.Length}. Disconnecting the connection.");
                Disconnect();
                return false;
            }

            // receive from kcp
            int received = kcp.Receive(kcpMessageBuffer, msgSize);
            if (received < 0)
            {
                // if receive failed, close everything
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.InvalidReceive,
                    $"KcpPeer: Receive failed with error={received}. closing connection.");
                Disconnect();
                return false;
            }

            // extract header & content without header
            // message = new ArraySegment<byte>(kcpMessageBuffer, 1, msgSize - 1);
            message = new ArraySegment<byte>(kcpMessageBuffer, 0, msgSize );
            lastReceiveTime = (uint)watch.ElapsedMilliseconds;
            return true;
        }


        public void TickIncoming()
        {
            uint time = (uint)watch.ElapsedMilliseconds;

            try
            {
                // detect common events & ping
                HandleTimeout(time);
                HandleDeadLink();
                HandlePing(time);
                HandleChoked();

                // process all received messages
                while (ReceiveNextReliable(out ArraySegment<byte> message))
                {
                    // call OnData IF the message contained actual data
                    if (message.Count > 0)
                    {
                        //Log.Warning($"Kcp recv msg: {BitConverter.ToString(message.Array, message.Offset, message.Count)}");
                        OnData?.Invoke(message);
                    }
                    // empty data = attacker, or something went wrong
                    else
                    {
                        // pass error to user callback. no need to log it manually.
                        // GetType() shows Server/ClientConn instead of just Connection.
                        OnError(ErrorCode.InvalidReceive,
                            $"KcpPeer: received empty Data message while Authenticated. Disconnecting the connection.");
                        Disconnect();
                    }
                }
            }
            // TODO KcpConnection is IO agnostic. move this to outside later.
            catch (SocketException exception)
            {
                // this is ok, the connection was closed
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.ConnectionClosed, $"KcpPeer: Disconnecting because {exception}. This is fine.");
                Disconnect();
            }
            catch (ObjectDisposedException exception)
            {
                // fine, socket was closed
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.ConnectionClosed, $"KcpPeer: Disconnecting because {exception}. This is fine.");
                Disconnect();
            }
            catch (Exception exception)
            {
                // unexpected
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.Unexpected, $"KcpPeer: unexpected Exception: {exception}");
                Disconnect();
            }
        }

        public void TickOutgoing()
        {
            uint time = (uint)watch.ElapsedMilliseconds;

            try
            {
                switch (state)
                {
                    case KcpState.Connected:
                    case KcpState.Authenticated:
                    {
                        // update flushes out messages
                        kcp.Update(time);
                        break;
                    }
                    case KcpState.Disconnected:
                    {
                        // do nothing while disconnected
                        break;
                    }
                }
            }
            // TODO KcpConnection is IO agnostic. move this to outside later.
            catch (SocketException exception)
            {
                // this is ok, the connection was closed
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.ConnectionClosed, $"KcpPeer: Disconnecting because {exception}. This is fine.");
                Disconnect();
            }
            catch (ObjectDisposedException exception)
            {
                // fine, socket was closed
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.ConnectionClosed, $"KcpPeer: Disconnecting because {exception}. This is fine.");
                Disconnect();
            }
            catch (Exception exception)
            {
                // unexpected
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.Unexpected, $"KcpPeer: unexpected exception: {exception}");
                Disconnect();
            }
        }

        // 
        void OnRawInputReliable(ArraySegment<byte> message)
        {
            // input into kcp, but skip channel byte
            int input = kcp.Input(message.Array, message.Offset, message.Count);
            if (input != 0)
            {
                // GetType() shows Server/ClientConn instead of just Connection.
                Log.Warning($"KcpPeer: Input failed with error={input} for buffer with length={message.Count - 1}");
            }
        }
        

        // insert raw IO. usually from socket.Receive.
        // offset is useful for relays, where we may parse a header and then
        // feed the rest to kcp. 
        public void RawInput(ArraySegment<byte> segment)
        {
            // ensure valid size 消息长度4+消息id4+序列号4+时间戳8+protobuf消息体
            if (segment.Count <= 20) return;
            
            // parse message
            ArraySegment<byte> message =
                new ArraySegment<byte>(segment.Array, segment.Offset , segment.Count);
            OnRawInputReliable(message);
        }

        // raw send called by kcp 
        void RawSendReliable(byte[] data, int length)
        {
            Buffer.BlockCopy(data, 0, rawSendBuffer, 0, length);
            // IO send
            ArraySegment<byte> segment = new ArraySegment<byte>(rawSendBuffer, 0, length);
            RawSend(segment);
        }

        void SendReliable(ArraySegment<byte> content)
        {
            // 1 byte header + content needs to fit into send buffer
            if (content.Count > kcpSendBuffer.Length) 
            {
                // otherwise content is larger than MaxMessageSize. let user know!
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.InvalidSend,
                    $"KcpPeer: Failed to send reliable message of size {content.Count} because it's larger than ReliableMaxMessageSize={reliableMax}");
                return;
            }


            // write data (if any)
            if (content.Count > 0)
                Buffer.BlockCopy(content.Array, content.Offset, kcpSendBuffer, 0, content.Count);

            // send to kcp for processing
            int sent = kcp.Send(kcpSendBuffer, 0,  content.Count);
            if (sent < 0)
            {
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.InvalidSend,
                    $"KcpPeer: Send failed with error={sent} for content with length={content.Count}");
            }
        }

      

      

        public void SendData(ArraySegment<byte> data)
        {
            // sending empty segments is not allowed.
            // nobody should ever try to send empty data.
            // it means that something went wrong, e.g. in Mirror/DOTSNET.
            // let's make it obvious so it's easy to debug.
            if (data.Count == 0)
            {
                // pass error to user callback. no need to log it manually.
                // GetType() shows Server/ClientConn instead of just Connection.
                OnError(ErrorCode.InvalidSend,
                    $"KcpPeer: tried sending empty message. This should never happen. Disconnecting.");
                Disconnect();
                return;
            }

            SendReliable( data);
        }



        // disconnect this connection
        public void Disconnect()
        {
            // only if not disconnected yet
            if (state == KcpState.Disconnected)
                return;

            // send a disconnect message
            try
            {
                kcp.Flush();
            }
            // TODO KcpConnection is IO agnostic. move this to outside later.
            catch (SocketException)
            {
                // this is ok, the connection was already closed
            }
            catch (ObjectDisposedException)
            {
                // this is normal when we stop the server
                // the socket is stopped so we can't send anything anymore
                // to the clients

                // the clients will eventually timeout and realize they
                // were disconnected
            }

            // set as Disconnected, call event
            // GetType() shows Server/ClientConn instead of just Connection.
            Log.Info($"KcpPeer: Disconnected.");
            state = KcpState.Disconnected;
            OnDisconnected?.Invoke();
        }
    }
}