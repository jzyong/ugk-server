using UnityEngine;
using System;
using System.Collections.Generic;
using System.Reflection;
using Google.Protobuf;
using kcp2k;
using System.Linq;
using Tools;
using UnityEngine.SceneManagement;

namespace Network
{
    /// <summary>
    /// <para>网络管理</para>
    /// KcpClient，KcpPeer，NetworkClient应该自己封装，在Mirror上的基础开发层层回调，头都绕晕了
    /// 
    /// </summary>
    [DisallowMultipleComponent]
    public class NetworkManager : MonoBehaviour
    {
        // transport layer 
        [Header("Network Info")]
        [Tooltip("Transport component attached to this object that server and client will use to connect")]
        public Transport transport;

        /// <summary>Server's address for clients to connect to. </summary>
        [Tooltip(
            "Network Address where the client should connect to the server. Server does not use this for anything.")]
        public string networkAddress = "192.168.110.2";

        [Tooltip("服务器端口")] public ushort port = 5000;

        // 消息序列号
        private UInt32 _seq;

        /// <summary>
        /// 消息处理
        /// </summary>
        public delegate void MessageHandler(UgkMessage ugkMessage);

        /// <summary>
        /// 消息处理器
        /// </summary>
        private Dictionary<MID, MessageHandler> messageHandlers;


        /// <summary>The one and only NetworkManager </summary>
        public static NetworkManager Singleton { get; internal set; }


        /// <summary>True if the server is running or client is connected/connecting.</summary>
        public bool isNetworkActive => NetworkClient.active;


        public void Awake()
        {
            // Don't allow collision-destroyed second instance to continue.
            Log.Info = Debug.Log;
            Log.Error = Debug.LogError;
            Log.Warning = Debug.LogWarning;
            if (!InitializeSingleton()) return;
        }

        public void Start()
        {
            //TODO 临时测试
            StartClient();
        }

        public void Update()
        {
        }

        // virtual so that inheriting classes' LateUpdate() can call base.LateUpdate() too
        public void LateUpdate()
        {
        }


        //
        void SetupClient()
        {
            InitializeSingleton();
            CreateMessageHandlersDictionary();
        }

        /// <summary>Starts the client, connects it to the server with networkAddress. </summary>
        public void StartClient()
        {
            if (NetworkClient.active)
            {
                Debug.LogWarning("Client already started.");
                return;
            }


            SetupClient();


            if (string.IsNullOrWhiteSpace(networkAddress))
            {
                Debug.LogError("Must set the Network Address field in the manager");
                return;
            }
            // Debug.Log($"NetworkManager StartClient address:{networkAddress}");

            NetworkClient.Connect(networkAddress, port);
        }


        /// <summary>Stops and disconnects the client. </summary>
        public void StopClient()
        {
            // ask client -> transport to disconnect.
            // handle voluntary and involuntary disconnects in OnClientDisconnect.
            //
            //   StopClient
            //     NetworkClient.Disconnect
            //       Transport.Disconnect
            //         ...
            //       Transport.OnClientDisconnect
            //     NetworkClient.OnTransportDisconnect
            //   NetworkManager.OnClientDisconnect
            NetworkClient.Disconnect();

            // UNET invoked OnDisconnected cleanup immediately.
            // let's keep it for now, in case any projects depend on it.
            // TODO simply remove this in the future.
        }

        // called when quitting the application by closing the window / pressing
        // stop in the editor. virtual so that inheriting classes'
        // OnApplicationQuit() can call base.OnApplicationQuit() too
        public virtual void OnApplicationQuit()
        {
            // stop client first
            // (we want to send the quit packet to the server instead of waiting
            //  for a timeout)
            if (NetworkClient.isConnected)
            {
                StopClient();
                //Debug.Log("OnApplicationQuit: stopped client");
            }


            // Call ResetStatics to reset statics and singleton
            ResetStatics();
        }


        // @
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

            //Debug.Log("NetworkManager created singleton (DontDestroyOnLoad)");
            Singleton = this;
            if (Application.isPlaying)
            {
                // Force the object to scene root, in case user made it a child of something
                // in the scene since DDOL is only allowed for scene root objects
                transform.SetParent(null);
                DontDestroyOnLoad(gameObject);
            }

            // set active transport AFTER setting singleton.
            // so only if we didn't destroy ourselves.

            // This tries to avoid missing transport errors and more clearly tells user what to fix.
            if (transport == null)
                if (TryGetComponent(out Transport newTransport))
                {
                    Debug.LogWarning(
                        $"No Transport assigned to Network Manager - Using {newTransport} found on same object.");
                    transport = newTransport;
                }
                else
                {
                    Debug.LogError("No Transport on Network Manager...add a transport and assign it.");
                    return false;
                }

            Transport.active = transport;
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

        public void OnDestroy()
        {
            //Debug.Log("NetworkManager destroyed");
        }

        /// <summary>
        /// 发送消息
        /// </summary>
        /// <param name="mid"></param>
        /// <param name="message"></param>
        public void Send(MID mid, IMessage message)
        {
            Send(mid, message.ToByteArray());
        }

        public void Send(MID mid, byte[] data)
        {
            // Send((int) mid, data);

            // 消息长度4+消息id4+序列号4+时间戳8+protobuf消息体
            byte[] msgLength = BitConverter.GetBytes(data.Length + 16);
            byte[] msgId = BitConverter.GetBytes((int)mid);
            ++_seq;
            byte[] seq = BitConverter.GetBytes(_seq);
            long time = 0; //TODO 时间戳生成
            byte[] timeStamp = BitConverter.GetBytes(time);
            byte[] datas = new byte[20 + data.Length];

            Array.Copy(msgLength, 0, datas, 0, msgLength.Length);
            Array.Copy(msgId, 0, datas, 4, msgId.Length);
            Array.Copy(seq, 0, datas, 8, seq.Length);
            Array.Copy(timeStamp, 0, datas, 12, seq.Length);
            Array.Copy(data, 0, datas, 20, data.Length);
            ArraySegment<byte> segment = new ArraySegment<byte>(datas);
            Transport.active.ClientSend(segment);
        }

        /// <summary>
        /// 创建消息处理
        ///
        /// 参考：
        /// </summary>
        /// <param name="messageHandlerGroupId"></param>
        /// <exception cref="NonStaticHandlerException"></exception>
        protected void CreateMessageHandlersDictionary()
        {
            MethodInfo[] methods = FindMessageHandlers();

            messageHandlers = new Dictionary<MID, MessageHandler>(methods.Length);
            foreach (MethodInfo method in methods)
            {
                MessageMapAttribute attribute = method.GetCustomAttribute<MessageMapAttribute>();

                if (!method.IsStatic)
                    throw new NonStaticHandlerException(method.DeclaringType, method.Name);

                Delegate clientMessageHandler = Delegate.CreateDelegate(typeof(MessageHandler), method, false);
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
                        messageHandlers.Add(attribute.mid, (MessageHandler)clientMessageHandler);
                        Debug.Log($"消息:${attribute.mid}  添加处理器成功 ：${clientMessageHandler.Method.Name}");
                    }
                }
            }
        }

        /// <summary>查找消息处理方法</summary>
        /// <returns>An array containing message handler methods.</returns>
        protected MethodInfo[] FindMessageHandlers()
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
        public MessageHandler GetMessageHandler(UInt32 messageId)
        {
            MID mid = (MID)messageId;
            return messageHandlers[mid];
        }
    }
}