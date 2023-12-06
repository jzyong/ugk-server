using Game.Room.Player;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 不发射子弹，只有碰撞的敌人
    /// </summary>
    public class SpaceGhostEnemyBehavior : BaseEnemyBehavior
    {

        private bool m_IsFlashingFromHit = false;
        private float m_FlashFromHitTime = 0.7f;




        protected override void UpdateActive()
        {
            if (m_IsFlashingFromHit)
            {
                m_FlashFromHitTime -= Time.deltaTime;
                if (m_FlashFromHitTime <= 0f)
                {
                }
            }

            MoveEnemy();
        }

        protected override void UpdateDefeatedAnimation()
        {
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
