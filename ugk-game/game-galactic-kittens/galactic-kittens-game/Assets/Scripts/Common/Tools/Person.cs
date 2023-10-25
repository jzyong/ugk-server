using System;
using Common.Network;

namespace Common.Tools
{
    public class Person
    {
        /// <summary>
        /// 玩家ID
        /// </summary>
        public Int64 Id { get; set; }

        /// <summary>
        /// 昵称
        /// </summary>
        public String Nick { get; set; }

        /// <summary>
        /// 等级
        /// </summary>
        public UInt32 Level { get; set; }

        /// <summary>
        /// 经验
        /// </summary>
        public UInt32 Exp { get; set; }
        
        /// <summary>
        /// 网关客户端
        /// </summary>
        public NetworkClient GateClient;


        /// <summary>
        /// 网关地址
        /// </summary>
        public String GateUrl { get; set; }

        /// <summary>
        /// 大厅ID
        /// </summary>
        public uint LobbyId { get; set; }

        
        

    }
}