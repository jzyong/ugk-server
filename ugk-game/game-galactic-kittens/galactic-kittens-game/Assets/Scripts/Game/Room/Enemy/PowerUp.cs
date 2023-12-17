using Common.Network.Sync;
using UnityEngine;

namespace Game.Room.Enemy
{
    /// <summary>
    /// 能量提升道具
    /// </summary>
    public class PowerUp : SnapTransform
    {
        [Tooltip("线速度")] public Vector3 linearVelocity = Vector3.left * 2;

        private void Update()
        {
            transform.Translate(linearVelocity* Time.deltaTime);
        }
    }
}