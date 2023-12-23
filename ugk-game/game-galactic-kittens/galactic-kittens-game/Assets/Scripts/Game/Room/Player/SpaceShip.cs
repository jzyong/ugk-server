using Game.Room;
using Game.Room.Enemy;
using UGK.Common.Network.Sync;
using UGK.Game.Manager;
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


        public uint UsePopwerCount { get; set; }

        public uint KillEnemyCount { get; set; }


        void OnTriggerEnter2D(Collider2D collider2D)
        {
            // If the collider2D hit a power-up
            if (collider2D.TryGetComponent(out PowerUp powerUp))
            {
                // Check if I have space to take the special
                if (powerUpCount < 2)
                {
                    // Update var
                    powerUpCount++;

                    // Update UI 
                    RoomManager.Instance.BroadcastPlayerProperty(this);

                    // Remove the power-up
                    RoomManager.Instance.DespawnObject(0, collider2D.gameObject.GetComponent<SnapTransform>().Id);
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