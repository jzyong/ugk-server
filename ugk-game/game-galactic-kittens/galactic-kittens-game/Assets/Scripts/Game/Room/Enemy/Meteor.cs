using System.Collections;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 陨石
    /// </summary>
    public class Meteor : MonoBehaviour, IDamagable
    {
        [SerializeField]
        private int m_damage = 1;
    
        [SerializeField]
        private int m_health = 1;


        [Header("Range for random scale value")]
        [SerializeField]
        private float m_scaleMin = 0.8f;

        [SerializeField]
        private float m_scaleMax = 1.5f;

        private void Start()
        {

          
        }

        public void SpawnInit()
        {
            // Randomly scale the meteor
            float randomScale = Random.Range(m_scaleMin, m_scaleMax);
            transform.localScale = new Vector3(randomScale, randomScale, 1f);
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {

            if (collider.TryGetComponent(out IDamagable damagable))
            {
                // Hit the object that collide with me
                damagable.Hit(m_damage);

                // Hit me too!
                Hit(m_damage);
            }
        }


        public void Hit(int damage)
        {
           //TODO
        }
    }
}