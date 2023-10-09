using System;
using System.Collections.Generic;
using Common.Network;
using Google.Protobuf;
using Grpc.Core;
using kcp2k;
using UnityEngine;

namespace Game.Manager
{
    /// <summary>
    /// 自定义网络管理
    /// </summary>
    public class GalacticKittensNetworkManager : NetworkManager<Player>
    {
        public static GalacticKittensNetworkManager singleton { get; private set; }

        // 后面从agent通过参数传输过来？
        [SerializeField] [Tooltip("匹配服Grpc地址")]
        private String matchGrpcUrl = "127.0.0.1:4000";


        //匹配服channel
        private Channel matchChannel;

        //大厅服channel
        private Dictionary<uint, Channel> lobbyChannels = new Dictionary<uint, Channel>(2);


        public override void Awake()
        {
            base.Awake();
            Application.targetFrameRate = 30;
            singleton = this;

           
        }

        public override void Start()
        {
            base.Start();
            // 初始化Grpc
            ServerInfoRequest();
            // // 开发环境读取静态配置， 需要连接多个网关，网关地址从外部传入
            // if (Application.platform == RuntimePlatform.WindowsEditor)
            // {
            //     StartClient();
            // }
        }

        //获取消息并处理 
        protected override void OnTransportData(ArraySegment<byte> data)
        {
            using (UgkMessage ugkMessage = UgkMessagePool.Get())
            {
                //  `消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
                var bytes = data.Array;
                Int32 messageLength = BitConverter.ToInt32(bytes, 0);
                ugkMessage.PlayerId = BitConverter.ToInt64(bytes, 4);
                ugkMessage.MessageId = BitConverter.ToUInt32(bytes, 12);
                ugkMessage.Seq = BitConverter.ToUInt32(bytes, 16);
                ugkMessage.TimeStamp = BitConverter.ToInt64(bytes, 20);

                // Debug.Log(
                //     $"{ugkMessage.PlayerId}收到消息 ID={ugkMessage.MessageId} Seq={ugkMessage.Seq} timeStamp={ugkMessage.TimeStamp}");
                var handler = Singleton.GetMessageHandler((int)ugkMessage.MessageId);
                if (handler == null)
                {
                    Debug.LogWarning($"消息{(MID)ugkMessage.MessageId}处理方法未实现");
                }
                else
                {
                    //TODO 获取玩家
                    var protoData = new byte[messageLength - 24];
                    Array.Copy(bytes, 28, protoData, 0, protoData.Length);
                    ugkMessage.Bytes = protoData;
                    var player = PlayerManager.Singleton.GetPlayer(ugkMessage.PlayerId);
                    handler(player, ugkMessage);
                }
            }
        }

        protected override UgkMessage GetServerHeartRequest()
        {
            if (heartRequest == null)
            {
                var ugkMessage = UgkMessagePool.Get();

                ServerHeartRequest request = new ServerHeartRequest()
                {
                    //TODO 完整信息，ID通过agent-manager 生成
                    Server = new ServerInfo()
                    {
                        Id = 1,
                        Name = "GalacticKittensGame",
                    }
                };
                ugkMessage.Bytes = request.ToByteArray();
                ugkMessage.MessageId = (int)MID.ServerHeartReq;
                heartRequest = ugkMessage;
            }

            return heartRequest;
        }


        /// <summary>
        /// 匹配服grpc连接
        /// </summary>
        public Channel MatchChannel
        {
            get
            {
                if (matchChannel == null || matchChannel.State == ChannelState.Shutdown ||
                    matchChannel.State == ChannelState.TransientFailure)
                {
                    var urlPort = matchGrpcUrl.Split(":");
                    matchChannel = new Channel(urlPort[0], Int32.Parse(urlPort[1]), ChannelCredentials.Insecure);
                    Debug.Log($"创建Match连接：{matchGrpcUrl}");
                }

                return matchChannel;
            }
        }

        /// <summary>
        /// 匹配服grpc连接
        /// </summary>
        public Channel GetLobbyChannel(uint id)
        {
            Channel channel;
            lobbyChannels.TryGetValue(id, out channel);
            return channel;
        }

        /// <summary>
        /// 请求服务器信息,并创建相应的grpc和kcp连接
        /// </summary>
        private void ServerInfoRequest()
        {
            var client = new ServerService.ServerServiceClient(MatchChannel);
            var response = client.getServerInfoAsync(new GetServerInfoRequest()).ResponseAsync.Result;
            Debug.Log($"服务器信息：{response}");
            foreach (var serverInfo in response.ServerInfo)
            {
                if (serverInfo.Name.Equals("lobby"))
                {
                    var urlPort = serverInfo.GrpcUrl.Split(":");
                    var lobbyChannel = new Channel(urlPort[0], Int32.Parse(urlPort[1]), ChannelCredentials.Insecure);
                    lobbyChannels.Add(serverInfo.Id, lobbyChannel);
                    Debug.Log($"创建{serverInfo.Name}连接：{serverInfo.GrpcUrl}");
                }
                else if (serverInfo.Name.Equals("gate"))
                {
                    var kcpTransport = gameObject.AddComponent<KcpTransport>();
                    var urlPort = serverInfo.GrpcUrl.Split(":");
                    kcpTransport.networkAddress = urlPort[0];
                    kcpTransport.port = UInt16.Parse(urlPort[1]);
                    Debug.Log($"创建{serverInfo.Name}连接：{serverInfo.GrpcUrl}");
                }
            }
            StartClient();
        }
    }
}