using Common.Network.Sync;
using Game.Manager;
using Game.Room.Player;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 发射子弹的敌人
    /// </summary>
    public class SpaceShooterEnemy : BaseEnemyBehavior
    {

        [SerializeField] private float m_ShootingCooldown = 2;


        private float m_CurrentCooldownTime = 0f;


        protected override void Update()
        {
            ChangeVelocity();
            //发射子弹
            m_CurrentCooldownTime += Time.deltaTime;
            if (m_CurrentCooldownTime >= m_ShootingCooldown)
            {
                m_CurrentCooldownTime = 0f;
                RoomManager.Instance.SpawnEnemyBullet(this);
            }
        }

        private void OnTriggerEnter2D(Collider2D otherObject)
        {
            // check if it's collided with a player spaceship
            var spaceShip = otherObject.gameObject.GetComponent<SpaceShip>();
            if (spaceShip != null)
            {
                // tell the spaceship that it's taken damage
                spaceShip.Hit(1);
                RoomManager.Instance.DespawnObject(spaceShip.GetComponent<SnapTransform>().Id, gameObject.GetComponent<SnapTransform>().Id);
                Destroy(gameObject);
            }
        }
    }
}