using System;
using System.Collections.Generic;
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
        /// <summary>
        /// 选择的角色id
        /// </summary>
        public Int32 CharacterId { get; set; }
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
                player.GateClient = GalacticKittensNetworkManager.Instance.GetGateClient(player.GateUrl);
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
        /// 请求玩家列表，并初始化网络
        /// <param name="roomId">0 编辑器测试模式，匹配服分配，其他需要正确的id</param>
        /// </summary>
        public void PlayerListReq(uint roomId=1)
        {
            var client =
                new GalacticKittensMatchService.GalacticKittensMatchServiceClient(GalacticKittensNetworkManager
                    .Instance.MatchChannel);
            var request = new GalacticKittensPlayerServerListRequest()
            {
                RoomId =roomId==0?roomId: GalacticKittensNetworkManager.Instance.ServerId
            };
            
            
            var response = client.playerServerListAsync(request).ResponseAsync.Result;
            Log.Info($"player list :{response}");
            if (response.Result != null && response.Result.Status != 200)
            {
                Log.Error($"get server list error:{response.Result.Msg}");
                return;
            }

            GalacticKittensNetworkManager.Instance.ServerId = response.RoomId;

            //网关
            Dictionary<uint, ServerInfo> gateServers = new Dictionary<uint, ServerInfo>(2);
            foreach (var info in response.PlayerGateServers)
            {
                var player = new Player
                {
                    Id = info.Key
                };
                var serverInfo = info.Value;
                player.GateUrl = serverInfo.GrpcUrl;
                var playerInfo = response.PlayerInfos[player.Id];
                player.CharacterId = playerInfo.CharacterId;
                players.Add(player);
                gateServers[serverInfo.Id] = serverInfo;
            }

            GalacticKittensNetworkManager.Instance.ConnectToGate(gateServers);

            //大厅 不用连接大厅，通过match中转连接大厅可能更好，集中处理相关逻辑
            Dictionary<uint, ServerInfo> lobbyServers = new Dictionary<uint, ServerInfo>(2);
            foreach (var info in response.PlayerLobbyServers)
            {
                var serverInfo = info.Value;
                lobbyServers[serverInfo.Id] = serverInfo;
                var player = GetPlayer(info.Key);
                player.LobbyId = serverInfo.Id;
            }

            GalacticKittensNetworkManager.Instance.ConnectToLobby(lobbyServers);

            // 向大厅请求玩家基础信息 
            PlayerInfoReq();

            //  创建网关连接
            BindGateGameMapReq();

            // 切换场景 
            SceneManager.LoadScene("GalacticKittensGamePlay");

            RoomManager.Instance.SpawnPlayers(players);
        }

        private void PlayerInfoReq()
        {
            foreach (var player in players)
            {
                var channel = GalacticKittensNetworkManager.Instance.GetLobbyChannel(player.LobbyId);
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