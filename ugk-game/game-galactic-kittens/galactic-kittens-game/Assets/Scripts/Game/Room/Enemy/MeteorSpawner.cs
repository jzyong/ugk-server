using UnityEngine;

namespace Game.Room.Enemy
{
    public class MeteorSpawner : MonoBehaviour
    {
        [SerializeField]
        private GameObject m_meteorPrefab;

        [SerializeField]
        private float m_spawingTime = 2f;

        private float m_timer;

        private void Update()
        {

            m_timer += Time.deltaTime;
            if (m_timer > m_spawingTime)
            {
                m_timer = 0f;
                SpawnMeteor();
            }
        }

        private void SpawnMeteor()
        {
            // The min and max Y pos for spawning the meteors
            float randomYpos = Random.Range(-5f, 6f);
            var newMeteorPosition = new Vector3(transform.position.x, randomYpos, 0f);
            //TODO 广播陨石产出
           
        }

    }
}
