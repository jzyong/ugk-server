using UnityEngine;

namespace Game.Messages
{
    /// <summary>
    /// 协议工具
    /// </summary>
    public static class ProtoUtil
    {
        /// <summary>
        /// proto坐标转换为Unity坐标
        /// </summary>
        /// <param name="vector3D"></param>
        /// <returns></returns>
        public static Vector3 BuildVector3(Vector3D vector3D)
        {
            return new Vector3(vector3D.X, vector3D.Y, vector3D.Z);
        }

        /// <summary>
        /// unity坐标转proto坐标
        /// </summary>
        /// <param name="vector3"></param>
        /// <returns></returns>
        public static Vector3D BuildVector3D(Vector3 vector3)
        {
            return new Vector3D()
            {
                X = vector3.x,
                Y = vector3.y,
                Z = vector3.z
            };
        }
    }
}