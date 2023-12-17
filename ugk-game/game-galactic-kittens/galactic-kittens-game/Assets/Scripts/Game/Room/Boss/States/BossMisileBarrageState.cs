using System.Collections;
using Game.Manager;
using UnityEngine;

namespace Game.Room.Boss.States
{
    /// <summary>
    /// 跟踪导弹 
    /// </summary>
    public class BossMisileBarrageState : BaseBossState
    {
        [SerializeField] Transform[] m_misileSpawningArea;

        [SerializeField] [Range(0f, 1f)] float m_misileDelayBetweenSpawns;

        IEnumerator RunMisileBarrageState()
        {
            // Spawn the missiles
            foreach (Transform spawnPosition in m_misileSpawningArea)
            {
                FireMisiles(spawnPosition.position);
                yield return new WaitForSeconds(m_misileDelayBetweenSpawns);
            }

            // Go idle from a moment
            M.SetState(BossState.idle);
        }

        // Spawn the missile prefab
        void FireMisiles(Vector3 position)
        {
            // 产出导弹
            RoomManager.Instance.SpawnBossBullet(35, position,Vector3.zero);
        }

        // Run state
        public override void RunState()
        {
            StartCoroutine(RunMisileBarrageState());
        }
    }
}