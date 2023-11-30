using System;
using System.Runtime.CompilerServices;
using kcp2k;

namespace Common.Network
{
    /// <summary>
    /// 自定义消息
    /// </summary>
    public class UgkMessage : IDisposable
    {
        /// <summary>
        /// 消息id
        /// </summary>
        public UInt32 MessageId { get; set; }
        
        /// <summary>
        /// 玩家ID
        /// </summary>
        public Int64 PlayerId { get; set; }
        
        /// <summary>
        /// 序列号
        /// </summary>
        public UInt32 Seq { get; set; }
        /// <summary>
        /// 时间戳
        /// </summary>
        public Int64 TimeStamp { get; set; }
        /// <summary>
        /// proto内容
        /// </summary>
        public byte[] Bytes { get; set; }
        
        public void Dispose() => UgkMessagePool.Return(this);

        public void Reset()
        {
            MessageId = 0;
            Seq = 0;
            TimeStamp = 0;
            PlayerId = 0;
            Bytes = null;
        }
        /// <summary>
        /// long时间转 double
        /// </summary>
        /// <returns></returns>
        public double GetTime()
        {
            return TimeStamp / 1000d;
        }
    }
    
   

    /// <summary>
    /// 消息缓冲池
    /// </summary>
    public static class UgkMessagePool
    {
        static readonly Pool<UgkMessage> Pool = new Pool<UgkMessage>(
            () => new UgkMessage(),
            (message)=>message.Reset(),
            // initial capacity to avoid allocations in the first few frames
            200
        );

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static UgkMessage Get()
        {
            UgkMessage message = Pool.Take();
            return message;
        }
        /// <summary>Returns a reader to the pool. </summary>
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static void Return(UgkMessage message)
        {
            Pool.Return(message);
        }
    }
}