using System;
using Common.Network.Serialize;
using Common.Tools;
using Google.Protobuf;
using UnityEngine;

namespace Common.Network
{
    /// <summary>
    /// 网络连接状态  
    /// </summary>
    public enum ConnectState
    {
        None,

        // connecting between Connect() and OnTransportConnected()
        Connecting,
        Connected,

        // disconnecting between Disconnect() and OnTransportDisconnected()
        Disconnecting,
        Disconnected
    }


    /// <summary>NetworkClient with connection to server.</summary>
    public class NetworkClient
    {
        /// <summary>
        /// 对应的传输层
        /// </summary>
        public Transport Transport { get; set; }

        ConnectState connectState = ConnectState.None;
        public bool active => connectState == ConnectState.Connecting || connectState == ConnectState.Connected;
        public bool isConnecting => connectState == ConnectState.Connecting;
        public bool isConnected => connectState == ConnectState.Connected;

        public Action OnConnectedEvent;
        public Action OnDisconnectedEvent;
        public Action<TransportError, string> OnErrorEvent;

        public UgkMessage HeartRequest { get; set; }

        //消息批量缓存
        private Batcher _batcher = new Batcher();


        void AddTransportHandlers()
        {
            RemoveTransportHandlers();

            Transport.OnClientConnected += OnTransportConnected;
            //Transport.OnClientDataReceived += OnTransportData; //Game自定义实现
            Transport.OnClientDisconnected += OnTransportDisconnected;
            Transport.OnClientError += OnTransportError;
            Transport.SendHeart += SendHeart;
        }

        void RemoveTransportHandlers()
        {
            Transport.OnClientConnected -= OnTransportConnected;
            // Transport.OnClientDataReceived -= OnTransportData; //Game自定义实现
            Transport.OnClientDisconnected -= OnTransportDisconnected;
            Transport.OnClientError -= OnTransportError;
            Transport.SendHeart -= SendHeart;
        }


        /// <summary>Connect client to a NetworkServer by address. @</summary>
        public void Connect(string address, ushort port)
        {
            Transport.enabled = true;
            AddTransportHandlers();
            connectState = ConnectState.Connecting;
            Transport.ClientConnect(address, port);
        }

        public void Connect(Uri uri, ushort port)
        {
            Transport.enabled = true;
            AddTransportHandlers();
            connectState = ConnectState.Connecting;
            Transport.ClientConnect(uri, port);
        }

        public void Disconnect()
        {
            if (connectState != ConnectState.Connecting &&
                connectState != ConnectState.Connected)
                return;
            connectState = ConnectState.Disconnecting;

            Transport.ClientDisconnect();
        }

        /// <summary>
        /// 连接创建
        /// </summary>
        void OnTransportConnected()
        {
            connectState = ConnectState.Connected;
            OnConnectedEvent?.Invoke();
        }

        void OnTransportDisconnected()
        {
            if (connectState == ConnectState.Disconnected) return;
            OnDisconnectedEvent?.Invoke();
            connectState = ConnectState.Disconnected;
            RemoveTransportHandlers();
        }

        void OnTransportError(TransportError error, string reason)
        {
            Log.Warn($"Client Transport Error: {error}: {reason}. This is fine.");
            OnErrorEvent?.Invoke(error, reason);
        }

        /// <summary>
        /// 发送心跳
        /// </summary>
        private void SendHeart()
        {
            if (HeartRequest != null)
            {
                SendMsg(HeartRequest);
                // Debug.Log("请求心跳");
            }
        }

        /// <summary>
        /// 发送消息
        /// </summary>
        /// <param name="playerId"></param>
        /// <param name="mid"></param>
        /// <param name="seq"></param>
        /// <param name="message"></param>
        /// <returns></returns>
        public bool SendMsg(long playerId, int mid, IMessage message, uint seq = 0)
        {
            var data = message.ToByteArray();
            // 消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体
            byte[] msgLength = BitConverter.GetBytes(data.Length + 24);
            byte[] playerIdBytes = BitConverter.GetBytes(playerId);
            byte[] msgId = BitConverter.GetBytes(mid);
            byte[] seqBytes = BitConverter.GetBytes(seq);
            long time = (long)(Time.unscaledTime * 1000);
            byte[] timeStamp = BitConverter.GetBytes(time);
            byte[] datas = new byte[28 + data.Length];

            Array.Copy(msgLength, 0, datas, 0, msgLength.Length);
            Array.Copy(playerIdBytes, 0, datas, 4, playerIdBytes.Length);
            Array.Copy(msgId, 0, datas, 12, msgId.Length);
            Array.Copy(seqBytes, 0, datas, 16, seqBytes.Length);
            Array.Copy(timeStamp, 0, datas, 20, timeStamp.Length);
            Array.Copy(data, 0, datas, 28, data.Length);
            ArraySegment<byte> segment = new ArraySegment<byte>(datas);
            // Transport.ClientSend(segment);
            _batcher.AddMessage(segment);
            return true;
        }

        /// <summary>
        /// 发送消息
        /// </summary>
        /// <param name="ugkMessage"></param>
        /// <returns></returns>
        public bool SendMsg(UgkMessage ugkMessage)
        {
            //心跳消息忽略，用于进行是否连接判断
            if (!isConnected && ugkMessage.MessageId != 1)
            {
                Log.Warn($"message{ugkMessage.MessageId} seq={ugkMessage.Seq} send fail,network close");
                return false;
            }

            var data = ugkMessage.Bytes;
            // 消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体
            byte[] msgLength = BitConverter.GetBytes(data.Length + 24);
            byte[] playerIdBytes = BitConverter.GetBytes(ugkMessage.PlayerId);
            byte[] msgId = BitConverter.GetBytes(ugkMessage.MessageId);
            byte[] seq = BitConverter.GetBytes(ugkMessage.Seq);
            long time = (long)(Time.time * 1000);
            byte[] timeStamp = BitConverter.GetBytes(time);
            byte[] datas = new byte[28 + data.Length];

            Array.Copy(msgLength, 0, datas, 0, msgLength.Length);
            Array.Copy(playerIdBytes, 0, datas, 4, playerIdBytes.Length);
            Array.Copy(msgId, 0, datas, 12, msgId.Length);
            Array.Copy(seq, 0, datas, 16, seq.Length);
            Array.Copy(timeStamp, 0, datas, 20, seq.Length);
            Array.Copy(data, 0, datas, 28, data.Length);
            ArraySegment<byte> segment = new ArraySegment<byte>(datas);
            // Transport.ClientSend(segment);
            _batcher.AddMessage(segment);
            return true;
        }

        /// <summary>
        /// 批量发送消息
        /// </summary>
        public void BatchSendMsg()
        {
            //批量发送数据
            if (_batcher.HasMessage())
            {
                // make and send as many batches as necessary from the stored
                // messages.
                using (NetworkWriterPooled writer = NetworkWriterPool.Get())
                {
                    // make a batch with our local time (double precision)
                    while (_batcher.GetBatch(writer))
                    {
                        ArraySegment<byte> segment = writer.ToArraySegment();
                        Transport.ClientSend(segment);
                        // reset writer for each new batch
                        writer.Position = 0;
                    }
                }
            }
        }


        // shutdown ////////////////////////////////////////////////////////////
        /// <summary>Shutdown the client.</summary>
        // RuntimeInitializeOnLoadMethod -> fast playmode without domain reload 
        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.BeforeSceneLoad)]
        public void Shutdown()
        {
            // reset statics
            connectState = ConnectState.None;

            // clear events. someone might have hooked into them before, but
            // we don't want to use those hooks after Shutdown anymore.
            OnConnectedEvent = null;
            OnDisconnectedEvent = null;
            OnErrorEvent = null;
        }

        // // GUI /////////////////////////////////////////////////////////////////
        // // called from NetworkManager to display timeline interpolation status.
        // // useful to indicate catchup / slowdown / dynamic adjustment etc.
        // public static void OnGUI()
        // {
        //     // only if in world
        //     if (!ready) return;
        //
        //     GUILayout.BeginArea(new Rect(10, 5, 800, 50));
        //
        //     GUILayout.BeginHorizontal("Box");
        //     GUILayout.Label("Snapshot Interp.:");
        //     // color while catching up / slowing down
        //     if (localTimescale > 1) GUI.color = Color.green; // green traffic light = go fast
        //     else if (localTimescale < 1) GUI.color = Color.red; // red traffic light = go slow
        //     else GUI.color = Color.white;
        //     GUILayout.Box($"timeline: {localTimeline:F2}");
        //     GUILayout.Box($"buffer: {snapshots.Count}");
        //     GUILayout.Box($"DriftEMA: {NetworkClient.driftEma.Value:F2}");
        //     GUILayout.Box($"DelTimeEMA: {NetworkClient.deliveryTimeEma.Value:F2}");
        //     GUILayout.Box($"timescale: {localTimescale:F2}");
        //     GUILayout.Box($"BTM: {snapshotSettings.bufferTimeMultiplier:F2}");
        //     GUILayout.Box($"RTT: {NetworkTime.rtt * 1000:000}");
        //     GUILayout.EndHorizontal();
        //
        //     GUILayout.EndArea();
        // }
    }
}