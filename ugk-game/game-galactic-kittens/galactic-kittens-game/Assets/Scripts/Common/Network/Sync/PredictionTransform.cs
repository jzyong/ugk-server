using System;
using Common.Network.Serialize;
using Common.Tools;
using Google.Protobuf;
using UnityEngine;

namespace Common.Network.Sync
{
    /// <summary>
    /// 预测同步   
    /// <para>本地控制对象直接应用，然后对结果进行矫正；</para>
    /// <para>其他玩家或服务器控制对象使用航位推测进行计算，然后对结果进行矫正</para>
    /// <remarks>每帧都改变方向的曲线运动效果不理想，抖动厉害</remarks>
    /// </summary>
    public class PredictionTransform : NetworkTransform
    {
        [SerializeField] [Tooltip("是否同步角速度")] private bool syncAngularVelocity;

        [SerializeField] [Range(0.00_01f, 1f)] [Tooltip("速度精度")]
        private float velocityPrecision = 0.01f;

        /// <summary>
        /// 最后一次消息同步序列号，本地控制对象立即进行插值应用位置动画等，然后发送给权威服务器，如果客户端连续发送多个位置信息，可能应用老的服务器位置，因此通过序列号判断
        /// https://www.gabrielgambetta.com/client-side-prediction-server-reconciliation.html
        /// </summary>
        public uint LastMessageSeq { get; set; }

        /// <summary>
        /// 线速度
        /// </summary>
        public Vector3 LinearVelocity { get; set; }

        /// <summary>
        /// 角速度
        /// </summary>
        public Vector3 AngularVelocity { get; set; }
        
        /// <summary>
        /// 同步数据
        /// </summary>
        public ByteString SyncData { get; set; }


        protected Vector3Long lastSerializedLinearVelocity = Vector3Long.zero;
        protected Vector3Long lastDeserializedLinearVelocity = Vector3Long.zero;

        protected Vector3Long lastSerializedAngularVelocity = Vector3Long.zero;
        protected Vector3Long lastDeserializedAngularVelocity = Vector3Long.zero;


        protected override void Awake()
        {
            if (sendInterval < 1)
            {
                sendInterval = 1;
            }

            nextSendTime = sendInterval+Time.unscaledTime;
        }

        public void Update()
        {
            // 计算位置和方向并应用
            CalculateTransform(Time.deltaTime);
        }

        public void LateUpdate()
        {
            //超时强制同步一下
            if (Onwer && Time.unscaledTime > nextSendTime)
            {
                nextSendTime += sendInterval;
                OnSerialize(false);
            }
        }

        /// <summary>
        /// 当速度发生改变
        /// </summary>
        public void OnVelocityChange()
        {
            if (Onwer)
            {
                OnSerialize(false);
            }
        }

        /// <summary>
        /// 发送同步数据,序列化
        /// [位置，旋转，缩放，线速度，角速度]
        /// </summary>
        protected void OnSerialize(bool initialState)
        {
            using (NetworkWriterPooled writer = NetworkWriterPool.Get())
            {
                // initial
                if (initialState)
                {
                    if (syncPosition) writer.WriteVector3(target.position);
                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            writer.WriteUInt(Compression.CompressQuaternion(target.rotation));
                        else
                            writer.WriteQuaternion(target.rotation);
                    }

                    if (syncScale) writer.WriteVector3(target.lossyScale);
                    writer.WriteVector3(LinearVelocity);
                    if (syncAngularVelocity) writer.WriteVector3(AngularVelocity);
                }
                // delta
                else
                {
                    if (syncPosition)
                    {
                        // quantize -> delta -> varint
                        Compression.ScaleToLong(target.position, positionPrecision, out Vector3Long quantized);
                        DeltaCompression.Compress(writer, lastSerializedPosition, quantized);
                    }

                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            writer.WriteUInt(Compression.CompressQuaternion(target.rotation));
                        else
                            writer.WriteQuaternion(target.rotation);
                    }

                    if (syncScale)
                    {
                        // quantize -> delta -> varint
                        Compression.ScaleToLong(target.lossyScale, scalePrecision, out Vector3Long quantized);
                        DeltaCompression.Compress(writer, lastSerializedScale, quantized);
                    }

                    //线速度
                    Compression.ScaleToLong(LinearVelocity, velocityPrecision, out Vector3Long quantizedLinearVelocity);
                    DeltaCompression.Compress(writer, lastSerializedLinearVelocity, quantizedLinearVelocity);

                    if (syncAngularVelocity)
                    {
                        Compression.ScaleToLong(AngularVelocity, velocityPrecision,
                            out Vector3Long quantizedAngularVelocity);
                        DeltaCompression.Compress(writer, lastSerializedAngularVelocity, quantizedAngularVelocity);
                    }
                }

                // save serialized as 'last' for next delta compression
                if (syncPosition)
                    Compression.ScaleToLong(target.position, positionPrecision, out lastSerializedPosition);
                if (syncScale) Compression.ScaleToLong(target.lossyScale, scalePrecision, out lastSerializedScale);
                Compression.ScaleToLong(LinearVelocity, velocityPrecision, out lastSerializedLinearVelocity);
                if (syncAngularVelocity)
                    Compression.ScaleToLong(AngularVelocity, velocityPrecision, out lastSerializedAngularVelocity);

                //发送数据
               SyncData = ByteString.CopyFrom(writer.ToArray());
                
            }
        }

        /// <summary>
        /// 接收同步数据
        /// </summary>
        /// <param name="ugkMessage"></param>
        /// <param name="data"></param>
        /// <param name="initialState"></param>
        public void OnDeserialize(UgkMessage ugkMessage, ByteString data, bool initialState)
        {

            var segment = new ArraySegment<byte>(data.ToByteArray());
            using (NetworkReaderPooled reader = NetworkReaderPool.Get(segment))
            {
                //客户端 判断自己控制对象的位置，需要使用服务器的权威位置，但是又不能被老数据覆盖，因此通过序列号判断数据是否为老数据
                if (Onwer && ugkMessage.Seq <= LastMessageSeq)
                {
                    return;
                }

                Vector3? position = null;
                Quaternion? rotation = null;
                Vector3? scale = null;
                Vector3? linearVelocity = null;
                Vector3? angularVelocity = null;


                // initial
                if (initialState)
                {
                    if (syncPosition) position = reader.ReadVector3();
                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            rotation = Compression.DecompressQuaternion(reader.ReadUInt());
                        else
                            rotation = reader.ReadQuaternion();
                    }

                    if (syncScale) scale = reader.ReadVector3();
                    linearVelocity = reader.ReadVector3();
                    if (syncAngularVelocity) reader.ReadVector3();
                }
                // delta
                else
                {
                    // varint -> delta -> quantize
                    if (syncPosition)
                    {
                        Vector3Long quantized = DeltaCompression.Decompress(reader, lastDeserializedPosition);
                        position = Compression.ScaleToFloat(quantized, positionPrecision);
                    }

                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            rotation = Compression.DecompressQuaternion(reader.ReadUInt());
                        else
                            rotation = reader.ReadQuaternion();
                    }

                    if (syncScale)
                    {
                        Vector3Long quantized = DeltaCompression.Decompress(reader, lastDeserializedScale);
                        scale = Compression.ScaleToFloat(quantized, scalePrecision);
                    }

                    Vector3Long quantizedLineVelocity =
                        DeltaCompression.Decompress(reader, lastDeserializedLinearVelocity);
                    linearVelocity = Compression.ScaleToFloat(quantizedLineVelocity, velocityPrecision);

                    if (syncAngularVelocity)
                    {
                        Vector3Long quantizedAngularVelocity =
                            DeltaCompression.Decompress(reader, lastDeserializedAngularVelocity);
                        angularVelocity = Compression.ScaleToFloat(quantizedAngularVelocity, velocityPrecision);
                    }
                }

                ReconciliationTransform(position, rotation, scale, linearVelocity, angularVelocity, ugkMessage);
                // save deserialized as 'last' for next delta compression
                if (syncPosition)
                    Compression.ScaleToLong(position.Value, positionPrecision, out lastDeserializedPosition);
                if (syncScale) Compression.ScaleToLong(scale.Value, scalePrecision, out lastDeserializedScale);
                Compression.ScaleToLong(linearVelocity.Value, velocityPrecision, out lastDeserializedLinearVelocity);
                if (syncAngularVelocity)
                {
                    Compression.ScaleToLong(angularVelocity.Value, velocityPrecision,
                        out lastDeserializedAngularVelocity);
                }
            }
        }

        /// <summary>
        /// 校对位置
        /// </summary>
        /// <param name="position"></param>
        /// <param name="rotation"></param>
        /// <param name="scale"></param>
        /// <param name="linearVelocity"></param>
        /// <param name="angularVelocity"></param>
        private void ReconciliationTransform(Vector3? position, Quaternion? rotation, Vector3? scale,
            Vector3? linearVelocity, Vector3? angularVelocity, UgkMessage ugkMessage)
        {
            if (linearVelocity.HasValue) LinearVelocity = linearVelocity.Value;
            if (angularVelocity.HasValue) AngularVelocity = angularVelocity.Value;
            if (scale.HasValue) target.localScale = scale.Value;

            //直接赋值，是否会抖动？
            //进行时间补偿，使与服务器一致
            double delayTime = Time.unscaledTime - ugkMessage.GetTime();
            if (delayTime > 0)
            {
                // 进行位置计算
                CalculateTransform((float)delayTime);
            }
            else
            {
                if (position.HasValue) target.position = position.Value;
                if (rotation.HasValue) target.rotation = rotation.Value;
            }
        }

        /// <summary>
        /// 计算位置方向
        /// </summary>
        /// <param name="deltaTime"></param>
        private void CalculateTransform(float deltaTime)
        {
            target.position += LinearVelocity * deltaTime;
            if (syncAngularVelocity)
            {
                // var rotation = Quaternion.Euler(AngularVelocity * deltaTime);
                // target.transform.rotation = Quaternion.LerpUnclamped(target.transform.rotation, rotation, 1);
                target.transform.Rotate(deltaTime*AngularVelocity);
            }
        }
        
        /// <summary>
        /// 设置最后一次反序列化缓存的线速度，增量压缩还原需要
        /// </summary>
        /// <param name="position"></param>
        public void SetLastDeserializedLinearVelocity(Vector3 position)
        {
            Compression.ScaleToLong(position, velocityPrecision, out lastDeserializedLinearVelocity);
        }
        
        /// <summary>
        /// 设置最后一次反序列化缓存的线速度，增量压缩还原需要
        /// </summary>
        /// <param name="position"></param>
        public void SetLastDeserializedAngularVelocity(Vector3 position)
        {
            Compression.ScaleToLong(position, velocityPrecision, out lastDeserializedAngularVelocity);
        }
        
    }
}