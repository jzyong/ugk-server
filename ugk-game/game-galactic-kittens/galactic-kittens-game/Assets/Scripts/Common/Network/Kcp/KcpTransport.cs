//#if MIRROR <- commented out because MIRROR isn't defined on first import yet
using System;
using System.Net;
using UnityEngine;
using Unity.Collections;
using UnityEngine.Serialization;
using Common.Network;

namespace kcp2k
{
    /// <summary>
    /// kcp 传输层逻辑
    /// </summary>
    public class KcpTransport : Transport
    {
        // scheme used by this transport
        public const string Scheme = "kcp";

        // common
        [Header("Transport Configuration")]
        [Tooltip("NoDelay is recommended to reduce latency. This also scales better without buffers getting full.")]
        public bool NoDelay = true;
        [Tooltip("KCP internal update interval. 100ms is KCP default, but a lower interval is recommended to minimize latency and to scale to more networked entities.")]
        public uint Interval = 10;
        [Tooltip("KCP timeout in milliseconds. Note that KCP sends a ping automatically.")]
        public int Timeout = 10000;
        [Tooltip("Socket receive buffer size. Large buffer helps support more connections. Increase operating system socket buffer size limits if needed.")]
        public int RecvBufferSize = 1024 * 1027 * 7;
        [Tooltip("Socket send buffer size. Large buffer helps support more connections. Increase operating system socket buffer size limits if needed.")]
        public int SendBufferSize = 1024 * 1027 * 7;

        [Header("Advanced")]
        [Tooltip("KCP fastresend parameter. Faster resend for the cost of higher bandwidth. 0 in normal mode, 2 in turbo mode.")]
        public int FastResend = 2;
        [Tooltip("KCP congestion window. Restricts window size to reduce congestion. Results in only 2-3 MTU messages per Flush even on loopback. Best to keept his disabled.")]
        /*public*/ bool CongestionWindow = false; // KCP 'NoCongestionWindow' is false by default. here we negate it for ease of use.
        [Tooltip("KCP window size can be modified to support higher loads. This also increases max message size.")]
        public uint ReceiveWindowSize = 4096; //Kcp.WND_RCV; 128 by default. Mirror sends a lot, so we need a lot more.
        [Tooltip("KCP window size can be modified to support higher loads.")]
        public uint SendWindowSize = 4096; //Kcp.WND_SND; 32 by default. Mirror sends a lot, so we need a lot more.
        
        [Tooltip("KCP will try to retransmit lost messages up to MaxRetransmit (aka dead_link) before disconnecting.")]
        public uint MaxRetransmit = Kcp.DEADLINK * 2; // default prematurely disconnects a lot of people (#3022). use 2x.
        [Tooltip("Enable to automatically set client & server send/recv buffers to OS limit. Avoids issues with too small buffers under heavy load, potentially dropping connections. Increase the OS limit if this is still too small.")]
        [FormerlySerializedAs("MaximizeSendReceiveBuffersToOSLimit")]
        public bool MaximizeSocketBuffers = true;

        [Header("Allowed Max Message Sizes\nBased on Receive Window Size")]
        [Tooltip("KCP reliable max message size shown for convenience. Can be changed via ReceiveWindowSize.")]
        [ReadOnly] public int ReliableMaxMessageSize = 0; // readonly, displayed from OnValidate
        
      

        // config is created from the serialized properties above.
        // we can expose the config directly in the future.
        // for now, let's not break people's old settings. @
        protected KcpConfig config;

        // use default MTU for this transport.
        const int MTU = Kcp.MTU_DEF;

        //  client
        protected KcpClient client;

        // debugging
        [Header("Debug")]
        public bool debugLog;
        // show statistics in OnGUI
        public bool statisticsGUI;
        // log statistics for headless servers that can't show them in GUI
        public bool statisticsLog;

      

        /// <summary>
        /// Kcp异常转为传输层异常
        /// </summary>
        /// <param name="error"></param>
        /// <returns></returns>
        /// <exception cref="InvalidCastException"></exception>
        public static TransportError ToTransportError(ErrorCode error)
        {
            switch(error)
            {
                case ErrorCode.DnsResolve: return TransportError.DnsResolve;
                case ErrorCode.Timeout: return TransportError.Timeout;
                case ErrorCode.Congestion: return TransportError.Congestion;
                case ErrorCode.InvalidReceive: return TransportError.InvalidReceive;
                case ErrorCode.InvalidSend: return TransportError.InvalidSend;
                case ErrorCode.ConnectionClosed: return TransportError.ConnectionClosed;
                case ErrorCode.Unexpected: return TransportError.Unexpected;
                default: throw new InvalidCastException($"KCP: missing error translation for {error}");
            }
        }

        protected virtual void Awake()
        {
            // logging
            //   Log.Info should use Debug.Log if enabled, or nothing otherwise
            //   (don't want to spam the console on headless servers)
            if (debugLog)
                Log.Info = Debug.Log;
            else
                Log.Info = _ => {};
            Log.Warning = Debug.LogWarning;
            Log.Error = Debug.LogError;

            // create config from serialized settings
            config = new KcpConfig( RecvBufferSize, SendBufferSize, MTU, NoDelay, Interval, FastResend, CongestionWindow, SendWindowSize, ReceiveWindowSize, Timeout, MaxRetransmit);

            // client (NonAlloc version is not necessary anymore)
            client = new KcpClient(
                () => OnClientConnected.Invoke(),
                (message) => OnClientDataReceived.Invoke(message),
                () => OnClientDisconnected.Invoke(),
                (error, reason) => OnClientError.Invoke(ToTransportError(error), reason),
                () =>SendHeart.Invoke(),
                config
            );

          

            if (statisticsLog)
                InvokeRepeating(nameof(OnLogStatistics), 1, 1);

            Log.Info("KcpTransport initialized!");
        }

        protected virtual void OnValidate()
        {
            // show max message sizes in inspector for convenience.
            // 'config' isn't available in edit mode yet, so use MTU define.
            ReliableMaxMessageSize = KcpPeer.ReliableMaxMessageSize(MTU, ReceiveWindowSize);
        }

        // all except WebGL
        public override bool Available() =>
            Application.platform != RuntimePlatform.WebGLPlayer;

        // client
        public override bool ClientConnected() => client.connected;
        public override void ClientConnect(string address,ushort port)
        {
            client.Connect(address, port);
        }
        public override void ClientSend(ArraySegment<byte> segment)
        {
            client.Send(segment);

            // call event. might be null if no statistics are listening etc.
            OnClientDataSent?.Invoke(segment);
        }
        public override void ClientDisconnect() => client.Disconnect();
        // process incoming in early update
        public override void ClientEarlyUpdate()
        {
           
             client.TickIncoming();
        }
        // process outgoing in late update
        public override void ClientLateUpdate() => client.TickOutgoing();

        

        // common
        public override void Shutdown() {}

        // max message size
        public override int GetMaxPacketSize()
        {
            return KcpPeer.ReliableMaxMessageSize(config.Mtu, ReceiveWindowSize);
            
        }


        // PrettyBytes function from DOTSNET
        // pretty prints bytes as KB/MB/GB/etc.
        // long to support > 2GB
        // divides by floats to return "2.5MB" etc. 
        public static string PrettyBytes(long bytes)
        {
            // bytes
            if (bytes < 1024)
                return $"{bytes} B";
            // kilobytes
            else if (bytes < 1024L * 1024L)
                return $"{(bytes / 1024f):F2} KB";
            // megabytes
            else if (bytes < 1024 * 1024L * 1024L)
                return $"{(bytes / (1024f * 1024f)):F2} MB";
            // gigabytes
            return $"{(bytes / (1024f * 1024f * 1024f)):F2} GB";
        }

        /// <summary>
        /// 显示统计数据
        /// </summary>
        protected virtual void OnGUIStatistics()
        {
            GUILayout.BeginArea(new Rect(5, 110, 300, 300));
            if (ClientConnected())
            {
                GUILayout.BeginVertical("Box");
                GUILayout.Label("CLIENT");
                GUILayout.Label($"  MaxSendRate: {PrettyBytes(client.peer.MaxSendRate)}/s");
                GUILayout.Label($"  MaxRecvRate: {PrettyBytes(client.peer.MaxReceiveRate)}/s");
                GUILayout.Label($"  SendQueue: {client.peer.SendQueueCount}");
                GUILayout.Label($"  ReceiveQueue: {client.peer.ReceiveQueueCount}");
                GUILayout.Label($"  SendBuffer: {client.peer.SendBufferCount}");
                GUILayout.Label($"  ReceiveBuffer: {client.peer.ReceiveBufferCount}");
                GUILayout.EndVertical();
            }

            GUILayout.EndArea();
        }

// OnGUI allocates even if it does nothing. avoid in release.
#if UNITY_EDITOR || DEVELOPMENT_BUILD
        protected virtual void OnGUI()
        {
            if (statisticsGUI) OnGUIStatistics();
        }
#endif

        /// <summary>
        /// 日志统计
        /// </summary>
        protected virtual void OnLogStatistics()
        {

            if (ClientConnected())
            {
                //TODO 待完善 NetworkTime
                // string log = "kcp CLIENT @ time: " + NetworkTime.localTime + "\n";
                // log += $"  MaxSendRate: {PrettyBytes(client.peer.MaxSendRate)}/s\n";
                // log += $"  MaxRecvRate: {PrettyBytes(client.peer.MaxReceiveRate)}/s\n";
                // log += $"  SendQueue: {client.peer.SendQueueCount}\n";
                // log += $"  ReceiveQueue: {client.peer.ReceiveQueueCount}\n";
                // log += $"  SendBuffer: {client.peer.SendBufferCount}\n";
                // log += $"  ReceiveBuffer: {client.peer.ReceiveBufferCount}\n\n";
                // Debug.Log(log);
            }
        }

        public override string ToString() => "KCP";
    }
}
//#endif MIRROR <- commented out because MIRROR isn't defined on first import yet
