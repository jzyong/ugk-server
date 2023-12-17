using System.Collections;
using System.Collections.Generic;
using Game.Manager;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public class BossDeathState : BaseBossState
    {
        IEnumerator RunDeath()
        {
            RoomManager.Instance.GameFinish();
            yield return new WaitForSeconds(10);
            
            
            //关闭docker
            Application.Quit();
        }

        public override void RunState()
        {
            StartCoroutine(RunDeath());
        }
    }
}