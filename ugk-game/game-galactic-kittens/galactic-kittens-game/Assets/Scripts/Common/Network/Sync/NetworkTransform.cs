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
        [Header("Target")]
        [Tooltip("The Transform component to sync. May be on on this GameObject, or on a child.")]
        public Transform target;
        

        [Header("Selective Sync\nDon't change these at Runtime")]
        public bool syncPosition = true;  // do not change at runtime!
        public bool syncRotation = true;  // do not change at runtime!
        public bool syncScale = false; // do not change at runtime! rare. off by default.

  

   

       
        protected virtual void Awake() { }

        protected  void OnValidate()
        {

            // set target to self if none yet
            if (target == null) target = transform;

            // time snapshot interpolation happens globally.
            // value (transform) happens in here.
            // both always need to be on the same send interval.
            // force the setting to '0' in OnValidate to make it obvious that we
            // actually use NetworkServer.sendInterval.
            syncInterval = 0;

            // obsolete clientAuthority compatibility:
            // if it was used, then set the new SyncDirection automatically.
            // if it wasn't used, then don't touch syncDirection.
#pragma warning disable CS0618
            if (clientAuthority)
            {
                syncDirection = SyncDirection.ClientToServer;
                Debug.LogWarning($"{name}'s NetworkTransform component has obsolete .clientAuthority enabled. Please disable it and set SyncDirection to ClientToServer instead.");
            }
#pragma warning restore CS0618
        }

        // snapshot functions //////////////////////////////////////////////////
        // construct a snapshot of the current state
        // => internal for testing
        protected virtual TransformSnapshot Construct()
        {
            // NetworkTime.localTime for double precision until Unity has it too
            return new TransformSnapshot(
                // our local time is what the other end uses as remote time
                NetworkTime.localTime, // Unity 2019 doesn't have timeAsDouble yet
                0,                     // the other end fills out local time itself
                target.localPosition,
                target.localRotation,
                target.localScale
            );
        }

        /// <summary>
        /// @
        /// </summary>
        /// <param name="snapshots"></param>
        /// <param name="timeStamp"></param>
        /// <param name="position"></param>
        /// <param name="rotation"></param>
        /// <param name="scale"></param>
        protected void AddSnapshot(SortedList<double, TransformSnapshot> snapshots, double timeStamp, Vector3? position, Quaternion? rotation, Vector3? scale)
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
            if (!position.HasValue) position = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].position : target.localPosition;
            if (!rotation.HasValue) rotation = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].rotation : target.localRotation;
            if (!scale.HasValue) scale = snapshots.Count > 0 ? snapshots.Values[snapshots.Count - 1].scale : target.localScale;

            // insert transform snapshot
            SnapshotInterpolation.InsertIfNotExists(snapshots, new TransformSnapshot(
                timeStamp, // arrival remote timestamp. NOT remote time.
                NetworkTime.localTime, // Unity 2019 doesn't have timeAsDouble yet
                position.Value,
                rotation.Value,
                scale.Value
            ));
        }

        // apply a snapshot to the Transform.
        // -> start, end, interpolated are all passed in caes they are needed
        // -> a regular game would apply the 'interpolated' snapshot
        // -> a board game might want to jump to 'goal' directly
        // (it's easier to always interpolate and then apply selectively,
        //  instead of manually interpolating x, y, z, ... depending on flags)
        // => internal for testing
        //
        // NOTE: stuck detection is unnecessary here.
        //       we always set transform.position anyway, we can't get stuck. @
        protected virtual void Apply(TransformSnapshot interpolated, TransformSnapshot endGoal)
        {
            // local position/rotation for VR support
            //
            // if syncPosition/Rotation/Scale is disabled then we received nulls
            // -> current position/rotation/scale would've been added as snapshot
            // -> we still interpolated
            // -> but simply don't apply it. if the user doesn't want to sync
            //    scale, then we should not touch scale etc.

            if (syncPosition)
                target.localPosition = interpolatePosition ? interpolated.position : endGoal.position;

            if (syncRotation)
                target.localRotation = interpolateRotation ? interpolated.rotation : endGoal.rotation;

            if (syncScale)
                target.localScale = interpolateScale ? interpolated.scale : endGoal.scale;
        }

        // client->server teleport to force position without interpolation.
        // otherwise it would interpolate to a (far away) new position.
        // => manually calling Teleport is the only 100% reliable solution. @
        [Command]
        public void CmdTeleport(Vector3 destination)
        {
            // client can only teleport objects that it has authority over.
            if (syncDirection != SyncDirection.ClientToServer) return;

            // TODO what about host mode?
            OnTeleport(destination);

            // if a client teleports, we need to broadcast to everyone else too
            // TODO the teleported client should ignore the rpc though.
            //      otherwise if it already moved again after teleporting,
            //      the rpc would come a little bit later and reset it once.
            // TODO or not? if client ONLY calls Teleport(pos), the position
            //      would only be set after the rpc. unless the client calls
            //      BOTH Teleport(pos) and target.position=pos
            RpcTeleport(destination);
        }

        // client->server teleport to force position and rotation without interpolation.
        // otherwise it would interpolate to a (far away) new position.
        // => manually calling Teleport is the only 100% reliable solution. @
        [Command]
        public void CmdTeleport(Vector3 destination, Quaternion rotation)
        {
            // client can only teleport objects that it has authority over.
            if (syncDirection != SyncDirection.ClientToServer) return;

            // TODO what about host mode?
            OnTeleport(destination, rotation);

            // if a client teleports, we need to broadcast to everyone else too
            // TODO the teleported client should ignore the rpc though.
            //      otherwise if it already moved again after teleporting,
            //      the rpc would come a little bit later and reset it once.
            // TODO or not? if client ONLY calls Teleport(pos), the position
            //      would only be set after the rpc. unless the client calls
            //      BOTH Teleport(pos) and target.position=pos
            RpcTeleport(destination, rotation);
        }

        // server->client teleport to force position without interpolation.
        // otherwise it would interpolate to a (far away) new position.
        // => manually calling Teleport is the only 100% reliable solution. @
        [ClientRpc]
        public void RpcTeleport(Vector3 destination)
        {
            // NOTE: even in client authority mode, the server is always allowed
            //       to teleport the player. for example:
            //       * CmdEnterPortal() might teleport the player
            //       * Some people use client authority with server sided checks
            //         so the server should be able to reset position if needed.

            // TODO what about host mode?
            OnTeleport(destination);
        }

        // server->client teleport to force position and rotation without interpolation.
        // otherwise it would interpolate to a (far away) new position.
        // => manually calling Teleport is the only 100% reliable solution. @
        [ClientRpc]
        public void RpcTeleport(Vector3 destination, Quaternion rotation)
        {
            // NOTE: even in client authority mode, the server is always allowed
            //       to teleport the player. for example:
            //       * CmdEnterPortal() might teleport the player
            //       * Some people use client authority with server sided checks
            //         so the server should be able to reset position if needed.

            // TODO what about host mode?
            OnTeleport(destination, rotation);
        }

        [ClientRpc]
        void RpcReset()
        {
            Reset();
        }

        // common Teleport code for client->server and server->client @
        protected virtual void OnTeleport(Vector3 destination)
        {
            // reset any in-progress interpolation & buffers
            Reset();

            // set the new position.
            // interpolation will automatically continue.
            target.position = destination;

            // TODO
            // what if we still receive a snapshot from before the interpolation?
            // it could easily happen over unreliable.
            // -> maybe add destination as first entry?
        }

        // common Teleport code for client->server and server->client @
        protected virtual void OnTeleport(Vector3 destination, Quaternion rotation)
        {
            // reset any in-progress interpolation & buffers
            Reset();

            // set the new position.
            // interpolation will automatically continue.
            target.position = destination;
            target.rotation = rotation;

            // TODO
            // what if we still receive a snapshot from before the interpolation?
            // it could easily happen over unreliable.
            // -> maybe add destination as first entry?
        }

        /// <summary>
        /// @
        /// </summary>
        public virtual void Reset()
        {
            // disabled objects aren't updated anymore.
            // so let's clear the buffers.
            snapshots.Clear();
            clientSnapshots.Clear();
        }

        protected virtual void OnEnable()
        {
            Reset();

            if (NetworkServer.active)
                NetworkIdentity.clientAuthorityCallback += OnClientAuthorityChanged;
        }

        protected virtual void OnDisable()
        {
            Reset();

            if (NetworkServer.active)
                NetworkIdentity.clientAuthorityCallback -= OnClientAuthorityChanged;
        }

        /// <summary>
        /// @
        /// </summary>
        /// <param name="conn"></param>
        /// <param name="identity"></param>
        /// <param name="authorityState"></param>
        [ServerCallback]
        void OnClientAuthorityChanged(NetworkConnectionToClient conn, NetworkIdentity identity, bool authorityState)
        {
            if (identity != netIdentity) return;

            // If server gets authority or syncdirection is server to client,
            // we don't reset buffers.
            // This is because if syncdirection is S to C, we will never have
            // snapshot issues since there is only ever 1 source.

            if (syncDirection == SyncDirection.ClientToServer)
            {
                Reset();
                RpcReset();
            }
        }

        // OnGUI allocates even if it does nothing. avoid in release.
#if UNITY_EDITOR || DEVELOPMENT_BUILD
        // debug ///////////////////////////////////////////////////////////////
        protected virtual void OnGUI()
        {
            if (!showOverlay) return;
            if (!Camera.main) return;

            // show data next to player for easier debugging. this is very useful!
            // IMPORTANT: this is basically an ESP hack for shooter games.
            //            DO NOT make this available with a hotkey in release builds
            if (!Debug.isDebugBuild) return;

            // project position to screen
            Vector3 point = Camera.main.WorldToScreenPoint(target.position);

            // enough alpha, in front of camera and in screen?
            if (point.z >= 0 && Utils.IsPointInScreen(point))
            {
                GUI.color = overlayColor;
                GUILayout.BeginArea(new Rect(point.x, Screen.height - point.y, 200, 100));

                // always show both client & server buffers so it's super
                // obvious if we accidentally populate both.
                GUILayout.Label($"Server Buffer:{snapshots.Count}");
                GUILayout.Label($"Client Buffer:{clientSnapshots.Count}");

                GUILayout.EndArea();
                GUI.color = Color.white;
            }
        }

        protected virtual void DrawGizmos(SortedList<double, TransformSnapshot> buffer)
        {
            // only draw if we have at least two entries
            if (buffer.Count < 2) return;

            // calculate threshold for 'old enough' snapshots
            double threshold = NetworkTime.localTime - NetworkClient.bufferTime;
            Color oldEnoughColor = new Color(0, 1, 0, 0.5f);
            Color notOldEnoughColor = new Color(0.5f, 0.5f, 0.5f, 0.3f);

            // draw the whole buffer for easier debugging.
            // it's worth seeing how much we have buffered ahead already
            for (int i = 0; i < buffer.Count; ++i)
            {
                // color depends on if old enough or not
                TransformSnapshot entry = buffer.Values[i];
                bool oldEnough = entry.localTime <= threshold;
                Gizmos.color = oldEnough ? oldEnoughColor : notOldEnoughColor;
                Gizmos.DrawWireCube(entry.position, Vector3.one);
            }

            // extra: lines between start<->position<->goal
            Gizmos.color = Color.green;
            Gizmos.DrawLine(buffer.Values[0].position, target.position);
            Gizmos.color = Color.white;
            Gizmos.DrawLine(target.position, buffer.Values[1].position);
        }

        protected virtual void OnDrawGizmos()
        {
            // This fires in edit mode but that spams NRE's so check isPlaying
            if (!Application.isPlaying) return;
            if (!showGizmos) return;

            if (isServer) DrawGizmos(snapshots);
            if (isClient) DrawGizmos(clientSnapshots);
        }
#endif
        
        
        
        
        
        
        [Header("Sync Only If Changed")]
        [Tooltip("When true, changes are not sent unless greater than sensitivity values below.")]
        public bool onlySyncOnChange = true;

        uint sendIntervalCounter = 0;
        double lastSendIntervalTime = double.MinValue;

        [Tooltip("If we only sync on change, then we need to correct old snapshots if more time than sendInterval * multiplier has elapsed.\n\nOtherwise the first move will always start interpolating from the last move sequence's time, which will make it stutter when starting every time.")]
        public float onlySyncOnChangeCorrectionMultiplier = 2;

        [Header("Rotation")]
        [Tooltip("Sensitivity of changes needed before an updated state is sent over the network")]
        public float rotationSensitivity = 0.01f;
        [Tooltip("Apply smallest-three quaternion compression. This is lossy, you can disable it if the small rotation inaccuracies are noticeable in your project.")]
        public bool compressRotation = false;

        // delta compression is capable of detecting byte-level changes.
        // if we scale float position to bytes,
        // then small movements will only change one byte.
        // this gives optimal bandwidth.
        //   benchmark with 0.01 precision: 130 KB/s => 60 KB/s
        //   benchmark with 0.1  precision: 130 KB/s => 30 KB/s
        [Header("Precision")]
        [Tooltip("Position is rounded in order to drastically minimize bandwidth.\n\nFor example, a precision of 0.01 rounds to a centimeter. In other words, sub-centimeter movements aren't synced until they eventually exceeded an actual centimeter.\n\nDepending on how important the object is, a precision of 0.01-0.10 (1-10 cm) is recommended.\n\nFor example, even a 1cm precision combined with delta compression cuts the Benchmark demo's bandwidth in half, compared to sending every tiny change.")]
        [Range(0.00_01f, 1f)]                   // disallow 0 division. 1mm to 1m precision is enough range.
        public float positionPrecision = 0.01f; // 1 cm
        [Range(0.00_01f, 1f)]                   // disallow 0 division. 1mm to 1m precision is enough range.
        public float scalePrecision = 0.01f; // 1 cm

        // delta compression needs to remember 'last' to compress against
        protected Vector3Long lastSerializedPosition = Vector3Long.zero;
        protected Vector3Long lastDeserializedPosition = Vector3Long.zero;

        protected Vector3Long lastSerializedScale = Vector3Long.zero;
        protected Vector3Long lastDeserializedScale = Vector3Long.zero;

        // Used to store last sent snapshots
        protected TransformSnapshot last;

        protected int lastClientCount = 1;

        // update //////////////////////////////////////////////////////////////
        void Update()
        {
            // if server then always sync to others.
            if (isServer) UpdateServer();
            // 'else if' because host mode shouldn't send anything to server.
            // it is the server. don't overwrite anything there.
            else if (isClient) UpdateClient();
        }

        void LateUpdate()
        {
            // set dirty to trigger OnSerialize. either always, or only if changed.
            // It has to be checked in LateUpdate() for onlySyncOnChange to avoid
            // the possibility of Update() running first before the object's movement
            // script's Update(), which then causes NT to send every alternate frame
            // instead.
            if (isServer || (IsClientWithAuthority && NetworkClient.ready))
            {
                if (sendIntervalCounter == sendIntervalMultiplier && (!onlySyncOnChange || Changed(Construct())))
                    SetDirty();

                CheckLastSendTime();
            }
        }

        protected virtual void UpdateServer()
        {
            // apply buffered snapshots IF client authority
            // -> in server authority, server moves the object
            //    so no need to apply any snapshots there.
            // -> don't apply for host mode player objects either, even if in
            //    client authority mode. if it doesn't go over the network,
            //    then we don't need to do anything.
            // -> connectionToClient is briefly null after scene changes:
            //    https://github.com/MirrorNetworking/Mirror/issues/3329
            if (syncDirection == SyncDirection.ClientToServer &&
                connectionToClient != null &&
                !isOwned)
            {
                if (snapshots.Count > 0)
                {
                    // step the transform interpolation without touching time.
                    // NetworkClient is responsible for time globally.
                    SnapshotInterpolation.StepInterpolation(
                        snapshots,
                        connectionToClient.remoteTimeline,
                        out TransformSnapshot from,
                        out TransformSnapshot to,
                        out double t);

                    // interpolate & apply
                    TransformSnapshot computed = TransformSnapshot.Interpolate(from, to, t);
                    Apply(computed, to);
                }
            }
        }

        protected virtual void UpdateClient()
        {
            // client authority, and local player (= allowed to move myself)?
            if (!IsClientWithAuthority)
            {
                // only while we have snapshots
                if (clientSnapshots.Count > 0)
                {
                    // step the interpolation without touching time.
                    // NetworkClient is responsible for time globally.
                    SnapshotInterpolation.StepInterpolation(
                        clientSnapshots,
                        NetworkTime.time, // == NetworkClient.localTimeline from snapshot interpolation
                        out TransformSnapshot from,
                        out TransformSnapshot to,
                        out double t);

                    // interpolate & apply
                    TransformSnapshot computed = TransformSnapshot.Interpolate(from, to, t);
                    Apply(computed, to);
                }

                lastClientCount = clientSnapshots.Count;
            }
        }

        /// <summary>
        /// @
        /// </summary>
        protected virtual void CheckLastSendTime()
        {
            // timeAsDouble not available in older Unity versions.
            if (AccurateInterval.Elapsed(NetworkTime.localTime, NetworkServer.sendInterval, ref lastSendIntervalTime))
            {
                if (sendIntervalCounter == sendIntervalMultiplier)
                    sendIntervalCounter = 0;
                sendIntervalCounter++;
            }
        }

        // check if position / rotation / scale changed since last sync @
        protected virtual bool Changed(TransformSnapshot current) =>
            // position is quantized and delta compressed.
            // only consider it changed if the quantized representation is changed.
            // careful: don't use 'serialized / deserialized last'. as it depends on sync mode etc.
            QuantizedChanged(last.position, current.position, positionPrecision) ||
            // rotation isn't quantized / delta compressed.
            // check with sensitivity.
            Quaternion.Angle(last.rotation, current.rotation) > rotationSensitivity ||
            // scale is quantized and delta compressed.
            // only consider it changed if the quantized representation is changed.
            // careful: don't use 'serialized / deserialized last'. as it depends on sync mode etc.
            QuantizedChanged(last.scale, current.scale, scalePrecision);

        // helper function to compare quantized representations of a Vector3 @
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        protected bool QuantizedChanged(Vector3 u, Vector3 v, float precision)
        {
            Compression.ScaleToLong(u, precision, out Vector3Long uQuantized);
            Compression.ScaleToLong(v, precision, out Vector3Long vQuantized);
            return uQuantized != vQuantized;
        }

        // NT may be used on client/server/host to Owner/Observers with
        // ServerToClient or ClientToServer.
        // however, OnSerialize should always delta against last.
        public override void OnSerialize(NetworkWriter writer, bool initialState)
        {
            // get current snapshot for broadcasting.
            TransformSnapshot snapshot = Construct();

            // ClientToServer optimization:
            // for interpolated client owned identities,
            // always broadcast the latest known snapshot so other clients can
            // interpolate immediately instead of catching up too

            // TODO dirty mask? [compression is very good w/o it already]
            // each vector's component is delta compressed.
            // an unchanged component would still require 1 byte.
            // let's use a dirty bit mask to filter those out as well.

            // initial
            if (initialState)
            {
                // If there is a last serialized snapshot, we use it.
                // This prevents the new client getting a snapshot that is different
                // from what the older clients last got. If this happens, and on the next
                // regular serialisation the delta compression will get wrong values.
                // Notes:
                // 1. Interestingly only the older clients have it wrong, because at the end
                //    of this function, last = snapshot which is the initial state's snapshot
                // 2. Regular NTR gets by this bug because it sends every frame anyway so initialstate
                //    snapshot constructed would have been the same as the last anyway.
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
                // int before = writer.Position;

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
            if (syncPosition) Compression.ScaleToLong(snapshot.position, positionPrecision, out lastSerializedPosition);
            if (syncScale) Compression.ScaleToLong(snapshot.scale, scalePrecision, out lastSerializedScale);

            // set 'last'
            last = snapshot;
        }

        public override void OnDeserialize(NetworkReader reader, bool initialState)
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

            // handle depending on server / client / host.
            // server has priority for host mode.
            if (isServer) OnClientToServerSync(position, rotation, scale);
            else if (isClient) OnServerToClientSync(position, rotation, scale);

            // save deserialized as 'last' for next delta compression
            if (syncPosition) Compression.ScaleToLong(position.Value, positionPrecision, out lastDeserializedPosition);
            if (syncScale) Compression.ScaleToLong(scale.Value, scalePrecision, out lastDeserializedScale);
        }

        // sync ////////////////////////////////////////////////////////////////

        // local authority client sends sync message to server for broadcasting
        protected virtual void OnClientToServerSync(Vector3? position, Quaternion? rotation, Vector3? scale)
        {
            // only apply if in client authority mode
            if (syncDirection != SyncDirection.ClientToServer) return;

            // protect against ever growing buffer size attacks
            if (snapshots.Count >= connectionToClient.snapshotBufferSizeLimit) return;

            // 'only sync on change' needs a correction on every new move sequence.
            if (onlySyncOnChange &&
                NeedsCorrection(snapshots, connectionToClient.remoteTimeStamp, NetworkServer.sendInterval * sendIntervalMultiplier, onlySyncOnChangeCorrectionMultiplier))
            {
                RewriteHistory(
                    snapshots,
                    connectionToClient.remoteTimeStamp,
                    NetworkTime.localTime,                                  // arrival remote timestamp. NOT remote timeline.
                    NetworkServer.sendInterval * sendIntervalMultiplier,    // Unity 2019 doesn't have timeAsDouble yet
                    target.localPosition,
                    target.localRotation,
                    target.localScale);
            }

            // add a small timeline offset to account for decoupled arrival of
            // NetworkTime and NetworkTransform snapshots.
            // needs to be sendInterval. half sendInterval doesn't solve it.
            // https://github.com/MirrorNetworking/Mirror/issues/3427
            // remove this after LocalWorldState.
            AddSnapshot(snapshots, connectionToClient.remoteTimeStamp + timeStampAdjustment + offset, position, rotation, scale);
        }

        // server broadcasts sync message to all clients
        protected virtual void OnServerToClientSync(Vector3? position, Quaternion? rotation, Vector3? scale)
        {
            // don't apply for local player with authority
            if (IsClientWithAuthority) return;

            // 'only sync on change' needs a correction on every new move sequence.
            if (onlySyncOnChange &&
                NeedsCorrection(clientSnapshots, NetworkClient.connection.remoteTimeStamp, NetworkClient.sendInterval * sendIntervalMultiplier, onlySyncOnChangeCorrectionMultiplier))
            {
                RewriteHistory(
                    clientSnapshots,
                    NetworkClient.connection.remoteTimeStamp,               // arrival remote timestamp. NOT remote timeline.
                    NetworkTime.localTime,                                  // Unity 2019 doesn't have timeAsDouble yet
                    NetworkClient.sendInterval * sendIntervalMultiplier,
                    target.localPosition,
                    target.localRotation,
                    target.localScale);
            }

            // add a small timeline offset to account for decoupled arrival of
            // NetworkTime and NetworkTransform snapshots.
            // needs to be sendInterval. half sendInterval doesn't solve it.
            // https://github.com/MirrorNetworking/Mirror/issues/3427
            // remove this after LocalWorldState.
            AddSnapshot(clientSnapshots, NetworkClient.connection.remoteTimeStamp + timeStampAdjustment + offset, position, rotation, scale);
        }

        // only sync on change /////////////////////////////////////////////////
        // snap interp. needs a continous flow of packets.
        // 'only sync on change' interrupts it while not changed.
        // once it restarts, snap interp. will interp from the last old position.
        // this will cause very noticeable stutter for the first move each time.
        // the fix is quite simple.

        // 1. detect if the remaining snapshot is too old from a past move. @
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
                localTime - sendInterval,       // Unity 2019 doesn't have timeAsDouble yet
                position,
                rotation,
                scale
            ));
        }

        public override void Reset()
        {
            base.Reset();

            // reset delta
            lastSerializedPosition = Vector3Long.zero;
            lastDeserializedPosition = Vector3Long.zero;

            lastSerializedScale = Vector3Long.zero;
            lastDeserializedScale = Vector3Long.zero;

            // reset 'last' for delta too
            last = new TransformSnapshot(0, 0, Vector3.zero, Quaternion.identity, Vector3.zero);
        } 
        
        
        
        
        
        
        
        
    }
}