using System.Collections;
using Game.Manager;
using UnityEngine;

namespace Game.Room.Boss.States
{
    public class BossFireState : BaseBossState
    {
        [SerializeField] private Transform[] _fireCannonSpawningArea;

        [SerializeField] private float _normalShootRateOfFire;

        [SerializeField] private float _idleSpeed;

        [SerializeField] [Tooltip("子弹选择方向")] private Vector3[] trangleBulletRotation;


        public override void RunState()
        {
            StartCoroutine(FireState());
        }

        private IEnumerator FireState()
        {
            // Setup initial vars
            float shootTimer = 0f;
            float normalStateTimer = 0f;
            float normalStateExitTime = Random.Range(7f, 21f);

            // We have a random time on this state so while we are in this state we proceed to fire
            while (normalStateTimer <= normalStateExitTime)
            {
                // Small movement on the boss
                transform.position = new Vector2(
                    transform.position.x,
                    Mathf.Sin(Time.time) * _idleSpeed);

                // every x time shoot (trio or circular)
                shootTimer += Time.deltaTime;
                if (shootTimer >= _normalShootRateOfFire)
                {
                    var nextBulletPrefabToShoot = GetNextBulletPrefabToShoot();

                    FireBulletPrefab(nextBulletPrefabToShoot);

                    shootTimer = 0f;
                }

                yield return new WaitForEndOfFrame();

                normalStateTimer += Time.deltaTime;
            }

            // When we end the time on this state call the special attack, it can be a different state
            // or a random for different states
            M.SetState(BossState.misileBarrage);
        }

        /// <summary>
        /// 返回子弹类型
        /// </summary>
        /// <returns>32 boss三角形小子弹，33 boss环形分裂后小子弹，34 boss环形分裂子弹，35 boss导弹</returns>
        private uint GetNextBulletPrefabToShoot()
        {
            int randomBulletChoice = Random.Range(0, 10);

            // trio -> 7/10, circular -> 3/10
            if (randomBulletChoice < 7)
            {
                return 32;
            }

            return 34;
        }

        private void FireBulletPrefab(uint type)
        {
            // Because the cannon positions are lower on the sprite with increase the rotation up
            float randomZrotation = Random.Range(-25f, 45f);

            foreach (Transform laserCannon in _fireCannonSpawningArea)
            {
                if (type == 32)
                {
                    foreach (var t in trangleBulletRotation)
                    {
                        RoomManager.Instance.SpawnBossBullet(type, laserCannon.position,
                            new Vector3(t.x, t.y, t.z + randomZrotation));
                    }
                }
                else
                {
                    RoomManager.Instance.SpawnBossBullet(type, laserCannon.position,
                        new Vector3(0, 0, randomZrotation));
                }
            }
        }
    }
}