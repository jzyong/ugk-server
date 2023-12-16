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

        private void Start()
        {
            M = FindFirstObjectByType<Boss>();
        }
    
        // Method that should be run on all states
        public virtual void RunState() { }
    
        public virtual void StopState() 
        {
            StopAllCoroutines();
        }
    }
}