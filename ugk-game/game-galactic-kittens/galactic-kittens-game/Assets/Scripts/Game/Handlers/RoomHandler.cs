using Common.Network;
using Common.Tools;
using Game.Manager;
using Google.Protobuf;

namespace Game.Handlers
{
    /// <summary>
    /// 房间内消息处理
    /// </summary>
    public class RoomHandler
    {
        /// <summary>
        /// 玩家心跳，用于同步服务器时间，检测ping值
        /// </summary>
        [MessageMap((int)MID.GalacticKittensFireReq)]
        private static void Heart(Player player, UgkMessage ugkMessage)
        {
            var request = new GalacticKittensFireRequest();
            request.MergeFrom(ugkMessage.Bytes);
            Log.Trace($" receive player fire {player.Id}-{player.Nick}");

            RoomManager.Instance.SpawnBullet(player);
            GalacticKittensFireResponse response = new GalacticKittensFireResponse()
            {
                Result = new MessageResult()
                {
                    Status = 200,
                    Msg = "Success"
                }
            };

            PlayerManager.Singleton.SendMsg(player, MID.GalacticKittensFireRes, response, ugkMessage.Seq);
        }
    }
}