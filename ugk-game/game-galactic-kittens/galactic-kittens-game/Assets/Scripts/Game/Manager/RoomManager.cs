using System.Collections.Generic;
using Common.Network.Sync;
using Common.Tools;
using Game.Room.Player;
using UnityEngine;

namespace Game.Manager
{
    /// <summary>
    /// 房间管理
    /// </summary>
    public class RoomManager : SingletonPersistent<RoomManager>
    {
        [SerializeField] [Tooltip("飞船，服务器只需要一个简单对象即可")] private SpaceShip _spaceShip;

        /// <summary>
        /// 飞船出生坐标
        /// </summary>
        private readonly Vector3[] shipSpawnPositions = new[]
            { new Vector3(-8, 4), new Vector3(-8, 1.5f), new Vector3(-8, -1f), new Vector3(-8, -3.5f) };

        /// <summary>
        /// 对象同步Id
        /// </summary>
        private long SyncId { get; set; }

        private void Update()
        {
            SpawnEnemy();
        }

        /// <summary>
        /// 出生玩家
        /// </summary>
        public void SpawnPlayers(List<Player> players)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();
            for (int i = 0; i < players.Count; i++)
            {
                // 根据角色创建对应的实体对象，添加SnapTransform组件
                var player = players[i];
                var spawnPosition = shipSpawnPositions[i];
                var spaceShip = Instantiate(_spaceShip, spawnPosition, Quaternion.identity,
                    Instance.transform);
                var snapTransform = spaceShip.GetComponent<SnapTransform>();
                snapTransform.Id = player.Id;
                GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                    new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                    {
                        OwnerId = player.Id,
                        Id = player.Id,
                        ConfigId = (uint)player.CharacterId, //0-3
                        Position = new Vector3D()
                        {
                            X = spawnPosition.x,
                            Y = spawnPosition.y,
                            Z = spawnPosition.z
                        }
                    };
                SyncManager.Instance.AddSnapTransform(snapTransform); //添加同步对象
                spawnResponse.Spawn.Add(spawnInfo);
                Log.Info($"{player.Id} born in {spawnPosition}");
            }

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }

        /// <summary>
        /// 出生敌人
        /// </summary>
        private void SpawnEnemy()
        {
            //TODO 规则是什么？
        }

        /// <summary>
        /// 创建子弹
        /// </summary>
        /// <param name="player"></param>
        public void SpawnBullet(Player player)
        {
            //TODO 获取玩家位置，产出子弹 ，子弹碰撞监测,prefab
        }

        /// <summary>
        /// 对象死亡
        /// </summary>
        public void DespawnObject(long killerId, long dieId, bool removeObject = true)
        {
            if (removeObject)
            {
                //移除的对象全部使用预测
                SyncManager.Instance.RemovePredictionTransform(dieId);
            }


            //TODO 清除对象，发送消息
            GalacticKittensObjectDieResponse response = new GalacticKittensObjectDieResponse()
            {
                KillerId = killerId,
                Id = dieId,
                // OwnerId = //TOD
            };

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectDieRes, response);
        }


        /// <summary>
        /// 游戏结束 TODO 待测试
        /// </summary>
        private void GameFinishReq()
        {
            var client =
                new GalacticKittensMatchService.GalacticKittensMatchServiceClient(GalacticKittensNetworkManager
                    .Instance.MatchChannel);
            var request = new GalacticKittensGameFinishRequest()
            {
                RoomId = GalacticKittensNetworkManager.Instance.ServerId
            };
            var response = client.gameFinishAsync(request).ResponseAsync.Result;
            Log.Info($"game finish :{response}");
            if (response.Result?.Status != 200)
            {
                Log.Error($"game finish error:{response.Result?.Msg}");
                return;
            }

            // 解绑玩家网关映射
            PlayerManager.Instance.BindGateGameMapReq(false);

            Application.Quit();
        }
    }
}