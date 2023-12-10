using System;
using System.Collections.Generic;
using kcp2k;

namespace Common.Network.Serialize
{
    /// <summary>
    /// 小的消息进行合并，在每帧结束后批量发送
    /// <para>Mirror 有对应的Unbatcher,先取出消息放入队列，再从队列中依次去拿取，我们直接解析执行</para>
    /// <para>服务器进行批量是否多余？因为这里合并了，网关又拆开了，再合并发送给客户端？</para>
    /// </summary>
    public class Batcher
    {
        //Mirror 阈值使用mtu或者接收端可接受的最大容器，因为批量消息共用一个时间戳。但是我们是每个消息共用一个时间戳，因此最大消息不超过mtu
        //大消息包分片再合并可能会额外增加延迟 ， 内部阈值可以直接使用接收端的容量大小？
        readonly int threshold = Kcp.MTU_DEF - 50;



        // full batches ready to be sent.
        readonly Queue<NetworkWriterPooled> batches = new Queue<NetworkWriterPooled>();

        // current batch in progress
        NetworkWriterPooled batch;


        // add a message for batching
        public void AddMessage(ArraySegment<byte> message)
        {
            // when appending to a batch in progress, check final size.
            // if it expands beyond threshold, then we should finalize it first.
            // => less than or exactly threshold is fine.
            //    GetBatch() will finalize it.
            // => see unit tests.
            if (batch != null &&
                batch.Position + message.Count > threshold)
            {
                batches.Enqueue(batch);
                batch =NetworkWriterPool.Get();
            }

            // initialize a new batch if necessary
            if (batch == null)
            {
                // borrow from pool. we return it in GetBatch.
                batch = NetworkWriterPool.Get();
            }

            batch.WriteBytes(message.Array, message.Offset, message.Count);
        }

        // helper function to copy a batch to writer and return it to pool
        static void CopyAndReturn(NetworkWriterPooled batch, NetworkWriter writer)
        {
            // make sure the writer is fresh to avoid uncertain situations
            if (writer.Position != 0)
                throw new ArgumentException($"GetBatch needs a fresh writer!");

            // copy to the target writer
            ArraySegment<byte> segment = batch.ToArraySegment();
            writer.WriteBytes(segment.Array, segment.Offset, segment.Count);

            // return batch to pool for reuse
            NetworkWriterPool.Return(batch);
        }

        /// <summary>
        /// 是否有待发送到消息
        /// </summary>
        /// <returns></returns>
        public bool HasMessage()
        {
            if (batches.Count > 0 || batch != null)
            {
                return true;
            }

            return false;
        }

        // get the next batch which is available for sending (if any).
        public bool GetBatch(NetworkWriter writer)
        {
            // get first batch from queue (if any)
            if (batches.TryDequeue(out NetworkWriterPooled first))
            {
                CopyAndReturn(first, writer);
                return true;
            }

            // if queue was empty, we can send the batch in progress.
            if (batch != null)
            {
                CopyAndReturn(batch, writer);
                batch = null;
                return true;
            }

            // nothing was written
            return false;
        }
    }
}