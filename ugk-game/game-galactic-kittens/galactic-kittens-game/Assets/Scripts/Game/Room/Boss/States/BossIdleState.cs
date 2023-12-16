using System.Collections;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public class BossIdleState : BaseBossState
    {
        [SerializeField]
        [Range(0.1f, 2f)]
        float m_idleTime;

        IEnumerator RunIdleState()
        {
            // Wait for a moment
            yield return new WaitForSeconds(m_idleTime);

            // Call the fire state
            M.SetState(BossState.fire);
        }

        public override void RunState()
        {
            StartCoroutine(RunIdleState());
        }
    }
}
