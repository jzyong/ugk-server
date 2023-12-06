using Game.Room.Player;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 发射子弹的敌人
    /// </summary>
    public class SpaceShooterEnemyBehavior : BaseEnemyBehavior
    {
        [SerializeField]
        public GameObject m_EnemyBulletPrefab;

        [SerializeField]
        private float m_ShootingCooldown =2 ;


        private float m_CurrentCooldownTime = 0f;



        protected override void UpdateActive()
        {
            MoveEnemy();

            m_CurrentCooldownTime += Time.deltaTime;
            if (m_CurrentCooldownTime >= m_ShootingCooldown)
            {
                m_CurrentCooldownTime = 0f;
                ShootLaserServerRpc();
            }
        }
        

        private void ShootLaserServerRpc()
        {
            // TODO 发射子弹
            // var newEnemyLaser = NetworkObjectSpawner.SpawnNewNetworkObject(m_EnemyBulletPrefab);
            //
            // var bulletController = newEnemyLaser.GetComponent<BulletController>();
            // if (bulletController != null)
            // {
            //     bulletController.m_Owner = gameObject;
            // }
            //
            // newEnemyLaser.transform.position = this.gameObject.transform.position;
        }


        private void OnTriggerEnter2D(Collider2D otherObject)
        {

            // check if it's collided with a player spaceship
            var spaceShip = otherObject.gameObject.GetComponent<SpaceShip>();
            if (spaceShip != null)
            {
                // tell the spaceship that it's taken damage
                spaceShip.Hit(1);

                // enemy explodes when it collides with the a player's ship
                m_EnemyState = EnemyState.defeatAnimation;
            }

            // check if it's collided with a player's bullet
            var shipBulletBehavior = otherObject.gameObject.GetComponent<BulletController>();
            if (shipBulletBehavior != null && shipBulletBehavior.m_Owner != this.gameObject)
            {
                // if so, take one health point away from enemy
                m_EnemyHealthPoints -= 1;
            }
        }

        private void OnEnemyHealthPointsChange(int oldHP, int newHP)
        {
            // if enemy's health is 0, then time to start enemy dead animation
            if (newHP <= 0)
            {
                m_EnemyState = EnemyState.defeatAnimation;
            }
        }
    }
}