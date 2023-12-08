using Common.Network;
using Game.Manager;
using Google.Protobuf;

namespace Game.Handlers
{
    /// <summary>
    /// 同步消息处理器
    /// </summary>
    internal class SyncHandler
    {
        /// <summary>
        /// 快照插值同步
        /// </summary>
        [MessageMap((int)MID.SnapSyncReq)]
        private static void SnapSync(Player player, UgkMessage ugkMessage)
        {
            var request = new SnapSyncRequest();
            request.MergeFrom(ugkMessage.Bytes);
            SyncManager.Instance.OnSnapSyncReceive(player, ugkMessage, request);
        }

        /// <summary>
        /// 快照插值同步
        /// </summary>
        [MessageMap((int)MID.PredictionSyncReq)]
        private static void PredictionSync(Player player, UgkMessage ugkMessage)
        {
            var request = new PredictionSyncRequest();
            request.MergeFrom(ugkMessage.Bytes);
            SyncManager.Instance.OnPredictionSyncReceive(ugkMessage, request);
        }
    }
}