using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Common.Tools;
using Google.Protobuf;
using kcp2k;
using UnityEngine;
using UnityEngine.Serialization;

namespace Common.Network
{
    /// <summary>
    /// <para>网络管理</para>
    ///  <para>KcpClient，KcpPeer，NetworkClient应该自己封装，在Mirror上的基础开发层层回调，头都绕晕了</para>
    /// 管理多个玩家,需要修改 
    /// </summary>
    [DisallowMultipleComponent]
    public class NetworkManager<T> : MonoBehaviour where T : Person
    {
        // 传输层   需要通过参数传入，且是多个网关
        [FormerlySerializedAs("transport")] [Header("网络信息")] [Tooltip("连接多个网关的传输层配置")]
        public Transport[] transports;

        /// <summary>
        /// 消息处理
        /// </summary>
        public delegate void MessageHandler<T>(T player, UgkMessage ugkMessage) where T : Person;

        /// <summary>
        /// 消息处理器
        /// </summary>
        private Dictionary<MID, MessageHandler<T>> messageHandlers;

        /// <summary>
        /// 网关客户端
        /// </summary>
        static Dictionary<String, NetworkClient> gateClients = new Dictionary<string, NetworkClient>(2);

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
        
        
        public static NetworkManager<T> Singleton { get; internal set; }

        public void Awake()
        {
            Log.Info = Debug.Log;
            Log.Error = Debug.LogError;
            Log.Warning = Debug.LogWarning;
            if (!InitializeSingleton()) return;
        }

        public void Start()
        {
            //TODO 临时测试,需要连接多个网关，网关地址从外部传入（怎么传）？
            StartClient();
        }


        /// <summary>
        /// 连接网关
        /// </summary>
        public void StartClient()
        {
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
                networkClient.SendHeart = SendHeart;
                networkClient.Connect(transport.networkAddress,transport.port);
                String url = $"{transport.networkAddress}:{transport.port}";
                gateClients[url] = networkClient;
            }
        }


        /// <summary>Stops and disconnects the client. </summary>
        private void StopClient()
        {
            foreach (var gateClient in gateClients)
            {
                //TODO  NetworkClient.Disconnect();
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
        /// 发送消息 TODO 需要争对每个连接
        /// </summary>
        /// <param name="mid"></param>
        /// <param name="message"></param>
        public void Send(MID mid, IMessage message)
        {
            Send(mid, message.ToByteArray());
        }

        /// <summary>
        ///  TODO  需要争对每个连接
        /// </summary>
        /// <param name="mid"></param>
        /// <param name="data"></param>
        public void Send(MID mid, byte[] data)
        {
            // // Send((int) mid, data);
            // // TODO 移到client中
            //
            // // 消息长度4+消息id4+序列号4+时间戳8+protobuf消息体
            // byte[] msgLength = BitConverter.GetBytes(data.Length + 16);
            // byte[] msgId = BitConverter.GetBytes((int)mid);
            // byte[] seq = BitConverter.GetBytes(0);
            // long time = 0; //TODO 时间戳生成
            // byte[] timeStamp = BitConverter.GetBytes(time);
            // byte[] datas = new byte[20 + data.Length];
            //
            // Array.Copy(msgLength, 0, datas, 0, msgLength.Length);
            // Array.Copy(msgId, 0, datas, 4, msgId.Length);
            // Array.Copy(seq, 0, datas, 8, seq.Length);
            // Array.Copy(timeStamp, 0, datas, 12, seq.Length);
            // Array.Copy(data, 0, datas, 20, data.Length);
            // ArraySegment<byte> segment = new ArraySegment<byte>(datas);
            // Transport.active.ClientSend(segment);
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

            messageHandlers = new Dictionary<MID, MessageHandler<T>>(methods.Length);
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
                        Debug.Log($"消息:${attribute.mid}  添加处理器成功 ：${clientMessageHandler.Method.Name}");
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
        public MessageHandler<T> GetMessageHandler(UInt32 messageId)
        {
            MID mid = (MID)messageId;
            return messageHandlers[mid];
        }
        
        /// <summary>
        /// 发送心跳消息
        /// </summary>
        public  void SendHeart()
        {
            // HeartRequest request = new HeartRequest(); //TODO 发送服务器内部心跳消息 ，需要每个客户端单独发送
            // NetworkManager.Singleton.Send(MID.HeartReq,request);
        }
        
        /// <summary>
        /// 使用unity 主循环更新
        /// </summary>
        public static   void NetworkEarlyUpdate()
        {
            foreach (var pair in gateClients)
            {
                pair.Value.Transport.ClientEarlyUpdate();
            }
        }
        
        /// <summary>
        /// 使用unity 主循环更新
        /// </summary>
        public static   void NetworkLateUpdate()
        {
            foreach (var pair in gateClients)
            {
                pair.Value.Transport.ClientEarlyUpdate();
            }
        }
        
    }
}