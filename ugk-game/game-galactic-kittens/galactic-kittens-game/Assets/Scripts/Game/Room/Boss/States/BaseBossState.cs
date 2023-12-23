using Common.Tools;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public enum BossState
    { 
        fire,
        misileBarrage,
        death,
        idle,
        enter
    };

    [RequireComponent(typeof(ugk.Game.Room.Boss.Boss))]
    public class BaseBossState : MonoBehaviour
    {
        protected ugk.Game.Room.Boss.Boss M;

        protected void Awake()
        {
            M = FindFirstObjectByType<ugk.Game.Room.Boss.Boss>();
            if (M==null)
            {
                Log.Warn("未找到Boss脚本");
            }
        }
    
        // Method that should be run on all states
        public virtual void RunState() { }
    
        public virtual void StopState() 
        {
            StopAllCoroutines();
        }
    }
}