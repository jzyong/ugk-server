using System;
using System.Collections;
using System.Collections.Generic;
using Common.Network.Sync;
using Common.Tools;
using Game.Messages;
using Game.Room.Boss;
using Game.Room.Enemy;
using Game.Room.Player;
using UnityEngine;
using Random = UnityEngine.Random;

namespace Game.Manager
{
    /// <summary>
    /// 房间管理
    /// </summary>
    public class RoomManager : SingletonPersistent<RoomManager>
    {
        [SerializeField] [Tooltip("飞船，服务器只需要一个简单对象即可")]
        private SpaceShip _spaceShipPrefab;

        [SerializeField] [Tooltip("子弹")] private SpaceshipBullet _spaceshipBulletPrefab;
        [SerializeField] [Tooltip("敌人子弹")] private EnemyBullet _enemyBulletPrefab;
        [SerializeField] [Tooltip("Boss小子弹")] private BossSmallBullet _bossSmallBulletPrefab;
        [SerializeField] [Tooltip("Boss环形子弹")] private BossCircularBullet _bossCircularBulletPrefab;

        [SerializeField] [Tooltip("Boss自动跟踪导弹")]
        private BossHomingMisile _bossHomingMisilePrefab;


        [SerializeField] [Tooltip("不射击的敌人")] private SpaceGhostEnemy _spaceGhostEnemyPrefab;
        [SerializeField] [Tooltip("击的敌人")] private SpaceShooterEnemy _spaceShooterEnemyPrefab;
        [SerializeField] [Tooltip("陨石")] private Meteor _meteorPrefab;
        [SerializeField] [Tooltip("Boss")] private Boss _bossPrefab;
        [SerializeField] [Tooltip("能量道具")] private PowerUp _powerUpPrefab;


        [Header("Enemies")] [SerializeField] private float m_EnemySpawnTime = 1.8f;
        [SerializeField] private float m_meteorSpawningTime = 1f;
        [SerializeField] private float m_bossSpawnTime = 45;
        private Vector3 m_CurrentNewEnemyPosition = new Vector3();
        private float m_CurrentEnemySpawnTime = 0f;
        private Vector3 m_CurrentNewMeteorPosition = new Vector3();
        private float m_CurrentMeteorSpawnTime = 0f;
        private float m_CurrentBossSpawnTime = 0f;
        private bool m_IsSpawning = true;
        private long bossId;


        /// <summary>
        /// 飞船出生坐标
        /// </summary>
        private readonly Vector3[] shipSpawnPositions = new[]
            { new Vector3(-8, 4), new Vector3(-8, 1.5f), new Vector3(-8, -1f), new Vector3(-8, -3.5f) };

        private Dictionary<long, SpaceShip> _spaceShips = new Dictionary<long, SpaceShip>(4);


        /// <summary>
        /// 对象同步Id
        /// </summary>
        private long SyncId { get; set; }

        private void Start()
        {
            // Initialize the enemy and meteor spawn position based on my owning GO's x position
            m_CurrentNewEnemyPosition.x = 9;
            m_CurrentNewEnemyPosition.z = 0f;

            m_CurrentNewMeteorPosition.x = 10; //客户端对象出生计算位置时会快速移动一小段，因此从屏幕外出生
            m_CurrentNewMeteorPosition.z = 0f;
            SyncId = short.MaxValue; //防止和玩家id冲突
        }

        private void Update()
        {
            if (m_IsSpawning)
            {
                SpawnEnemy();
                SpawnMeteor();
                SpawnBoss();
            }
        }

        /// <summary>
        /// 出生玩家
        /// </summary>
        public void SpawnPlayers(List<Player> players)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();
            for (int i = 0; i < players.Count; i++)
            {
                // 根据角色创建对应的实体对象，添加SnapTransform组件
                var player = players[i];
                var spawnPosition = shipSpawnPositions[i];
                var spaceShip = Instantiate(_spaceShipPrefab, spawnPosition, Quaternion.identity,
                    Instance.transform);
                var snapTransform = spaceShip.GetComponent<SnapTransform>();
                snapTransform.Id = player.Id;
                if (player.Id > SyncId) //防止ID重复
                {
                    SyncId = player.Id;
                }

                GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                    new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                    {
                        OwnerId = player.Id,
                        Id = player.Id,
                        ConfigId = (uint)player.CharacterId, //0-3
                        SyncInterval = snapTransform.sendInterval,
                        Position = ProtoUtil.BuildVector3D(spawnPosition)
                    };
                snapTransform.InitTransform(spawnPosition, null);
                SyncManager.Instance.AddSnapTransform(snapTransform); //添加同步对象
                spawnResponse.Spawn.Add(spawnInfo);
                _spaceShips[player.Id] = spaceShip;
                Log.Info($"{player.Id} born in {spawnPosition}");
            }

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }

        /// <summary>
        /// 出生敌人
        /// </summary>
        private void SpawnEnemy()
        {
            m_CurrentEnemySpawnTime += Time.deltaTime;
            if (m_CurrentEnemySpawnTime >= m_EnemySpawnTime)
            {
                // update the new enemy's spawn position(y value). This way we don't have to allocate
                // a new Vector3 each time.
                m_CurrentNewEnemyPosition.y = Random.Range(-5f, 5f);
                int randomPick = Random.Range(0, 99);

                GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();


                GameObject gameObject;


                uint configId = 40; //40射击敌人、41幽灵敌人、41陨石
                //射击
                if (randomPick < 50)
                {
                    gameObject = Instantiate(_spaceShooterEnemyPrefab, m_CurrentNewEnemyPosition,
                        Quaternion.identity,
                        Instance.transform).gameObject;
                }
                else
                {
                    gameObject = Instantiate(_spaceGhostEnemyPrefab, m_CurrentNewEnemyPosition,
                        Quaternion.identity,
                        Instance.transform).gameObject;
                    configId = 41;
                }

                var snapTransform = gameObject.GetComponent<SnapTransform>();
                snapTransform.Id = SyncId++;
                snapTransform.Onwer = true;
                gameObject.name = $"SpaceEnemy-{snapTransform.Id}";
                GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                    new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                    {
                        OwnerId = 0,
                        Id = snapTransform.Id,
                        ConfigId = configId,
                        SyncInterval = snapTransform.sendInterval,
                        Position = ProtoUtil.BuildVector3D(m_CurrentNewEnemyPosition),
                    };
                snapTransform.InitTransform(m_CurrentNewEnemyPosition, null);
                SyncManager.Instance.AddSnapTransform(snapTransform); //添加同步对象
                spawnResponse.Spawn.Add(spawnInfo);
                Log.Info($"enemy {snapTransform.Id}  born in {m_CurrentNewEnemyPosition}");

                PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);

                m_CurrentEnemySpawnTime = 0f;
            }
        }

        /// <summary>
        ///  出生能量提升
        /// </summary>
        public void SpawnPowerUp (Vector3 position)
        {
            //概率控制
            int randomPick = Random.Range(1, 100);
            if (randomPick>10)
            {
                return;
            }
            
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();
            PowerUp powerUp = Instantiate(_powerUpPrefab, position, Quaternion.identity,
                Instance.transform);
            uint configId = 50;
            var predictionTransform = powerUp.GetComponent<PredictionTransform>();
            predictionTransform.Id = SyncId++;
            predictionTransform.Onwer = true;
            predictionTransform.LinearVelocity = powerUp.linearVelocity;
            powerUp.name = $"PowerUp-{predictionTransform.Id}";
            GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                {
                    OwnerId = 0,
                    Id = predictionTransform.Id,
                    ConfigId = configId,
                    Position = ProtoUtil.BuildVector3D(m_CurrentNewMeteorPosition),
                    LinearVelocity = ProtoUtil.BuildVector3D(predictionTransform.LinearVelocity),
                };
            SyncManager.Instance.AddPredictionTransform(predictionTransform); //添加同步对象
            spawnResponse.Spawn.Add(spawnInfo);
            Log.Info($"PowerUp {predictionTransform.Id}  born in {m_CurrentNewMeteorPosition}");

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);

            m_CurrentMeteorSpawnTime = 0f;
        }


        /// <summary>
        /// 出生陨石
        /// </summary>
        private void SpawnMeteor()
        {
            m_CurrentMeteorSpawnTime += Time.deltaTime;
            if (m_CurrentMeteorSpawnTime >= m_meteorSpawningTime)
            {
                m_CurrentNewMeteorPosition.y = Random.Range(-5f, 6f);

                GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();


                Meteor meteor = Instantiate(_meteorPrefab, m_CurrentNewMeteorPosition, Quaternion.identity,
                    Instance.transform);
                meteor.SpawnInit();
                uint configId = 42; //40射击敌人、41幽灵敌人、42陨石

                var predictionTransform = meteor.GetComponent<PredictionTransform>();
                predictionTransform.Id = SyncId++;
                predictionTransform.Onwer = true;
                predictionTransform.LinearVelocity = Vector3.left * 4;
                meteor.name = $"Meteor-{predictionTransform.Id}";
                GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                    new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                    {
                        OwnerId = 0,
                        Id = predictionTransform.Id,
                        ConfigId = configId,
                        Position = ProtoUtil.BuildVector3D(m_CurrentNewMeteorPosition),
                        LinearVelocity = ProtoUtil.BuildVector3D(predictionTransform.LinearVelocity),
                        Scale = ProtoUtil.BuildVector3D(meteor.transform.localScale)
                    };
                SyncManager.Instance.AddPredictionTransform(predictionTransform); //添加同步对象
                spawnResponse.Spawn.Add(spawnInfo);
                Log.Info($"meteor {predictionTransform.Id}  born in {m_CurrentNewMeteorPosition}");

                PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);

                m_CurrentMeteorSpawnTime = 0f;
            }
        }


        /// <summary>
        /// 出生Boss
        /// </summary>
        private void SpawnBoss()
        {
            m_CurrentBossSpawnTime += Time.deltaTime;
            if (m_CurrentBossSpawnTime >= m_bossSpawnTime)
            {
                m_IsSpawning = false;
                StartCoroutine(BossAppear());
            }
        }

        IEnumerator BossAppear()
        {
            // Warning title and sound
            SpawnBoss(20);

            // Same time as audio length
            yield return new WaitForSeconds(5.5f);
            SpawnBoss(21);
        }

        /// <summary>
        /// 
        /// </summary>
        /// <param name="type">20Boss预警，21Boss</param>
        public void SpawnBoss(uint type)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();
            GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                new GalacticKittensObjectSpawnResponse.Types.SpawnInfo();
            spawnInfo.ConfigId = type;
            spawnResponse.Spawn.Add(spawnInfo);
            if (type == 20)
            {
                PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
                return;
            }

            var spawnPosition = new Vector3(5, 0, 0);
            Boss boss = Instantiate(_bossPrefab, spawnPosition, Quaternion.identity,
                Instance.transform);
            boss.StartBoss(spawnPosition);

            var snapTransform = boss.GetComponent<SnapTransform>();
            snapTransform.Id = SyncId++;
            bossId = snapTransform.Id;
            snapTransform.Onwer = true;
            snapTransform.InitTransform(spawnPosition, null);
            boss.name = $"Boss-{snapTransform.Id}";
            spawnInfo.OwnerId = 0;
            spawnInfo.Id = snapTransform.Id;
            spawnInfo.Position = ProtoUtil.BuildVector3D(spawnPosition);
            SyncManager.Instance.AddSnapTransform(snapTransform); //添加同步对象
            Log.Info($"Boss {snapTransform.Id}  born in {m_CurrentNewMeteorPosition}");
            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }


        /// <summary>
        /// 创建子弹
        /// </summary>
        /// <param name="player"></param>
        public void SpawnBullet(Player player)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();
            SpaceShip spaceShip = _spaceShips[player.Id];

            var spawnPosition = spaceShip.transform.position;
            spawnPosition = new Vector3(spawnPosition.x + 1f, spawnPosition.y - 0.3f, spawnPosition.z); //y轴下移一点
            var spaceshipBullet = Instantiate(_spaceshipBulletPrefab, spawnPosition, Quaternion.identity,
                Instance.transform);
            var predictionTransform = spaceshipBullet.GetComponent<PredictionTransform>();
            predictionTransform.Id = SyncId++;
            spaceshipBullet.name = $"SpaceBullet{player.Id}-{predictionTransform.Id}";
            predictionTransform.LinearVelocity = spaceshipBullet.linearVelocity;
            GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                {
                    OwnerId = player.Id,
                    Id = predictionTransform.Id,
                    ConfigId = 30,
                    Position = ProtoUtil.BuildVector3D(spawnPosition),
                    LinearVelocity = ProtoUtil.BuildVector3D(spaceshipBullet.linearVelocity),
                };
            SyncManager.Instance.AddPredictionTransform(predictionTransform); //添加同步对象
            spawnResponse.Spawn.Add(spawnInfo);
            Log.Info($"{player.Id} bullet born in {spawnPosition}");

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }

        /// <summary>
        /// 创建敌人子弹
        /// </summary>
        /// <param name="enemy"></param>
        public void SpawnEnemyBullet(SpaceShooterEnemy enemy)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();

            var spawnPosition = enemy.transform.position;
            var bullet = Instantiate(_enemyBulletPrefab, spawnPosition, Quaternion.identity,
                Instance.transform);
            var predictionTransform = bullet.GetComponent<PredictionTransform>();
            predictionTransform.Id = SyncId++;
            bullet.name = $"EnemyBullet-{predictionTransform.Id}";
            predictionTransform.LinearVelocity = bullet.linearVelocity;
            GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                {
                    OwnerId = enemy.GetComponent<SnapTransform>().Id,
                    Id = predictionTransform.Id,
                    ConfigId = 31,
                    Position = ProtoUtil.BuildVector3D(spawnPosition),
                    LinearVelocity = ProtoUtil.BuildVector3D(bullet.linearVelocity),
                };
            SyncManager.Instance.AddPredictionTransform(predictionTransform); //添加同步对象
            spawnResponse.Spawn.Add(spawnInfo);
            Log.Info($"enemy bullet born in {spawnPosition}");

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }

        /// <summary>
        /// 产生boss子弹
        /// </summary>
        /// <param name="type">32 boss三角形小子弹，33 boss环形分裂后小子弹，34 boss环形分裂子弹，35 boss导弹</param>
        /// <param name="position"></param>
        /// <param name="rotation"></param>
        public void SpawnBossBullet(uint type, Vector3 position, Vector3 rotation)
        {
            GalacticKittensObjectSpawnResponse spawnResponse = new GalacticKittensObjectSpawnResponse();

            GameObject bullet = null;
            switch (type)
            {
                case 32:
                case 33:
                    bullet = Instantiate(_bossSmallBulletPrefab, position, Quaternion.Euler(rotation),
                        Instance.transform).gameObject;
                    break;
                case 34:
                    bullet = Instantiate(_bossCircularBulletPrefab, position, Quaternion.Euler(rotation),
                        Instance.transform).gameObject;
                    break;
                case 35:
                    bullet = Instantiate(_bossHomingMisilePrefab, position, _bossHomingMisilePrefab.transform.rotation,
                        Instance.transform).gameObject;
                    break;
                default:
                    Log.Warn($"bullet type：{type} not find");
                    return;
            }

            var snapTransform = bullet.GetComponent<SnapTransform>();
            snapTransform.Id = SyncId++;
            snapTransform.Onwer = true;
            bullet.name = $"BossBullet-{snapTransform.Id}";
            snapTransform.InitTransform(position, null);
            GalacticKittensObjectSpawnResponse.Types.SpawnInfo spawnInfo =
                new GalacticKittensObjectSpawnResponse.Types.SpawnInfo()
                {
                    OwnerId = bossId,
                    Id = snapTransform.Id,
                    ConfigId = type,
                    Position = ProtoUtil.BuildVector3D(position),
                };
            SyncManager.Instance.AddSnapTransform(snapTransform); //添加同步对象
            spawnResponse.Spawn.Add(spawnInfo);
            Log.Info($"Boss bullet born in {position}");
            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectSpawnRes, spawnResponse);
        }


        /// <summary>
        /// 对象死亡
        /// </summary>
        public void DespawnObject(long killerId, long dieId, bool removeObject = true)
        {
            if (removeObject)
            {
                //移除的对象全部使用预测
                SyncManager.Instance.RemoveSyncObject(dieId);
            }

            GalacticKittensObjectDieResponse response = new GalacticKittensObjectDieResponse()
            {
                KillerId = killerId,
                Id = dieId,
                // OwnerId = //TOD
            };

            PlayerManager.Instance.BroadcastMsg(MID.GalacticKittensObjectDieRes, response);
        }


        /// <summary>
        /// 游戏结束 TODO 待测试
        /// </summary>
        private void GameFinishReq()
        {
            var client =
                new GalacticKittensMatchService.GalacticKittensMatchServiceClient(GalacticKittensNetworkManager
                    .Instance.MatchChannel);
            var request = new GalacticKittensGameFinishRequest()
            {
                RoomId = GalacticKittensNetworkManager.Instance.ServerId
            };
            var response = client.gameFinishAsync(request).ResponseAsync.Result;
            Log.Info($"game finish :{response}");
            if (response.Result?.Status != 200)
            {
                Log.Error($"game finish error:{response.Result?.Msg}");
                return;
            }

            // 解绑玩家网关映射
            PlayerManager.Instance.BindGateGameMapReq(false);

            Application.Quit();
        }

        public SpaceShip GetSpaceShip(long id)
        {
            if (_spaceShips.TryGetValue(id, out SpaceShip spaceShip))
            {
                return spaceShip;
            }

            Log.Warn($"ship {id} not find");
            return null;
        }

        public int PlayerCount()
        {
            return _spaceShips.Count;
        }

        public void GameFinish()
        {
            SyncManager.Instance.ResetData();
        }
    }
}