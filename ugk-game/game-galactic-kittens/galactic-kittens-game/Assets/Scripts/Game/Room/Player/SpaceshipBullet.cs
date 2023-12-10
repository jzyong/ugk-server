using System;
using Common.Network.Sync;
using Game.Manager;
using Game.Room.Enemy;
using kcp2k;
using UnityEngine;

namespace Game.Room.Player
{
    /// <summary>
    /// 飞船子弹
    /// </summary>
    public class SpaceshipBullet : MonoBehaviour
    {
        [Tooltip("线速度")] public Vector3 linearVelocity = Vector3.right * 2;

        public int damage = 1;
        private PredictionTransform _predictionTransform;

        private void Awake()
        {
            _predictionTransform = transform.GetComponent<PredictionTransform>();
            _predictionTransform.AngularVelocity = linearVelocity;
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
            }
        }
    }
}