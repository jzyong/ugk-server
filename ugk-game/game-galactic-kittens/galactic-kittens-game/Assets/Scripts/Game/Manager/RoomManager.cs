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
    public class RoomManager : SingletonInstance<RoomManager>
    {



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