using System;
using Common.Network;

namespace Common.Tools
{
    public class Person
    {
        /// <summary>
        /// 玩家ID
        /// </summary>
        public Int64 PlayerId { get; set; }
        
        /// <summary>
        /// 客户端对应的连接
        /// </summary>
        public NetworkClient Client { get; set; }
        
        
        //TODO 发送消息
    }
}