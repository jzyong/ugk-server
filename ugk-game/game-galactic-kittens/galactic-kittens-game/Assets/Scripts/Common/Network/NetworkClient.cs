using System;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using Object = UnityEngine.Object;

namespace Network
{
    /// <summary>
    /// 网络连接状态  TODO
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

    //消息处理委托 TODO 传入proto Message
    public delegate void MessageHandDelegate(int seq, long timeStamp);

    /// <summary>NetworkClient with connection to server.</summary>
    public static partial class NetworkClient
    {
        // time & value snapshot interpolation are separate.
        // -> time is interpolated globally on NetworkClient / NetworkConnection
        // -> value is interpolated per-component, i.e. NetworkTransform.
        // however, both need to be on the same send interval.
        //
        // additionally, server & client need to use the same send interval.
        // otherwise it's too easy to accidentally cause interpolation issues if
        // a component sends with client.interval but interpolates with
        // server.interval, etc.
        public static int sendRate => 30;
        public static float sendInterval => sendRate < int.MaxValue ? 1f / sendRate : 0; // for 30 Hz, that's 33ms
        static double lastSendTime;

        // message handlers by messageId 
        internal static readonly Dictionary<int, MessageHandDelegate> handlers =
            new Dictionary<int, MessageHandDelegate>();


        /// <summary>Client's NetworkConnection to server.  TODO </summary>
        // public static NetworkConnection connection { get; internal set; }

        // NetworkClient state 
        internal static ConnectState connectState = ConnectState.None;

        /// <summary>active is true while a client is connecting/connected either as standalone or as host client. @</summary>
        // (= while the network is active)
        public static bool active => connectState == ConnectState.Connecting ||
                                     connectState == ConnectState.Connected;


        /// <summary>Check if client is connecting (before connected).</summary>
        public static bool isConnecting => connectState == ConnectState.Connecting;

        /// <summary>Check if client is connected (after connecting). </summary>
        public static bool isConnected => connectState == ConnectState.Connected;


        // OnConnected / OnDisconnected used to be NetworkMessages that were
        // invoked. this introduced a bug where external clients could send
        // Connected/Disconnected messages over the network causing undefined
        // behaviour.
        // => public so that custom NetworkManagers can hook into it 
        public static Action OnConnectedEvent;
        public static Action OnDisconnectedEvent;
        public static Action<TransportError, string> OnErrorEvent;


        // initialization ////////////////////////////////////////////////////// @
        static void AddTransportHandlers()
        {
            // community Transports may forget to call OnDisconnected.
            // which could cause handlers to be added twice with +=.
            // ensure we always clear the old ones first.
            // fixes: https://github.com/vis2k/Mirror/issues/3152
            RemoveTransportHandlers();

            // += so that other systems can also hook into it (i.e. statistics)
            Transport.active.OnClientConnected += OnTransportConnected;
            Transport.active.OnClientDataReceived += OnTransportData;
            Transport.active.OnClientDisconnected += OnTransportDisconnected;
            Transport.active.OnClientError += OnTransportError;
            Transport.active.SendHeart += SendHeart;
        }

        //
        static void RemoveTransportHandlers()
        {
            // -= so that other systems can also hook into it (i.e. statistics)
            Transport.active.OnClientConnected -= OnTransportConnected;
            Transport.active.OnClientDataReceived -= OnTransportData;
            Transport.active.OnClientDisconnected -= OnTransportDisconnected;
            Transport.active.OnClientError -= OnTransportError;
            Transport.active.SendHeart -= SendHeart;
        }

        // connect /////////////////////////////////////////////////////////////
        // initialize is called before every connect @
        static void Initialize()
        {
            //TODO 消息注册通过类扫描，反射自动注册
            RegisterMessageHandlers();
            Transport.active.enabled = true;
        }

        /// <summary>Connect client to a NetworkServer by address. @</summary>
        public static void Connect(string address,ushort port)
        {
            Initialize();

            AddTransportHandlers();
            connectState = ConnectState.Connecting;
            Transport.active.ClientConnect(address,port);
        }

        /// <summary>Connect client to a NetworkServer by Uri.</summary>
        public static void Connect(Uri uri,ushort port)
        {
            Initialize();

            AddTransportHandlers();
            connectState = ConnectState.Connecting;
            Transport.active.ClientConnect(uri,port);
        }


        /// <summary>Disconnect from server. </summary>
        public static void Disconnect()
        {
            if (connectState != ConnectState.Connecting &&
                connectState != ConnectState.Connected)
                return;
            connectState = ConnectState.Disconnecting;

            // call Disconnect on the NetworkConnection
            Transport.active.ClientDisconnect();

        }

        // transport events ////////////////////////////////////////////////////
        // called by Transport @
        static void OnTransportConnected()
        {
                //TZODO
                // reset network time stats
                // NetworkTime.ResetStatics();

                // the handler may want to send messages to the client
                // thus we should set the connected state before calling the handler
                connectState = ConnectState.Connected;
                // NetworkTime.UpdateClient();
                OnConnectedEvent?.Invoke();
        }


        // called by Transport 获取消息并处理 
        internal static void OnTransportData(ArraySegment<byte> data)
        {
            using (UgkMessage ugkMessage =UgkMessagePool.Get())
            {
                //  `消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
                var bytes = data.Array;
                Int32 messageLength = BitConverter.ToInt32(bytes, 0);
                ugkMessage.MessageId = BitConverter.ToUInt32(bytes, 4);
                ugkMessage.Seq = BitConverter.ToUInt32(bytes, 8);
                ugkMessage.TimeStamp =  BitConverter.ToInt64(bytes, 12);
            
            
                // Debug.Log($"收到消息 ID={messageId} Seq={seq} timeStamp={timeStamp}");
                var handler = NetworkManager.Singleton.GetMessageHandler(ugkMessage.MessageId);
                if (handler==null)
                {
                    Debug.LogWarning($"消息{(MID)ugkMessage.MessageId}处理方法未实现");
                }
                else
                {
                    var protoData = new byte[messageLength-16];
                    Array.Copy(bytes, 20, protoData, 0, protoData.Length);
                    ugkMessage.Bytes = protoData;
                    handler(ugkMessage);
                }
            }
   
        }

        // called by Transport
        // IMPORTANT: often times when disconnecting, we call this from Mirror
        //            too because we want to remove the connection and handle
        //            the disconnect immediately.
        //            => which is fine as long as we guarantee it only runs once
        //            => which we do by setting the state to Disconnected! @
        internal static void OnTransportDisconnected()
        {
            // StopClient called from user code triggers Disconnected event
            // from transport which calls StopClient again, so check here
            // and short circuit running the Shutdown process twice.
            if (connectState == ConnectState.Disconnected) return;

            // Raise the event before changing ConnectState
            // because 'active' depends on this during shutdown
            //
            // previously OnDisconnected was only invoked if connection != null.
            // however, if DNS resolve fails in Transport.Connect(),
            // OnDisconnected would never be called because 'connection' is only
            // created after the Transport.Connect() call.
            // fixes: https://github.com/MirrorNetworking/Mirror/issues/3365
            OnDisconnectedEvent?.Invoke();

            connectState = ConnectState.Disconnected;
            
            //TODO

            // transport handlers are only added when connecting.
            // so only remove when actually disconnecting.
            RemoveTransportHandlers();
        }

        // transport errors are forwarded to high level @
        static void OnTransportError(TransportError error, string reason)
        {
            // transport errors will happen. logging a warning is enough.
            // make sure the user does not panic.
            Debug.LogWarning($"Client Transport Error: {error}: {reason}. This is fine.");
            OnErrorEvent?.Invoke(error, reason);
        }
        
        /// <summary>
        /// 发送心跳消息
        /// </summary>
        public static void SendHeart()
        {
            HeartRequest request = new HeartRequest();
            NetworkManager.Singleton.Send(MID.HeartReq,request);
        }

        // send ////////////////////////////////////////////////////////////////
        /// <summary>Send a NetworkMessage to the server over the given channel.</summary>
        public static void Send<T>(T message)
            where T :  Object
        {
            //TODO 继承protobuf
           //TODO
        }

        // message handlers //////////////////////////////////////////////////// @
        internal static void RegisterMessageHandlers()
        {
            //TODO 消息注册通过类扫描，反射自动注册
        }

     

        

      

       

        // update //////////////////////////////////////////////////////////////
        // NetworkEarlyUpdate called before any Update/FixedUpdate
        // (we add this to the UnityEngine in NetworkLoop)
        internal static void NetworkEarlyUpdate()
        {
            // process all incoming messages first before updating the world
            if (Transport.active != null)
                Transport.active.ClientEarlyUpdate();

        }

        // NetworkLateUpdate called after any Update/FixedUpdate/LateUpdate
        // (we add this to the UnityEngine in NetworkLoop)
        internal static void NetworkLateUpdate()
        {
            
            // process all outgoing messages after updating the world
            if (Transport.active != null)
                Transport.active.ClientLateUpdate();
        }

      

        // shutdown ////////////////////////////////////////////////////////////
        /// <summary>Shutdown the client.</summary>
        // RuntimeInitializeOnLoadMethod -> fast playmode without domain reload @
        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.BeforeSceneLoad)]
        public static void Shutdown()
        {
           
            handlers.Clear();

            // reset statics
            connectState = ConnectState.None;
            lastSendTime = 0;

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