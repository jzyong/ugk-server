using System.Collections;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 敌人刷新
    /// </summary>
    public class EnemySpawner : MonoBehaviour
    {
        //TODO 只需要共用一个prefab
        public GameObject spaceGhostEnemyPrefabToSpawn;
        public GameObject spaceShooterEnemyPrefabToSpawn;


        [SerializeField]
        private GameObject m_bossPrefabToSpawn;

        [SerializeField]
        private Transform m_bossPosition;

        [Header("Enemies")]
        [SerializeField]
        private float m_EnemySpawnTime = 1.8f;

        [SerializeField]
        private float m_bossSpawnTime =75;

        [Header("Meteors")]
        [SerializeField]
        private GameObject m_meteorPrefab;

        [SerializeField]
        private float m_meteorSpawningTime;



        private Vector3 m_CurrentNewEnemyPosition = new Vector3();
        private float m_CurrentEnemySpawnTime = 0f;
        private Vector3 m_CurrentNewMeteorPosition = new Vector3();
        private float m_CurrentMeteorSpawnTime = 0f;
        private float m_CurrentBossSpawnTime = 0f;
        private bool m_IsSpawning = true;

        private void Start()
        {
            // Initialize the enemy and meteor spawn position based on my owning GO's x position
            m_CurrentNewEnemyPosition.x = transform.position.x;
            m_CurrentNewEnemyPosition.z = 0f;

            m_CurrentNewMeteorPosition.x = transform.position.x;
            m_CurrentNewMeteorPosition.z = 0f;
        }

        // Update is called once per frame
        void Update()
        {
            if (! m_IsSpawning)
                return;

            UpdateEnemySpawning();

            UpdateMeteorSpawning();

            UpdateBossSpawning();
        }

        private void UpdateEnemySpawning()
        {
            m_CurrentEnemySpawnTime += Time.deltaTime;
            if (m_CurrentEnemySpawnTime >= m_EnemySpawnTime)
            {
                // update the new enemy's spawn position(y value). This way we don't have to allocate
                // a new Vector3 each time.
                m_CurrentNewEnemyPosition.y = Random.Range(-5f, 5f);

                var nextPrefabToSpawn = GetNextRandomEnemyPrefabToSpawn();
                
                // //广播敌人刷新产出
                // NetworkObjectSpawner.SpawnNewNetworkObject(
                //     nextPrefabToSpawn,
                //     m_CurrentNewEnemyPosition);

                m_CurrentEnemySpawnTime = 0f;
            }
        }

        GameObject GetNextRandomEnemyPrefabToSpawn()
        {
            int randomPick = Random.Range(0, 99);

            if (randomPick < 50)
            {
                return spaceGhostEnemyPrefabToSpawn;
            }

            // randomPick >= 50
            return spaceShooterEnemyPrefabToSpawn;
        }

        private void UpdateMeteorSpawning()
        {
            m_CurrentMeteorSpawnTime += Time.deltaTime;
            if (m_CurrentMeteorSpawnTime > m_meteorSpawningTime)
            {
                SpawnNewMeteor();

                m_CurrentMeteorSpawnTime = 0f;
            }
        }

        void SpawnNewMeteor()
        {
            // The min and max Y pos for spawning the meteors
            m_CurrentNewMeteorPosition.y = Random.Range(-5f, 6f);

            // TODO 广播陨石产出
            // NetworkObjectSpawner.SpawnNewNetworkObject(m_meteorPrefab, m_CurrentNewMeteorPosition);
        }

        private void UpdateBossSpawning()
        {
            m_CurrentBossSpawnTime += Time.deltaTime;
            if (m_CurrentBossSpawnTime >= m_bossSpawnTime)
            {
                m_IsSpawning = false;
                StartCoroutine(BossAppear());
            }
        }

        IEnumerator BossAppear()
        {
            // Warning title and sound
            

            yield return new WaitForSeconds(3);

            //BOSS状态同步、刷出对象 TODO
            

            // GameObject boss = NetworkObjectSpawner.SpawnNewNetworkObject(
            //     m_bossPrefabToSpawn,
            //     transform.position);

            // Boss bossController = boss.GetComponent<Boss>();
            // bossController.StartBoss(m_bossPosition.position);
            // bossController.SetUI(m_bossUI);
            // boss.name = "BOSS";
        }
    }
}