using System;
using Common.Network.Sync;
using Game.Manager;
using UGK.Common.Network.Sync;
using UGK.Game.Manager;
using ugk.Game.Room.Player;
using UnityEngine;

namespace Game.Room.Boss
{
    /// <summary>
    /// boss 普通小子弹
    /// </summary>
    public class BossSmallBullet : MonoBehaviour
    {
        private int m_damage = 1;
        [Tooltip("移动方向")]
        public Vector3 direction = Vector3.up;

        [SerializeField][Tooltip("速度")]
        private float speed;
        
        
        
        private void Update()
        {
            transform.Translate(speed * Time.deltaTime * direction);    
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {
            if (collider.TryGetComponent(out SpaceShip spaceShip))
            {
                spaceShip.Hit(m_damage);
                long killerId = 0;

                var snapTransform = collider.GetComponent<SnapTransform>();
                if (snapTransform != null)
                {
                    killerId = snapTransform.Id;
                }

                // 移除对象
                RoomManager.Instance.DespawnObject(killerId,GetComponent<SnapTransform>().Id);
            }
        }
    }
}
