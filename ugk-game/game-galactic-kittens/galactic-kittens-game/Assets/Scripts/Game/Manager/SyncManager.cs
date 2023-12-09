using System;
using System.Collections.Generic;
using Common.Network;
using Common.Network.Sync;
using Common.Tools;
using UnityEngine;

namespace Game.Manager
{
    /// <summary>
    /// 同步管理器，保存需要同步管理的对象
    /// <para>客户端玩家自己控制的对象使用快照同步，服务器出生的使用预测同步</para>
    /// </summary>
    public class SyncManager : SingletonPersistent<SyncManager>
    {
        /// <summary>
        /// 场景中所有快照同步对象
        /// </summary>
        private readonly Dictionary<long, SnapTransform> _snapTransforms = new Dictionary<long, SnapTransform>();

        /// <summary>
        /// 场景中所有预测同步对象
        /// </summary>
        private readonly Dictionary<long, PredictionTransform> _predictionTransforms =
            new Dictionary<long, PredictionTransform>();

        /// <summary>
        /// 批量同步快照的消息
        /// </summary>
        private SnapSyncResponse snapSyncMessage;

        /// <summary>
        /// 批量预测同步消息 
        /// </summary>
        private PredictionSyncResponse predictionSyncMessage;


        public override void Awake()
        {
            base.Awake();
            snapSyncMessage = new SnapSyncResponse();
            predictionSyncMessage = new PredictionSyncResponse();
        }

        private void OnEnable()
        {
            ResetData();
        }

        private void OnDisable()
        {
            ResetData();
        }

        private void ResetData()
        {
            _snapTransforms.Clear();
            _predictionTransforms.Clear();
        }

        /// <summary>
        /// 收到同步消息
        /// </summary>
        public void OnSnapSyncReceive(Player player,UgkMessage ugkMessage,SnapSyncRequest request)
        {
            if (!gameObject.activeSelf)
            {
                return;
            }

            try
            {
                foreach (var kv in request.Payload)
                {
                    if (!_snapTransforms.TryGetValue(kv.Key, out SnapTransform snapTransform))
                    {
                        Debug.LogWarning($"同步对象{kv.Key} 不存在");
                        continue;
                    }

                    if (player.Id!=snapTransform.Id)
                    {
                        Log.Warn($"{player.Id} push other people's position");
                        return;
                    }

                    snapTransform.OnDeserialize(ugkMessage,kv.Value, false);
                }
            }
            catch (Exception e)
            {
               Log.Error($"解析快照插值消息错误：{e}");
            }
           
        }

        /// <summary>
        /// 收到同步消息
        /// </summary>
        public void OnPredictionSyncReceive(UgkMessage ugkMessage, PredictionSyncRequest request)
        {
            if (!gameObject.activeSelf)
            {
                return;
            }

            try
            {
                foreach (var kv in request.Payload)
                {
                    if (!_predictionTransforms.TryGetValue(kv.Key, out PredictionTransform predictionTransform))
                    {
                        Debug.LogWarning($"同步对象{kv.Key} 不存在");
                        continue;
                    }

                    predictionTransform.OnDeserialize(ugkMessage, kv.Value, false);
                }
            }
            catch (Exception e)
            {
                Log.Error($"解析预测同步消息错误：{e}");
            }
           
        }


        /// <summary>
        /// 将同步消息发送给玩家
        /// <para>由于所有玩家都在同一个屏幕，所以没有AOI管理，将所有消息同步给所有玩家</para>
        /// </summary>
        private void SyncTransformToPlayers()
        {
            foreach (var kv in _snapTransforms)
            {
                if (kv.Value.SyncData != null)
                {
                    snapSyncMessage.Payload[kv.Key] = kv.Value.SyncData;
                    kv.Value.SyncData = null;
                }
            }

            foreach (var kv in _predictionTransforms)
            {
                if (kv.Value.SyncData != null)
                {
                    predictionSyncMessage.Payload[kv.Key] = kv.Value.SyncData;
                    kv.Value.SyncData = null;
                }
            }

            //批量同步消息
            if (snapSyncMessage.Payload.Count > 0)
            {
                PlayerManager.Instance.BroadcastMsg(MID.SnapSyncRes, snapSyncMessage);
                snapSyncMessage.Payload.Clear();
            }

            var predictionCount = predictionSyncMessage.Payload.Count;
            if (predictionCount > 0)
            {
                PlayerManager.Instance.BroadcastMsg(MID.PredictionSyncRes, predictionSyncMessage);
                if (predictionCount > 64)
                {
                    Debug.LogWarning($"同步消息太多{predictionCount} =>{predictionSyncMessage.Payload.Keys}");
                }

                predictionSyncMessage.Payload.Clear();
            }
        }

        /// <summary>
        /// 移除
        /// </summary>
        /// <param name="id"></param>
        /// <param name="type">0移除，1移除并隐藏，2移除并销毁</param>
        public bool RemovePredictionTransform(long id, int type = 0)
        {
            if (_predictionTransforms.Remove(id, out PredictionTransform predictionTransform))
            {
                if (type == 1)
                {
                    predictionTransform.gameObject.SetActive(false);
                }
                else if (type == 2)
                {
                    Destroy(predictionTransform);
                }

                return true;
            }

            return false;
        }

        /// <summary>
        /// 移除
        /// </summary>
        /// <param name="id"></param>
        /// <param name="type">0移除，1移除并隐藏，2移除并销毁</param>
        public bool RemoveSnapTransform(long id, int type = 0)
        {
            if (_snapTransforms.Remove(id, out SnapTransform snapTransform))
            {
                if (type == 1)
                {
                    snapTransform.gameObject.SetActive(false);
                }
                else if (type == 2)
                {
                    Destroy(snapTransform);
                }

                return true;
            }

            return false;
        }

        /// <summary>
        /// 移除同步对象
        /// </summary>
        /// <param name="id"></param>
        /// <param name="type">0移除，1移除并隐藏，2移除并销毁</param>
        /// <returns></returns>
        public bool RemoveSyncObject(long id, int type = 0)
        {
            if (RemoveSnapTransform(id, type))
            {
                return true;
            }

            return RemovePredictionTransform(id, type);
        }
        
        public void AddSnapTransform(SnapTransform snapTransform)
        {
            _snapTransforms[snapTransform.Id] = snapTransform;
        }


        public void AddPredictionTransform(PredictionTransform predictionTransform)
        {
            _predictionTransforms[predictionTransform.Id] = predictionTransform;
        }


        public void Update()
        {
            SyncTransformToPlayers();
        }
    }
}