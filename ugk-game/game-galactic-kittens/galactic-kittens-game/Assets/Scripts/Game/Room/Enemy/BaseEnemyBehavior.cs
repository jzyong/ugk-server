using System.Collections;
using Common.Network.Sync;
using Game.Manager;
using UGK.Common.Network.Sync;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 敌人基础类
    /// </summary>
    public class BaseEnemyBehavior : MonoBehaviour, IDamagable
    {
        protected enum EnemyMovementType
        {
            linear,
            sineWave,

            // you can add more movement types here

            COUNT //MAX - used to get random value
        }


        [SerializeField] protected float m_EnemySpeed = 4f;


        [SerializeField] protected bool m_UsesEnemyLifetime = true;

        // TODO 血量减少需要同步广播
        [SerializeField] protected int m_EnemyHealthPoints = 3;



        protected EnemyMovementType m_EnemyMovementType;

        protected Vector2 m_Direction = Vector2.left;

        protected float m_WaveAmplitude;

        [SerializeField] private float m_hitEffectDuration;


        public void Start()
        {
            m_WaveAmplitude = Random.Range(2f, 6f);
            m_EnemyMovementType = GetRandomEnemyMovementType();
        }

        protected virtual void Update()
        {
        }


        /// <summary>
        /// 改变速度
        /// </summary>
        protected virtual void ChangeVelocity()
        {
            if (m_EnemyMovementType == EnemyMovementType.sineWave)
            {
                m_Direction.x = -1f; //to move from right to left
                m_Direction.y = Mathf.Sin(Time.time * m_WaveAmplitude);
                m_Direction.Normalize();
            }


            // move the enemy in the desired direction  使用预测同步效果不理想，因此使用快照插值同步
            transform.Translate(m_Direction * m_EnemySpeed * Time.deltaTime);
        }

        /// <summary>
        /// 随机移动类型
        /// </summary>
        /// <returns></returns>
        protected EnemyMovementType GetRandomEnemyMovementType()
        {
            int randomValue = Random.Range(0, (int)EnemyMovementType.COUNT);

            return (EnemyMovementType)randomValue;
        }

        protected void DespawnEnemy()
        {
            gameObject.SetActive(false);
        }


        public virtual void Hit(int damage)
        {
            m_EnemyHealthPoints -= 1;
            if (m_EnemyHealthPoints < 1)
            {
                RoomManager.Instance.DespawnObject(0, GetComponent<SnapTransform>().Id);
                RoomManager.Instance.SpawnPowerUp(gameObject.transform.position);
                Destroy(gameObject);
               
            }
        }
    }
}