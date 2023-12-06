using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Common.Tools;
using Google.Protobuf;
using kcp2k;
using UnityEngine;
using Log = kcp2k.Log;

namespace Common.Network
{
    /// <summary>
    /// <para>网络管理</para>
    ///  <para>KcpClient，KcpPeer，NetworkClient应该自己封装，在Mirror上的基础开发层层回调，头都绕晕了</para>
    /// <para>子类继承写错了，会导致unity主线程Update卡一段时间，没有任何报错，然后网络超时，服务器断开连接？</para>
    /// </summary>
    [DisallowMultipleComponent]
    [AddComponentMenu("网络/Network Manager")]
    public class NetworkManager<T> : MonoBehaviour where T : Person
    {
        // // 传输层   需要通过参数传入，且是多个网关
        // [Header("网络信息")] [Tooltip("连接多个网关的传输层配置")]
        // public Transport[] transports;

        /// <summary>
        /// 消息处理
        /// </summary>
        public delegate void MessageHandler<in P>(P player, UgkMessage ugkMessage) where P : Person;

        /// <summary>
        /// 消息处理器
        /// </summary>
        private Dictionary<int, MessageHandler<T>> messageHandlers;

        /// <summary>
        /// 网关客户端
        /// </summary>
        Dictionary<String, NetworkClient> gateClients = new Dictionary<string, NetworkClient>(2);

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

        protected UgkMessage heartRequest;

        /// <summary>
        /// 服务存在时间
        /// </summary>
        protected float ExistMaxTime = 3600;


        public static NetworkManager<T> Singleton { get; internal set; }

        public virtual void Awake()
        {
            Log.Info = Tools.Log.Info;
            Log.Error = Tools.Log.Error;
            Log.Warning = Tools.Log.Warn;
            Application.targetFrameRate = 30;
            Application.runInBackground = true;
            if (!InitializeSingleton()) return;
        }

        public virtual void Start()
        {
        }

        public void Update()
        {
            //进程超过最大存在时间，直接退出，防止资源泄漏
            if (Time.time>ExistMaxTime)
            {
                Application.Quit();
            }
        }


        /// <summary>
        /// 连接网关
        /// </summary>
        protected void StartClient()
        {
            //连接地址是通过rpc请求动态添加的
            var transports = gameObject.GetComponents<KcpTransport>();
            if (transports == null)
            {
                Debug.LogError("没有可连接的网关配置");
                return;
            }

            //TODO 断线重连这些？
            foreach (var transport in transports)
            {
                NetworkClient networkClient = new NetworkClient();
                networkClient.Transport = transport;
                networkClient.HeartRequest = GetServerHeartRequest();
                networkClient.Transport.OnClientDataReceived = OnTransportData;
                String url = $"{transport.networkAddress}:{transport.port}";
                networkClient.Connect(transport.networkAddress, transport.port);
                gateClients[url] = networkClient;
            }
        }


        /// <summary>Stops and disconnects the client. </summary>
        private void StopClient()
        {
            foreach (var gateClient in gateClients)
            {
                gateClient.Value.Disconnect();
            }
        }

        public virtual void OnApplicationQuit()
        {
            StopClient();
            ResetStatics();
        }


        /// <summary>
        /// 初始化
        /// </summary>
        /// <returns></returns>
        bool InitializeSingleton()
        {
            if (Singleton != null && Singleton == this)
                return true;

            if (Singleton != null)
            {
                Debug.LogWarning(
                    "Multiple NetworkManagers detected in the scene. Only one NetworkManager can exist at a time. The duplicate NetworkManager will be destroyed.");
                Destroy(gameObject);

                // Return false to not allow collision-destroyed second instance to continue.
                return false;
            }

            Singleton = this;
            if (Application.isPlaying)
            {
                // Force the object to scene root, in case user made it a child of something
                // in the scene since DDOL is only allowed for scene root objects
                transform.SetParent(null);
                DontDestroyOnLoad(gameObject);
            }

            CreateMessageHandlersDictionary();
            NetworkLoop.OnEarlyUpdate = NetworkEarlyUpdate;
            NetworkLoop.OnLateUpdate = NetworkLateUpdate;
            return true;
        }


        // This is the only way to clear the singleton, so another instance can be created.
        // RuntimeInitializeOnLoadMethod -> fast playmode without domain reload
        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.BeforeSceneLoad)]
        public static void ResetStatics()
        {
            // and finally (in case it isn't null already)...
            Singleton = null;
        }


        /// <summary>
        ///  发送消息
        /// </summary>
        public bool Send(NetworkClient networkClient, long playerId, int mid, IMessage message, uint seq = 0)
        {
            return networkClient.SendMsg(playerId, mid, message, seq);
        }

        /// <summary>
        /// 创建消息处理
        ///
        /// 参考：
        /// </summary>
        /// <exception cref="NonStaticHandlerException"></exception>
        private void CreateMessageHandlersDictionary()
        {
            MethodInfo[] methods = FindMessageHandlers();

            messageHandlers = new Dictionary<int, MessageHandler<T>>(methods.Length);
            foreach (MethodInfo method in methods)
            {
                MessageMapAttribute attribute = method.GetCustomAttribute<MessageMapAttribute>();

                if (!method.IsStatic)
                    throw new NonStaticHandlerException(method.DeclaringType, method.Name);

                Delegate clientMessageHandler = Delegate.CreateDelegate(typeof(MessageHandler<T>), method, false);
                if (clientMessageHandler != null)
                {
                    // It's a message handler for Client instances
                    if (messageHandlers.ContainsKey(attribute.mid))
                    {
                        MethodInfo otherMethodWithId = messageHandlers[attribute.mid].GetMethodInfo();
                        throw new DuplicateHandlerException((Int32)attribute.mid, method, otherMethodWithId);
                    }
                    else
                    {
                        messageHandlers.Add(attribute.mid, (MessageHandler<T>)clientMessageHandler);
                        Tools.Log.Debug($"message {attribute.mid} add handler :{clientMessageHandler.Method.Name}");
                    }
                }
            }
        }

        /// <summary>查找消息处理方法</summary>
        /// <returns>An array containing message handler methods.</returns>
        private MethodInfo[] FindMessageHandlers()
        {
            // string thisAssemblyName = Assembly.GetExecutingAssembly().GetName().FullName;

            return Assembly.GetExecutingAssembly().GetTypes().SelectMany(t =>
                    t.GetMethods(BindingFlags.Public | BindingFlags.NonPublic | BindingFlags.Static |
                                 BindingFlags
                                     .Instance)) // Include instance methods in the search so we can show the developer an error instead of silently not adding instance methods to the dictionary
                .Where(m => m.GetCustomAttributes(typeof(MessageMapAttribute), false).Length > 0)
                .ToArray();
        }

        /// <summary>
        /// 获取消息处理器
        /// </summary>
        /// <param name="messageId"></param>
        /// <returns></returns>
        public MessageHandler<T> GetMessageHandler(int messageId)
        {
            MessageHandler<T> handler;
            if (messageHandlers.TryGetValue(messageId, out handler))
            {
                return handler;
            }

            return null;
        }

        /// <summary>
        /// 收到返回消息
        /// </summary>
        /// <param name="data"></param>
        protected virtual void OnTransportData(ArraySegment<byte> data)
        {
        }

        /// <summary>
        /// 发送心跳消息
        /// </summary>
        protected virtual UgkMessage GetServerHeartRequest()
        {
            return null;
        }


        /// <summary>
        /// 使用unity 主循环更新
        /// </summary>
        public void NetworkEarlyUpdate()
        {
            var time = Time.time;
            foreach (var pair in gateClients)
            {
                // Debug.Log($"NetworkEarlyUpdate：{Time.time}");
                pair.Value.Transport.ClientEarlyUpdate();
            }

            if (Time.time - time > 0.01)
            {
                Debug.LogWarning($"NetworkEarlyUpdate耗时：{Time.time - time}");
            }
        }

        /// <summary>
        /// 使用unity 主循环更新
        /// </summary>
        public void NetworkLateUpdate()
        {
            var time = Time.time;
            foreach (var pair in gateClients)
            {
                pair.Value.Transport.ClientLateUpdate();
            }

            if (Time.time - time > 0.01)
            {
                Debug.LogWarning($"NetworkLateUpdate耗时：{Time.time - time}");
            }
        }

        /// <summary>
        /// 获取网关客户端
        /// </summary>
        /// <param name="url"></param>
        /// <returns></returns>
        public NetworkClient GetGateClient(String url)
        {
            if (gateClients.TryGetValue(url,out var client))
            {
                return client;
            }

            return null;
        }
    }
}