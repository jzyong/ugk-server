using Common.Network.Sync;
using Game.Manager;
using UGK.Common.Network.Sync;
using ugk.Game.Room.Player;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 不发射子弹，只有碰撞的敌人
    /// </summary>
    public class SpaceGhostEnemy : BaseEnemyBehavior
    {
        protected override void Update()
        {
            ChangeVelocity();
        }


        private void OnTriggerEnter2D(Collider2D otherObject)
        {
            // check if it's collided with a player spaceship
            var spaceShip = otherObject.gameObject.GetComponent<SpaceShip>();
            if (spaceShip != null)
            {
                // tell the spaceship that it's taken damage
                spaceShip.Hit(1);

                RoomManager.Instance.DespawnObject(spaceShip.GetComponent<SnapTransform>().Id,
                    gameObject.GetComponent<SnapTransform>().Id);
                Destroy(gameObject);
            }
        }
    }
}