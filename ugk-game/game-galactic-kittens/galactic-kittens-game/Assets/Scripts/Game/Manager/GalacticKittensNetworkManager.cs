﻿using System;
using Common.Network;
using Google.Protobuf;
using UnityEngine;

namespace Game.Manager
{
    /// <summary>
    /// 自定义网络管理
    /// </summary>
    public class GalacticKittensNetworkManager : NetworkManager<Player>
    {
        public static GalacticKittensNetworkManager singleton { get; private set; }


        public override void Awake()
        {
            base.Awake();
            Application.targetFrameRate = 30;
            singleton = this;
        }

        public override void Start()
        {
            base.Start();
            // 开发环境读取静态配置， 需要连接多个网关，网关地址从外部传入
            if (Application.platform == RuntimePlatform.WindowsEditor)
            {
                StartClient();
            }
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
    }
}