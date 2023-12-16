using System.Collections;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public class BossMisileBarrageState : BaseBossState
    {
        [SerializeField]
        Transform[] m_misileSpawningArea;

        [SerializeField]
        GameObject m_misilePrefab;

        [SerializeField]
        [Range(0f, 1f)]
        float m_misileDelayBetweenSpawns;
    
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
           //TODO 产出导弹
        }

        // Run state
        public override void RunState()
        {
            StartCoroutine(RunMisileBarrageState());
        }
    }
}
