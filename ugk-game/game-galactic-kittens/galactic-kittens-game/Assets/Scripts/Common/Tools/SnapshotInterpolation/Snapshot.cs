namespace Common.Tools.SnapshotInterpolation
{
    public interface Snapshot
    {
        // the remote timestamp (when it was sent by the remote)
        double remoteTime { get; set; }

        // the local timestamp (when it was received on our end)
        // technically not needed for basic snapshot interpolation.
        // only for dynamic buffer time adjustment. 
        double localTime { get; set; }
    }
}
