package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/grpc"
	"sort"
	"sync"
	"time"
)

// DockerManager 管理agent创建docker容器服务,需要多台主机测试
type DockerManager struct {
	util.DefaultModule
	MachineInfos  []*message.HostMachineInfo    //主机信息
	MachiInfoChan chan *message.HostMachineInfo //上传主机信息chan
}

var dockerManager *DockerManager
var dockerSingletonOnce sync.Once

func GetDockerManager() *DockerManager {
	dockerSingletonOnce.Do(func() {
		dockerManager = &DockerManager{
			MachineInfos:  make([]*message.HostMachineInfo, 0, 5),
			MachiInfoChan: make(chan *message.HostMachineInfo, 1024),
		}
	})
	return dockerManager
}

func (m *DockerManager) Init() error {
	log.Info("DockerManager 初始化......")
	return nil
}
func (m *DockerManager) Run() {
	go m.run()
}

func (m *DockerManager) Stop() {
}

// 单独一个routine执行主机分配
func (m *DockerManager) run() {
	minuteTicker := time.Tick(time.Minute)
	for {
		select {
		case machineInfo := <-m.MachiInfoChan:
			m.updateMachineInfo(machineInfo)
		case <-minuteTicker:
			m.minuteUpdate()
		}
	}
}

// 更新主机信息
func (m *DockerManager) updateMachineInfo(machineInfo *message.HostMachineInfo) {
	info := m.GetMachineInfo(machineInfo.GetServerId())
	if info == nil {
		m.MachineInfos = append(m.MachineInfos, machineInfo)
	} else {
		info.CpuPercent = machineInfo.CpuPercent
		info.MemoryPercent = machineInfo.MemoryPercent
		info.AvailableMemorySize = machineInfo.AvailableMemorySize
		info.AvailableDiskSize = machineInfo.AvailableDiskSize
	}

}

func (m *DockerManager) minuteUpdate() {
	//对主机进行排序,安装内存进行倒序排
	sort.Slice(m.MachineInfos, func(i, j int) bool {
		return m.MachineInfos[i].GetAvailableMemorySize() > m.MachineInfos[j].GetAvailableMemorySize()
	})
}

func (m *DockerManager) GetMachineInfo(id uint32) *message.HostMachineInfo {
	for _, info := range m.MachineInfos {
		if info.GetServerId() == id {
			return info
		}
	}
	return nil
}

// GetBestAgentClient 获得运行Docker容器的客户端
func (m *DockerManager) GetBestAgentClient() *grpc.ClientConn {
	var machineInfo *message.HostMachineInfo
	for _, info := range m.MachineInfos {
		if info.GetAvailableMemorySize() > 500 && info.GetCpuPercent() < 0.8 && info.GetAvailableDiskSize() > 1000 {
			machineInfo = info
		}
	}

	if machineInfo == nil {
		log.Warn("没有满足条件的可用主机\n%+v", m.MachineInfos)
		return nil
	}
	client, err := manager.GetServiceClientManager().GetGrpcConn(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentName, 0), machineInfo.GetServerId())
	if err != nil {
		log.Error("获取agent客户端错误：%v", err)
	}
	return client
}
