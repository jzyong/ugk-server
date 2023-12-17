using System;
using Common.Network.Sync;
using Game.Manager;
using Game.Room.Boss.States;
using Game.Room.Player;
using UnityEngine;

/*
    Script that controls how the boss is going to work,
    the different behaviours are set on different scripts. 
    Here you can add new states
*/

namespace Game.Room.Boss
{
    public class Boss : MonoBehaviour, IDamagable
    {
        [SerializeField] private int m_damage;

        [Header("States for the boss")] [SerializeField]
        private BossEnterState m_enterState;

        [SerializeField] private BaseBossState m_fireState;

        [SerializeField] private BaseBossState m_misileBarrageState;

        [SerializeField] private BaseBossState m_idleState;

        [SerializeField] private BaseBossState m_deathState;

        [Header("For testing the boss states -> false for production")] [SerializeField]
        private bool m_isTesting;

        [SerializeField] private BossState m_testState;

        [SerializeField] [Tooltip("血量")] private float helath = 15;


        private void Awake()
        {
            helath = helath * RoomManager.Instance.PlayerCount();
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {
            // When the players get close to me do some damage
            if (collider.TryGetComponent(out SpaceShip playerShip))
            {
                playerShip.Hit(m_damage);
            }
        }


        // This will set the starting state for the boss -> enter state
        public void StartBoss(Vector3 initialPositionForEnterState)
        {
            m_enterState.initialPosition = initialPositionForEnterState;
            SetState(BossState.enter);
        }

        // Set the boss state to run
        // You can add more states to the boss
        //..
        public void SetState(BossState state)
        {
            switch (state)
            {
                case BossState.enter:
                    m_enterState.RunState();
                    break;
                case BossState.fire:
                    m_fireState.RunState();
                    break;
                case BossState.misileBarrage:
                    m_misileBarrageState.RunState();
                    break;
                case BossState.idle:
                    m_idleState.RunState();
                    break;
                case BossState.death:
                    // Stop all coroutines from other state
                    // because the death can override any state
                    m_enterState.StopState();
                    m_fireState.StopState();
                    m_misileBarrageState.StopState();
                    m_idleState.StopState();

                    m_deathState.RunState();
                    break;
            }
        }

        // // Set the boss UI
        // public void SetUI(BossUI bossUI)
        // {
        //
        //     BossHealth bossHealth = GetComponentInChildren<BossHealth>();
        //     this.bossUI = bossUI;
        //     bossUI.SetHealth(bossHealth.Health);
        // }

        // public override void OnNetworkSpawn()
        // {
        //     if (IsServer && m_isTesting)
        //     {
        //         // If you want to test the boss outside of the normal flow of the game
        //         SetState(m_testState);
        //     }
        //     base.OnNetworkSpawn();
        // }
        public void Hit(int damage)
        {
            helath -= damage;
            if (helath < 1)
            {
                SetState(BossState.death);
                RoomManager.Instance.DespawnObject(0, GetComponent<SnapTransform>().Id);
                Destroy(gameObject);
            }
        }
    }
}