using System;
using System.Collections.Generic;
using Common.Network;
using Common.Tools;
using UnityEngine;

namespace Game.Manager
{
    /// <summary>
    /// 房间管理
    /// </summary>
    public class RoomManager : SingletonPersistent<RoomManager>
    {
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
            //TODO 构建飞船对象,添加 SnapTransform
            foreach (var player in players)
            {
                //TODO 根据角色创建对应的实体对象，添加SnapTransform组件

                GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                    new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                    {
                        OwnerId = player.Id,
                        Id = SyncId++,
                        ConfigId = 1,
                        //
                        // SyncPayload = ; //TODO
                    };

                spawnResponse.Spawn.Add(spawnInfo);
            }

            PlayerManager.Singleton.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
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
            //TODO 获取玩家位置，产出子弹 ，子弹碰撞监测
        }

        /// <summary>
        /// 角色死亡
        /// </summary>
        /// <param name="killerId"></param>
        /// <param name="dieId"></param>
        public void RoleDie(long killerId, long dieId)
        {
            //TODO 清除对象，发送消息
            GalacticKittensObjectDieResponse response = new GalacticKittensObjectDieResponse()
            {
                KillerId = killerId,
                Id = dieId,
                // OwnerId = //TOD
            };

            PlayerManager.Singleton.BroadcastMsg(MID.GalacticKittensObjectDieRes, response);
        }


        /// <summary>
        /// 游戏结束 TODO 待测试
        /// </summary>
        private void GameFinishReq()
        {
            var client =
                new GalacticKittensMatchService.GalacticKittensMatchServiceClient(GalacticKittensNetworkManager
                    .singleton.MatchChannel);
            var request = new GalacticKittensGameFinishRequest()
            {
                RoomId = GalacticKittensNetworkManager.singleton.ServerId
            };
            var response = client.gameFinishAsync(request).ResponseAsync.Result;
            Log.Info($"game finish :{response}");
            if (response.Result?.Status != 200)
            {
                Log.Error($"game finish error:{response.Result?.Msg}");
                return;
            }

            // 解绑玩家网关映射
            PlayerManager.Singleton.BindGateGameMapReq(false);

            Application.Quit();
        }
    }
}