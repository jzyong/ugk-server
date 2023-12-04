// Snapshot Interpolation: https://gafferongames.com/post/snapshot_interpolation/

using System;
using System.Collections.Generic;
using System.Runtime.CompilerServices;
using Common.Network.Serialize;
using Common.Tools;
using Common.Tools.SnapshotInterpolation;
using Google.Protobuf;
using Network.Sync;
using UnityEngine;

namespace Common.Network.Sync
{
    /// <summary>
    /// 快照同步
    /// 客户端只同步玩家自己操作的  TODO 待测试
    /// </summary>
    public class SnapTransform : NetworkTransform
    {
        //缓存的快照
        public readonly SortedList<double, TransformSnapshot> snapshots = new SortedList<double, TransformSnapshot>();

        [Header("Sync Only If Changed")]
        [Tooltip("When true, changes are not sent unless greater than sensitivity values below.")]
        public bool onlySyncOnChange = true;

        // interpolation is on by default, but can be disabled to jump to
        // the destination immediately. some projects need this.
        [Header("Interpolation")] [Tooltip("Set to false to have a snap-like effect on position movement.")]
        public bool interpolatePosition = true;

        [Tooltip("Set to false to have a snap-like effect on rotations.")]
        public bool interpolateRotation = true;

        [Tooltip(
            "Set to false to remove scale smoothing. Example use-case: Instant flipping of sprites that use -X and +X for direction.")]
        public bool interpolateScale = true;

        [Header("Rotation")] [Tooltip("Sensitivity of changes needed before an updated state is sent over the network")]
        public float rotationSensitivity = 0.01f;

        // Used to store last sent snapshots
        protected TransformSnapshot last;

        /// <summary>
        /// 同步数据
        /// </summary>
        public ByteString SyncData { get; set; }

        protected void Awake()
        {
        }


        /// <summary>
        /// 构建当前的快照
        /// </summary>
        /// <returns></returns>
        protected TransformSnapshot Construct()
        {
            // NetworkTime.localTime for double precision until Unity has it too
            return new TransformSnapshot(
                // our local time is what the other end uses as remote time
                Time.unscaledTime, // Unity 2019 doesn't have timeAsDouble yet
                0, // the other end fills out local time itself
                target.position,
                target.rotation,
                target.lossyScale
            );
        }

        /// <summary>
        /// 
        /// </summary>
        /// <param name="timeStamp"></param>
        /// <param name="position"></param>
        /// <param name="rotation"></param>
        /// <param name="scale"></param>
        protected void AddSnapshot(double timeStamp, Vector3? position,
            Quaternion? rotation, Vector3? scale)
        {
            // position, rotation, scale can have no value if same as last time.
            // saves bandwidth.
            // but we still need to feed it to snapshot interpolation. we can't
            // just have gaps in there if nothing has changed. for example, if
            //   client sends snapshot at t=0
            //   client sends nothing for 10s because not moved
            //   client sends snapshot at t=10
            // then the server would assume that it's one super slow move and
            // replay it for 10 seconds.
            if (!position.HasValue)
                position = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].position : target.position;
            if (!rotation.HasValue)
                rotation = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].rotation : target.rotation;
            if (!scale.HasValue)
                scale = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].scale : target.lossyScale;

            //Log.Info($"收到Transform position={position.Value} rotation={rotation.Value} scale={scale.Value} remoteTime={timeStamp} localTime={Time.unscaledTime}");

            // insert transform snapshot
            SnapshotInterpolation.InsertIfNotExists(snapshots, new TransformSnapshot(
                timeStamp, // arrival remote timestamp. NOT remote time.
                Time.unscaledTime, // Unity 2019 doesn't have timeAsDouble yet
                position.Value,
                rotation.Value,
                scale.Value
            ));
        }


        protected  void Apply(TransformSnapshot from, TransformSnapshot to, float t)
        {
            if (syncPosition)
            {
                target.position = interpolatePosition
                    ? Vector3.LerpUnclamped(from.position, to.position, (float)t)
                    : to.position;
            }

            if (syncRotation)
            {
                target.rotation = interpolateRotation
                    ? Quaternion.SlerpUnclamped(from.rotation, to.rotation, (float)t)
                    : to.rotation;
            }


            // Unity doesn't support setting world scale.
            // OnValidate disables syncScale in world mode.
            // else
            // target.lossyScale = scale; //
            if (syncScale)
            {
                target.localScale = interpolateScale ? Vector3.LerpUnclamped(from.scale, to.scale, t) : to.scale;
            }
        }

        /// <summary>
        /// 
        /// </summary>
        public  void Reset()
        {
            // disabled objects aren't updated anymore.
            // so let's clear the buffers.
            snapshots.Clear();
            // reset 'last' for delta too
            last = new TransformSnapshot(0, 0, Vector3.zero, Quaternion.identity, Vector3.zero);
        }

        protected  void OnEnable()
        {
            Reset();
        }

        protected  void OnDisable()
        {
            Reset();
        }


        void Update()
        {
            UpdateClient();
        }

        void LateUpdate()
        {
            if (Time.unscaledTime > nextSendTime && (!onlySyncOnChange || Changed(Construct())))
            {
                OnSerialize(false);
                nextSendTime += sendInterval;
            }
        }

        protected  void UpdateClient()
        {
            // only while we have snapshots
            if (snapshots.Count > 0)
            {
                // step the interpolation without touching time.
                // NetworkClient is responsible for time globally.
                SnapshotInterpolation.StepInterpolation(
                    snapshots,
                    Time.unscaledTime, // == NetworkClient.localTimeline from snapshot interpolation
                    out TransformSnapshot from,
                    out TransformSnapshot to,
                    out double t);

                //Log.Info($"快照数量:{snapshots.Count} from:{from.position} to:{to.position} t={t}");
                //插值并引用
                Apply(from, to, (float)t);
            }
        }

        /// <summary>
        /// 发送同步数据,序列化
        /// </summary>
        public void OnSerialize(bool initialState)
        {
            using (NetworkWriterPooled writer = NetworkWriterPool.Get())
            {
                // get current snapshot for broadcasting.
                TransformSnapshot snapshot = Construct();

                // initial
                if (initialState)
                {
                    if (last.remoteTime > 0) snapshot = last;
                    if (syncPosition) writer.WriteVector3(snapshot.position);
                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            writer.WriteUInt(Compression.CompressQuaternion(snapshot.rotation));
                        else
                            writer.WriteQuaternion(snapshot.rotation);
                    }

                    if (syncScale) writer.WriteVector3(snapshot.scale);
                }
                // delta
                else
                {
                    if (syncPosition)
                    {
                        // quantize -> delta -> varint
                        Compression.ScaleToLong(snapshot.position, positionPrecision, out Vector3Long quantized);
                        DeltaCompression.Compress(writer, lastSerializedPosition, quantized);
                    }

                    if (syncRotation)
                    {
                        // (optional) smallest three compression for now. no delta.
                        if (compressRotation)
                            writer.WriteUInt(Compression.CompressQuaternion(snapshot.rotation));
                        else
                            writer.WriteQuaternion(snapshot.rotation);
                    }

                    if (syncScale)
                    {
                        // quantize -> delta -> varint
                        Compression.ScaleToLong(snapshot.scale, scalePrecision, out Vector3Long quantized);
                        DeltaCompression.Compress(writer, lastSerializedScale, quantized);
                    }
                }

                // save serialized as 'last' for next delta compression
                if (syncPosition)
                    Compression.ScaleToLong(snapshot.position, positionPrecision, out lastSerializedPosition);
                if (syncScale) Compression.ScaleToLong(snapshot.scale, scalePrecision, out lastSerializedScale);

                // set 'last'
                last = snapshot;

                //发送数据
                SyncData = ByteString.CopyFrom(writer.ToArray());
            }
        }

        /// <summary>
        /// 接收同步数据
        /// </summary>
        public void OnDeserialize(UgkMessage ugkMessage, ByteString data, bool initialState)
        {
            var segment = new ArraySegment<byte>(data.ToByteArray());
            using (NetworkReaderPooled reader = NetworkReaderPool.Get(segment))
            {
                Vector3? position = null;
                Quaternion? rotation = null;
                Vector3? scale = null;

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
                }

                //保存sendInterval时间段快照数据，让平滑移动 
                OnReceiveTransform(position, rotation, scale, ugkMessage.GetTime() + sendInterval);

                // save deserialized as 'last' for next delta compression
                if (syncPosition)
                    Compression.ScaleToLong(position.Value, positionPrecision, out lastDeserializedPosition);
                if (syncScale) Compression.ScaleToLong(scale.Value, scalePrecision, out lastDeserializedScale);
            }
        }


        // check if position / rotation / scale changed since last sync 
        protected bool Changed(TransformSnapshot current) =>
            // position is quantized and delta compressed.
            // only consider it changed if the quantized representation is changed.
            QuantizedChanged(last.position, current.position, positionPrecision) ||
            // rotation isn't quantized / delta compressed.
            // check with sensitivity.
            Quaternion.Angle(last.rotation, current.rotation) > rotationSensitivity ||
            // scale is quantized and delta compressed.
            // only consider it changed if the quantized representation is changed.
            // careful: don't use 'serialized / deserialized last'. as it depends on sync mode etc.
            QuantizedChanged(last.scale, current.scale, scalePrecision);

        // helper function to compare quantized representations of a Vector3 
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        protected bool QuantizedChanged(Vector3 u, Vector3 v, float precision)
        {
            Compression.ScaleToLong(u, precision, out Vector3Long uQuantized);
            Compression.ScaleToLong(v, precision, out Vector3Long vQuantized);
            return uQuantized != vQuantized;
        }


        // server broadcasts sync message to all clients
        protected void OnReceiveTransform(Vector3? position, Quaternion? rotation, Vector3? scale, double timeStamp)
        {
            // 'only sync on change' needs a correction on every new move sequence.
            if (onlySyncOnChange && NeedsCorrection(snapshots, timeStamp, sendInterval, 2))
            {
                RewriteHistory(
                    snapshots,
                    timeStamp, // arrival remote timestamp. NOT remote timeline.
                    Time.unscaledTime, // Unity 2019 doesn't have timeAsDouble yet
                    sendInterval,
                    target.position,
                    target.rotation,
                    target.lossyScale);
            }

            AddSnapshot(timeStamp, position, rotation, scale);
        }

        // only sync on change /////////////////////////////////////////////////
        // snap interp. needs a continous flow of packets.
        // 'only sync on change' interrupts it while not changed.
        // once it restarts, snap interp. will interp from the last old position.
        // this will cause very noticeable stutter for the first move each time.
        // the fix is quite simple.

        // 1. detect if the remaining snapshot is too old from a past move. 
        static bool NeedsCorrection(
            SortedList<double, TransformSnapshot> snapshots,
            double remoteTimestamp,
            double bufferTime,
            double toleranceMultiplier) =>
            snapshots.Count == 1 &&
            remoteTimestamp - snapshots.Keys[0] >= bufferTime * toleranceMultiplier;

        // 2. insert a fake snapshot at current position,
        //    exactly one 'sendInterval' behind the newly received one.
        static void RewriteHistory(
            SortedList<double, TransformSnapshot> snapshots,
            // timestamp of packet arrival, not interpolated remote time!
            double remoteTimeStamp,
            double localTime,
            double sendInterval,
            Vector3 position,
            Quaternion rotation,
            Vector3 scale)
        {
            // clear the previous snapshot
            snapshots.Clear();

            // insert a fake one at where we used to be,
            // 'sendInterval' behind the new one.
            SnapshotInterpolation.InsertIfNotExists(snapshots, new TransformSnapshot(
                remoteTimeStamp - sendInterval, // arrival remote timestamp. NOT remote time.
                localTime - sendInterval, // Unity 2019 doesn't have timeAsDouble yet
                position,
                rotation,
                scale
            ));
        }
    }
}