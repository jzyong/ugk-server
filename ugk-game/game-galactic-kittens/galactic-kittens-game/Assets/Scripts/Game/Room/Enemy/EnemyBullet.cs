using Common.Network.Sync;
using kcp2k;
using UGK.Common.Network.Sync;
using UGK.Game.Manager;
using ugk.Game.Room.Player;
using UnityEngine;

namespace UGK.Game.Room.Enemy
{
    /// <summary>
    /// 飞船子弹
    /// </summary>
    public class EnemyBullet : MonoBehaviour
    {
        [Tooltip("线速度")] public Vector3 linearVelocity = Vector3.left * 8;

        public int damage = 1;
        private PredictionTransform _predictionTransform;

        private void Awake()
        {
            _predictionTransform = transform.GetComponent<PredictionTransform>();
            _predictionTransform.LinearVelocity = linearVelocity;
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {
            if (collider.TryGetComponent(out SpaceShip spaceShip))
            {
                Log.Info($"命中敌人{collider.name}");
                spaceShip.Hit(damage);
                long killerId = 0;

                var snapTransform = collider.GetComponent<SnapTransform>();
                if (snapTransform != null)
                {
                    killerId = snapTransform.Id;
                }
                //  广播子弹消失
                RoomManager.Instance.DespawnObject(killerId, _predictionTransform.Id);
            }
            
        }
    }
}