using Common.Network;
using Common.Tools;
using Game.Manager;
using Google.Protobuf;
using UnityEngine;

namespace UGK.Game.Handlers
{
    /// <summary>
    /// 玩家消息处理器
    /// </summary>
    internal class PlayerHandler
    {
        /// <summary>
        /// 玩家心跳，用于同步服务器时间，检测ping值
        /// </summary>
        [MessageMap((int)MID.HeartReq)]
        private static void Heart(Player player, UgkMessage ugkMessage)
        {
            var request = new HeartRequest();
            request.MergeFrom(ugkMessage.Bytes);
            Log.Debug($" receive player heart: {ugkMessage.TimeStamp} clientTime={request.ClientTime} serverTime={Time.time}");
            var response = new HeartResponse()
            {
                ClientTime = request.ClientTime
            };
            PlayerManager.Instance.SendMsg(player, MID.HeartRes, response, ugkMessage.Seq);
        }
    }
}