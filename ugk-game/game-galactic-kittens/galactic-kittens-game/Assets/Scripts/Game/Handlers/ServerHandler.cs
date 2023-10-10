using System;
using Game.Manager;
using Google.Protobuf;
using Common.Network;
using Common.Tools;


namespace Game.Handlers
{
    /// <summary>
    /// 服务器消息处理器
    /// </summary>
    internal class ServerHandler
    {
        /// <summary>
        /// 心跳
        /// </summary>
        [MessageMap((int)MID.ServerHeartRes)]
        private static void Heart(Player player, UgkMessage ugkMessage)
        {
            var response = new ServerHeartResponse();
            response.MergeFrom(ugkMessage.Bytes);
            Log.Debug($" receive heart: {ugkMessage.TimeStamp} {response}");
        }
    }
}