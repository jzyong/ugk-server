using Common.Network.Sync;
using Game.Manager;
using UGK.Common.Network.Sync;
using UGK.Game.Manager;
using UnityEngine;

namespace Game.Room.Boss
{
    /// <summary>
    /// 环形子弹
    /// </summary>
    public class BossCircularBullet : MonoBehaviour
    {
        [SerializeField] private Transform[] _firePositions;
        [Tooltip("移动方向")] public Vector3 direction = Vector3.left;

        [SerializeField] [Tooltip("速度")] private float speed = 2;


        private void Update()
        {
            transform.Translate(speed * Time.deltaTime * direction);
        }

        public void Start()
        {
            // After a random time amount, blow up and spawn small bullets
            float randomSpawn = Random.Range(1.5f, 3f);
            Invoke(nameof(SpawnBullets), randomSpawn);
        }

        private void SpawnBullets()
        {
            // Spawn the bullets
            foreach (Transform firePosition in _firePositions)
            {
                RoomManager.Instance.SpawnBossBullet(33, firePosition.position, firePosition.rotation.eulerAngles);
            }

            // De-spawn me
            RoomManager.Instance.DespawnObject(0, GetComponent<SnapTransform>().Id, true);
        }
    }
}