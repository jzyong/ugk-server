using UnityEngine;

namespace Game.Room.Boss
{
    public class SmallBullet : MonoBehaviour
    {
        private int m_damage = 1;

        private void OnTriggerEnter2D(Collider2D collider)
        {
            //TODO
            // if (IsServer && collider.TryGetComponent(out IDamagable damagable))
            // {
            //     damagable.Hit(m_damage);
            //
            //     NetworkObjectDespawner.DespawnNetworkObject(NetworkObject);
            // }
        }
    }
}
