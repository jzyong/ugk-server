using UnityEngine;

namespace Game.Room.Player
{
    /// <summary>
    /// 飞船子弹
    /// </summary>
    public class SpaceshipBullet : MonoBehaviour
    {
         [Tooltip("线速度")] public Vector3 linearVelocity = Vector3.right * 10;
    }
}