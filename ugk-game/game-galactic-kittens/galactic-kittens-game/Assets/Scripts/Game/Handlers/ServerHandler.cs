using Common.Network;
using Common.Tools;
using Game.Manager;
using Google.Protobuf;
using UGK.Common.Tools;

namespace UGK.Game.Handlers
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
           // Log.Debug($" receive server heart: {ugkMessage.TimeStamp} {response}");
        }
        
        
        /// <summary>
        /// 绑定游戏返回
        /// </summary>
        [MessageMap((int)MID.BindGameConnectRes)]
        private static void BindGame(Player player, UgkMessage ugkMessage)
        {
            var response = new BindGameConnectResponse();
            response.MergeFrom(ugkMessage.Bytes);
            Log.Debug($"{player.Id} 绑定游戏:  {response}");
        }
    }
}