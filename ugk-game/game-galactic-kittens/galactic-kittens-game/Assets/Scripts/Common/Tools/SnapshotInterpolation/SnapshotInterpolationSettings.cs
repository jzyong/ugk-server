// snapshot interpolation settings struct.
// can easily be exposed in Unity inspectors.
using System;
using UnityEngine;

namespace Common.Tools.SnapshotInterpolation
{
    // class so we can define defaults easily 
    [Serializable]
    public class SnapshotInterpolationSettings
    {
        // decrease bufferTime at runtime to see the catchup effect.
        // increase to see slowdown.
        // 'double' so we can have very precise dynamic adjustment without rounding 
        //Mirror 是30帧同步，因此为2，但是我们是1帧，选择
        [Header("Buffering")]
        [Tooltip("Local simulation is behind by sendInterval * multiplier seconds.\n\nThis guarantees that we always have enough snapshots in the buffer to mitigate lags & jitter.\n\nIncrease this if the simulation isn't smooth. By default, it should be around 2.")]
        // public double bufferTimeMultiplier = 2;
        public double bufferTimeMultiplier = -0.033;

        // catchup /////////////////////////////////////////////////////////////
        // catchup thresholds in 'frames'.
        // half a frame might be too aggressive. 
        [Header("Catchup / Slowdown")]
        [Tooltip("Slowdown begins when the local timeline is moving too fast towards remote time. Threshold is in frames worth of snapshots.\n\nThis needs to be negative.\n\nDon't modify unless you know what you are doing.")]
        // public float catchupNegativeThreshold = -1; // careful, don't want to run out of snapshots
        public float catchupNegativeThreshold = -0.033f; //因为发送间隔设置的1

        [Tooltip("Catchup begins when the local timeline is moving too slow and getting too far away from remote time. Threshold is in frames worth of snapshots.\n\nThis needs to be positive.\n\nDon't modify unless you know what you are doing.")]
        // public float catchupPositiveThreshold = 1;
        public float catchupPositiveThreshold = 0.033f; //因为发送间隔设置的1

        [Tooltip("Local timeline acceleration in % while catching up.")]
        [Range(0, 1)]
        public double catchupSpeed = 0.02f; // see snap interp demo. 1% is too slow.

        [Tooltip("Local timeline slowdown in % while slowing down.")]
        [Range(0, 1)]
        public double slowdownSpeed = 0.04f; // slow down a little faster so we don't encounter empty buffer (= jitter)
        // 
        [Tooltip("Catchup/Slowdown is adjusted over n-second exponential moving average.")]
        // public int driftEmaDuration = 1; // shouldn't need to modify this, but expose it anyway
        public int driftEmaDuration = 30; //因为发送频率是1

        // dynamic buffer time adjustment //////////////////////////////////////
        // dynamically adjusts bufferTimeMultiplier for smooth results.
        // to understand how this works, try this manually:
        //
        // - disable dynamic adjustment
        // - set jitter = 0.2 (20% is a lot!)
        // - notice some stuttering
        // - disable interpolation to see just how much jitter this really is(!)
        // - enable interpolation again
        // - manually increase bufferTimeMultiplier to 3-4
        //   ... the cube slows down (blue) until it's smooth
        // - with dynamic adjustment enabled, it will set 4 automatically
        //   ... the cube slows down (blue) until it's smooth as well
        //
        // note that 20% jitter is extreme.
        // for this to be perfectly smooth, set the safety tolerance to '2'.
        // but realistically this is not necessary, and '1' is enough. @
        [Header("Dynamic Adjustment")]
        [Tooltip("Automatically adjust bufferTimeMultiplier for smooth results.\nSets a low multiplier on stable connections, and a high multiplier on jittery connections.")]
        //public bool dynamicAdjustment = true;
        public bool dynamicAdjustment = false;

        // 
        [Tooltip("Safety buffer that is always added to the dynamic bufferTimeMultiplier adjustment.")]
        public float dynamicAdjustmentTolerance = 1; // 1 is realistically just fine, 2 is very very safe even for 20% jitter. can be half a frame too. (see above comments)

        //
        [Tooltip("Dynamic adjustment is computed over n-second exponential moving average standard deviation.")]
        //public int deliveryTimeEmaDuration = 2;   // 1-2s recommended to capture average delivery time
        public int deliveryTimeEmaDuration = 60;//因为发送频率是1

    }
}
