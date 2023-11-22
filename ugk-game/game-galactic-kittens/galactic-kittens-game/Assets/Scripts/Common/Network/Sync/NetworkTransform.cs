using System.Collections;
using UnityEngine;

namespace Common.Network.Sync
{
    /// <summary>
    /// 同步
    /// </summary>
    public abstract class NetworkTransform : MonoBehaviour
    {
        // target transform to sync. can be on a child.
        [Header("Target")] [Tooltip("The Transform component to sync. May be on on this GameObject, or on a child.")]
        public Transform target;


        [Header("Selective Sync\nDon't change these at Runtime")]
        public bool syncPosition = true; // do not change at runtime!

        public bool syncRotation = true; // do not change at runtime!
        public bool syncScale = false; // do not change at runtime! rare. off by default.


        protected virtual void Awake()
        {
        }

        protected void OnValidate()
        {
            // set target to self if none yet
            if (target == null) target = transform;
        }


        protected virtual void OnEnable()
        {
        }

        protected virtual void OnDisable()
        {
        }
    }
}