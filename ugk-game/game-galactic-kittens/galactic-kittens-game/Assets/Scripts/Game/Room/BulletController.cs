using System;
using Game.Manager;
using Network.Sync;
using UnityEngine;

namespace Game.Room
{
    /// <summary>
    /// 子弹控制器 
    /// </summary>
    [Obsolete]
    public class BulletController : MonoBehaviour
    {
        private enum BulletOwner
        {
            enemy,
            player
        };

        [SerializeField]
        private BulletOwner m_owner;
        [HideInInspector] public GameObject m_Owner { get; set; } = null;

        public int damage = 1;

        private void OnTriggerEnter2D(Collider2D collider)
        {
            if (collider.TryGetComponent(out IDamagable damagable))
            {
                damagable.Hit(damage);
                if (m_owner==BulletOwner.player)
                {
                    // 增加玩家命中计数 TODO
                }

                //  广播子弹消失
                RoomManager.Instance.DespawnObject(0, GetComponent<PredictionTransform>().Id);
            }
        }
    }
}