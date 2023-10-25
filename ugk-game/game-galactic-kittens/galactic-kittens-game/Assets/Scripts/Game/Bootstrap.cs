using System;
using Game.Manager;
using UnityEngine;

namespace Game
{
    /// <summary>
    /// 启动入口
    /// </summary>
    public class Bootstrap : MonoBehaviour
    {
        private void Awake()
        {
            DontDestroyOnLoad(gameObject);
        }

        private void Start()
        {
            // 初始化
            PlayerManager.Singleton.PlayerListReq();
        }
    }
}