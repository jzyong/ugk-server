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

            PlayerManager.Instance.SendMsg(player, MID.GalacticKittensFireRes, response, ugkMessage.Seq);
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

            //TODO 使用护盾逻辑,获取玩家飞船，并添加护盾 ，广播护盾消息 GalacticKittensShipShieldStateResponse
            GalacticKittensUseShieldResponse response = new GalacticKittensUseShieldResponse()
            {
                Result = new MessageResult()
                {
                    Status = 200,
                    Msg = "Success"
                }
            };

            PlayerManager.Instance.SendMsg(player, MID.GalacticKittensFireRes, response, ugkMessage.Seq);
        }

        /// <summary>
        /// 使用护盾
        /// </summary>
        [MessageMap((int)MID.GalacticKittensShipMoveStateReq)]
        private static void ShipMoveState(Player player, UgkMessage ugkMessage)
        {
            var request = new GalacticKittensShipMoveStateRequest();
            request.MergeFrom(ugkMessage.Bytes);
            Log.Trace($"  {player.Id}-{player.Nick} ship state {request.State}");

            var spaceShip = RoomManager.Instance.GetSpaceShip(player.Id);
            if (spaceShip == null)
            {
                GalacticKittensShipMoveStateResponse response2 = new GalacticKittensShipMoveStateResponse()
                {
                    Result = new MessageResult()
                    {
                        Status = 500,
                        Msg = "Request ship id not exist"
                    }
                };
                PlayerManager.Instance.SendMsg(player, MID.GalacticKittensShipMoveStateRes, response2, ugkMessage.Seq);
                return;
            }

            GalacticKittensShipMoveStateResponse response = new GalacticKittensShipMoveStateResponse()
            {
                ShipId = player.Id, //飞船id等于玩家id
                State = request.State
            };
            PlayerManager.Instance.SendMsg(player, MID.GalacticKittensShipMoveStateRes, response, ugkMessage.Seq);
            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensShipMoveStateRes, response,
                excludePredicate: id => id == player.Id);
        }
    }
}