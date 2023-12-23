using Common.Network.Sync;
using Game.Manager;
using Game.Room;
using kcp2k;
using UGK.Common.Network.Sync;
using UGK.Game.Manager;
using ugk.Game.Room.Player;
using UnityEngine;

namespace UGK.Game.Room.Player
{
    /// <summary>
    /// 飞船子弹
    /// </summary>
    public class SpaceshipBullet : MonoBehaviour
    {
        [Tooltip("线速度")] public Vector3 linearVelocity = Vector3.right * 2;

        public int damage = 1;
        private PredictionTransform _predictionTransform;

        public long OwnerId { get; set; }

        private void Awake()
        {
            _predictionTransform = transform.GetComponent<PredictionTransform>();
            _predictionTransform.LinearVelocity = linearVelocity;
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {
            if (collider.TryGetComponent(out IDamagable damagable))
            {
                if (damagable is SpaceShip)
                {
                    return;
                }

                Log.Info($"命中敌人{collider.name}");
                damagable.Hit(damage);
                long killerId = 0;

                var snapTransform = collider.GetComponent<SnapTransform>();
                if (snapTransform != null)
                {
                    killerId = snapTransform.Id;
                }

                //  广播子弹消失
                RoomManager.Instance.DespawnObject(killerId, _predictionTransform.Id);
                var spaceShip = RoomManager.Instance.GetSpaceShip(OwnerId);
                if (spaceShip != null)
                {
                    spaceShip.KillEnemyCount++;
                }
            }
        }
    }
}