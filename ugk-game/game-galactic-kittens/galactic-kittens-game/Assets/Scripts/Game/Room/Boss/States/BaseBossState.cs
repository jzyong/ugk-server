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

    [RequireComponent(typeof(Boss))]
    public class BaseBossState : MonoBehaviour
    {
        protected Boss M;

        protected void Awake()
        {
            M = FindFirstObjectByType<Boss>();
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