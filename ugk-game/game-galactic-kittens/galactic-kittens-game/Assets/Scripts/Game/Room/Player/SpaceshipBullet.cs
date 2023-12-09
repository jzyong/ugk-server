using System;
using Game.Manager;
using Network.Sync;
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
             _predictionTransform.GetComponent<PredictionTransform>();
             _predictionTransform.AngularVelocity = linearVelocity;
         }

         private void OnTriggerEnter2D(Collider2D collider)
         {
             if (collider.TryGetComponent(out IDamagable damagable))
             {
                 damagable.Hit(damage);

                 //  广播子弹消失
                 RoomManager.Instance.DespawnObject(0, _predictionTransform.Id);
             }
         }
    }
}