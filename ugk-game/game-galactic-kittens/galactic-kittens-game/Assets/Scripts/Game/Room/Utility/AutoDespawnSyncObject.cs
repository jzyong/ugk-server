using System;
using Game.Manager;
using Network.Sync;
using UnityEngine;

namespace Game.Room.Utility
{
    /// <summary>
    /// 自动销毁同步对象，用于移动的敌人离开视野
    /// </summary>
    [RequireComponent(typeof(PredictionTransform))]
    public class AutoDespawnOnServer : MonoBehaviour
    {
        [Min(0f)] [SerializeField] [Header("Time alive in seconds (s)")]
        private float m_autoDestroyTime;

        [SerializeField] [Tooltip("销毁的对象")] private GameObject target;


        private void OnValidate()
        {
            if (target == null)
            {
                target = transform.gameObject;
            }
        }

        private void Update()
        {
            m_autoDestroyTime -= Time.deltaTime;

            if (m_autoDestroyTime <= 0f)
            {
                // 广播对象死亡，销毁GameObject和移动组件
                RoomManager.Instance.DespawnObject(0, target.GetComponent<PredictionTransform>().Id);
                DestroyImmediate(gameObject);
            }
        }
    }
}