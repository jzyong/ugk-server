using Common.Network.Sync;
using Game.Manager;
using Game.Room;
using Game.Room.Enemy;
using UGK.Common.Network.Sync;
using UnityEngine;

namespace ugk.Game.Room.Player
{
    /// <summary>
    /// 玩家飞船
    /// </summary>
    public class SpaceShip : MonoBehaviour, IDamagable
    {
        [HideInInspector] public uint powerUpCount;

        [Tooltip("血量")] public uint hp = 30;


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

                    // Update UI 
                    RoomManager.Instance.BroadcastPlayerProperty(this);

                    // Remove the power-up
                    RoomManager.Instance.DespawnObject(0, collider.gameObject.GetComponent<SnapTransform>().Id);
                }
            }
        }

        public void Hit(int damage)
        {
            hp--;
            if (hp < 1)
            {
                gameObject.SetActive(false);
                RoomManager.Instance.DespawnObject(0, GetComponent<SnapTransform>().Id, false);
                RoomManager.Instance.GameFinishFail();
            }

            RoomManager.Instance.BroadcastPlayerProperty(this);
        }
    }
}