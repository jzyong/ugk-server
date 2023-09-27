using System;

namespace Common.Network
{
    
    /// <summary>
    /// 消息处理映射
    /// </summary>
    [AttributeUsage(AttributeTargets.Method,Inherited = false,AllowMultiple = false)]
    public sealed class MessageMapAttribute :Attribute
    {
        /// <summary>
        /// 消息ID
        /// </summary>
        public readonly int mid;

        public MessageMapAttribute(int mid)
        {
            this.mid = mid;
        }
    }
}