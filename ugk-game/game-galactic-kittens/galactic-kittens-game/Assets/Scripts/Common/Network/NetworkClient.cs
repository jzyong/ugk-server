using System;
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

        double lastSendTime;

        ConnectState connectState = ConnectState.None;
        public bool active => connectState == ConnectState.Connecting || connectState == ConnectState.Connected;
        public bool isConnecting => connectState == ConnectState.Connecting;
        public bool isConnected => connectState == ConnectState.Connected;

        public Action OnConnectedEvent;
        public Action OnDisconnectedEvent;
        public Action<TransportError, string> OnErrorEvent;
        public Action SendHeart;


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
            Debug.LogWarning($"Client Transport Error: {error}: {reason}. This is fine.");
            OnErrorEvent?.Invoke(error, reason);
        }


        // shutdown ////////////////////////////////////////////////////////////
        /// <summary>Shutdown the client.</summary>
        // RuntimeInitializeOnLoadMethod -> fast playmode without domain reload 
        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.BeforeSceneLoad)]
        public void Shutdown()
        {
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