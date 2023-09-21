using System;
using Common.Network;
using Common.Tools;
using UnityEngine;

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
        public static PlayerManager Singleton { get; internal set; }


        public Player GetPlayer(Int64 playerId)
        {
            //TODO
            return null;
        }

    }
}