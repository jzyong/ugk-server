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

    [RequireComponent(typeof(BossController))]
    public class BaseBossState : MonoBehaviour
    {
        protected BossController m_controller;

        private void Start()
        {
            m_controller = FindFirstObjectByType<BossController>();
        }
    
        // Method that should be run on all states
        public virtual void RunState() { }
    
        public virtual void StopState() 
        {
            StopAllCoroutines();
        }
    }
}