using Common.Network.Sync;
using Game.Manager;
using Game.Room.Enemy;
using UnityEngine;

namespace Game.Room.Player
{
    /// <summary>
    /// 玩家飞船
    /// </summary>
    public class SpaceShip : MonoBehaviour, IDamagable
    {
        private int powerUpCount;

        void OnTriggerEnter2D(Collider2D collider)
        {
            // If the collider hit a power-up
            if (collider.TryGetComponent(out PowerUp powerUp))
            {
                // Check if I have space to take the special
                if (powerUpCount < 2)
                {
                    // Update var
                    powerUpCount++;

                    // Update UI TODO
                    // playerUI.UpdatePowerUp(m_specials.Value, true);

                    // Remove the power-up
                    RoomManager.Instance.DespawnObject(0, collider.gameObject.GetComponent<PredictionTransform>().Id);
                }
            }
        }

        public void Hit(int damage)
        {
            //TODO
        }
    }
}