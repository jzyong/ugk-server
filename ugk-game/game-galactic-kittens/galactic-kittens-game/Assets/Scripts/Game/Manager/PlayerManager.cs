using System;
using System.Collections.Generic;
using Common.Network;
using Common.Tools;
using Google.Protobuf;
using UnityEngine;
using UnityEngine.SceneManagement;

namespace Game.Manager
{
    /// <summary>
    /// 玩家
    /// </summary>
    public class Player : Person
    {
    }

    /// <summary>
    /// 玩家管理
    /// </summary>
    public class PlayerManager : SingletonInstance<PlayerManager>
    {
        private List<Player> players = new List<Player>(4);

        public Player GetPlayer(Int64 playerId)
        {
            foreach (var player in players)
            {
                if (player.Id == playerId)
                {
                    return player;
                }
            }

            return null;
        }

        /// <summary>
        /// 发送消息 
        /// </summary>
        /// <param name="player"></param>
        /// <param name="mid"></param>
        /// <param name="msg"></param>
        /// <param name="seq"></param>
        /// <returns></returns>
        public bool SendMsg(Player player, MID mid, IMessage msg, uint seq = 0)
        {
            if (player.GateClient == null)
            {
                player.GateClient = GalacticKittensNetworkManager.singleton.GetGateClient(player.GateUrl);
                if (player.GateClient == null)
                {
                    Log.Error($"{player.Id} message {mid} send fail: gate client {player.GateUrl} not find");
                    return false;
                }
            }

            return player.GateClient.SendMsg(player.Id, (int)mid, msg, seq);
        }

        /// <summary>
        /// 广播消息
        /// </summary>
        /// <param name="mid"></param>
        /// <param name="message"></param>
        /// <param name="seq"></param>
        /// <param name="excludePredicate"></param>
        /// <returns></returns>
        public void BroadcastMsg(MID mid, IMessage message, uint seq = 0, Predicate<long> excludePredicate = null)
        {
            foreach (var player in players)
            {
                if (excludePredicate != null && excludePredicate.Invoke(player.Id))
                {
                    continue;
                }

                SendMsg(player, mid, message, seq);
            }
        }


        /// <summary>
        /// 请求玩家列表，并初始化网络 TODO 待测试,会阻塞主线程吗？
        /// </summary>
        public void PlayerListReq()
        {
            var client =
                new GalacticKittensMatchService.GalacticKittensMatchServiceClient(GalacticKittensNetworkManager
                    .singleton.MatchChannel);
            var request = new GalacticKittensPlayerServerListRequest()
            {
                RoomId = GalacticKittensNetworkManager.singleton.ServerId
            };
            var response = client.playerServerListAsync(request).ResponseAsync.Result;
            Log.Info($"player list :{response}");
            if (response.Result != null && response.Result.Status != 200)
            {
                Log.Error($"get server list error:{response.Result.Msg}");
                return;
            }

            //网关
            Dictionary<uint, ServerInfo> gateServers = new Dictionary<uint, ServerInfo>(2);
            foreach (var info in response.PlayerGateServers)
            {
                var player = new Player();
                player.Id = info.Key;
                var serverInfo = info.Value;
                player.GateUrl = serverInfo.GrpcUrl;
                players.Add(player);
                gateServers[serverInfo.Id] = serverInfo;
            }

            GalacticKittensNetworkManager.singleton.ConnectToGate(gateServers);

            //大厅
            Dictionary<uint, ServerInfo> lobbyServers = new Dictionary<uint, ServerInfo>(2);
            foreach (var info in response.PlayerLobbyServers)
            {
                var serverInfo = info.Value;
                lobbyServers[serverInfo.Id] = serverInfo;
                var player = GetPlayer(info.Key);
                player.LobbyId = serverInfo.Id;
            }

            GalacticKittensNetworkManager.singleton.ConnectToLobby(lobbyServers);

            // 向大厅请求玩家基础信息 
            PlayerInfoReq();

            //  创建网关连接
            BindGateGameMapReq();

            // 切换场景 TODO 需要Unity开发快捷方式，待测试
            SceneManager.LoadScene("GalacticKittensGamePlay");

            RoomManager.Instance.SpawnPlayers(players);
        }

        private void PlayerInfoReq()
        {
            foreach (var player in players)
            {
                var channel = GalacticKittensNetworkManager.singleton.GetLobbyChannel(player.LobbyId);
                if (channel == null)
                {
                    Log.Error($"{player.Id} lobby {player.LobbyId} channel not exist");
                    return;
                }

                var client = new PlayerService.PlayerServiceClient(channel);
                var request = new PlayerInfoRequest()
                {
                    PlayerId = player.Id
                };

                var response = client.GetPlayerInfoAsync(request).ResponseAsync.Result;
                if (response.Result != null && response.Result.Status != 200)
                {
                    Log.Error($"{player.Id} get info error:{response.Result.Msg}");
                    if (response.Result.Msg.Equals("room no player"))
                    {
                        Log.Error("quit game");
                        Application.Quit();
                    }

                    return;
                }

                Log.Debug($"{player.Id} player info:{response}");
                var info = response.Player;
                player.Nick = info.Nick;
                player.Level = info.Level;
                player.Exp = info.Exp;
            }
        }

        /// <summary>
        /// 绑定玩家网关游戏映射
        /// </summary>
        public void BindGateGameMapReq(bool bind = true)
        {
            foreach (var player in players)
            {
                var request = new BindGameConnectRequest()
                {
                    PlayerId = player.Id,
                    Bind = bind
                };

                SendMsg(player, MID.BindGameConnectReq, request);
            }
        }
    }
}