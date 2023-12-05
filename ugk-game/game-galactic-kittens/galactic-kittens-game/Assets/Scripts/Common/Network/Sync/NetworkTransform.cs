using Common.Network.Serialize;
using Common.Tools;
using UnityEngine;

namespace Common.Network.Sync
{
    /// <summary>
    /// Transform 同步  TODO 添加网络同步统计信息
    /// </summary>
    public abstract class NetworkTransform : MonoBehaviour
    {
        // 作用同步对象
        [Header("Target")] [Tooltip("同步作用的对象")]
        public Transform target;

        // selective sync 
        [Header("Selective Sync\nDon't change these at Runtime")]
        public bool syncPosition = true; // do not change at runtime!
        public bool syncRotation = false; // do not change at runtime!
        public bool syncScale = false; // do not change at runtime! rare. off by default.

        [Tooltip("消息发送间隔")] public double sendInterval = 0.033;
        

        [Tooltip(
            "Apply smallest-three quaternion compression. This is lossy, you can disable it if the small rotation inaccuracies are noticeable in your project.")]
        public bool compressRotation = false;

        // delta compression is capable of detecting byte-level changes.
        // if we scale float position to bytes,
        // then small movements will only change one byte.
        // this gives optimal bandwidth.
        //   benchmark with 0.01 precision: 130 KB/s => 60 KB/s
        //   benchmark with 0.1  precision: 130 KB/s => 30 KB/s
        [Header("Precision")]
        [Tooltip(
            "Position is rounded in order to drastically minimize bandwidth.\n\nFor example, a precision of 0.01 rounds to a centimeter. In other words, sub-centimeter movements aren't synced until they eventually exceeded an actual centimeter.\n\nDepending on how important the object is, a precision of 0.01-0.10 (1-10 cm) is recommended.\n\nFor example, even a 1cm precision combined with delta compression cuts the Benchmark demo's bandwidth in half, compared to sending every tiny change.")]
        [Range(0.00_01f, 1f)]
        // disallow 0 division. 1mm to 1m precision is enough range.
        public float positionPrecision = 0.01f; // 1 cm

        [Range(0.00_01f, 1f)] // disallow 0 division. 1mm to 1m precision is enough range.
        public float scalePrecision = 0.01f; // 1 cm

        /**
         * 每个对象的唯一id
         */
        public long Id { get; set; }
        //下次消息发送时间
        protected double nextSendTime;
        // delta compression needs to remember 'last' to compress against
        protected Vector3Long lastSerializedPosition = Vector3Long.zero;
        protected Vector3Long lastDeserializedPosition = Vector3Long.zero;

        protected Vector3Long lastSerializedScale = Vector3Long.zero;
        protected Vector3Long lastDeserializedScale = Vector3Long.zero;
        
        /// <summary>
        /// 是否为本地玩家拥有者
        /// </summary>
        public bool Onwer { get; set; }


        // make sure to call this when inheriting too!
        protected virtual void Awake()
        {
        }

        protected void OnValidate()
        {
            // set target to self if none yet
            if (target == null) target = transform;
        }
        
        /// <summary>
        /// 设置最后一次反序列化缓存的坐标，增量压缩还原需要
        /// </summary>
        /// <param name="position"></param>
        public void SetLastDeserializedPositon(Vector3 position)
        {
            Compression.ScaleToLong(position, positionPrecision, out lastDeserializedPosition);
        }
        
        /// <summary>
        /// 设置最后一次反序列化缓存的缩放，增量压缩还原需要
        /// </summary>
        /// <param name="scale"></param>
        public void SetLastDeserializedScale(Vector3 scale)
        {
            Compression.ScaleToLong(scale, scalePrecision, out lastDeserializedScale);
        }


    }
}