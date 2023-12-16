using System.Collections;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public class BossSuperLaserState : BaseBossState
    {
        [SerializeField]
        GameObject m_superLaserPrefab;

        [SerializeField]
        Transform m_superLaserPosition;
    
        IEnumerator FireSuperLaser()
        {
            float randomRotation = Random.Range(-40f, 10f);

           //TODO 产出 Laser

            // TODO: Wait the time the vfx last
            yield return new WaitForSeconds(5f);
            M.SetState(BossState.idle);
        }
        
        public override void RunState()
        {
            StartCoroutine(FireSuperLaser());
        }

        public override void StopState()
        {
            StopCoroutine(FireSuperLaser());
        }  

    }
}
