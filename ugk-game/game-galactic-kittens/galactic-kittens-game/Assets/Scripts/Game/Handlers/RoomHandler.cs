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
        /// 开火请求 ,只有玩家控制的对象请求，子弹服务器生成推送
        /// </summary>
        [MessageMap((int)MID.GalacticKittensFireReq)]
        private static void Fire(Player player, UgkMessage ugkMessage)
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
        
        
        /// <summary>
        /// 使用护盾
        /// </summary>
        [MessageMap((int)MID.GalacticKittensUseShieldReq)]
        private static void UseShield(Player player, UgkMessage ugkMessage)
        {
            var request = new GalacticKittensUseShieldRequest();
            request.MergeFrom(ugkMessage.Bytes);
            Log.Trace($" receive use shield {player.Id}-{player.Nick}");

           //TODO 使用护盾逻辑,获取玩家飞船，并添加护盾
            GalacticKittensUseShieldResponse response = new GalacticKittensUseShieldResponse()
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