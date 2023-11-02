package manager

import (
	"context"
	"fmt"
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
	MachineInfos         []*message.HostMachineInfo    //主机信息
	RequestChan          chan func()                   //请求处理
	MachiInfoChan        chan *message.HostMachineInfo //上传主机信息chan
	gameAgentGrpcClients map[string]*grpc.ClientConn   //游戏服务器所在的agent key：游戏名称-游戏id,需要存储数据库，服务重启丢失了
}

var dockerManager *DockerManager
var dockerSingletonOnce sync.Once

func GetDockerManager() *DockerManager {
	dockerSingletonOnce.Do(func() {
		dockerManager = &DockerManager{
			MachineInfos:         make([]*message.HostMachineInfo, 0, 5),
			MachiInfoChan:        make(chan *message.HostMachineInfo, 1024),
			RequestChan:          make(chan func(), 1024),
			gameAgentGrpcClients: make(map[string]*grpc.ClientConn, 100),
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
		case requestHand := <-m.RequestChan:
			requestHand()
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
		if info.GetAvailableMemorySize() > 500 && info.GetCpuPercent() < 80 && info.GetAvailableDiskSize() > 1000 {
			machineInfo = info
		}
	}

	if machineInfo == nil {
		log.Warn("没有满足条件的可用主机\n%+v", m.MachineInfos)
		return nil
	}
	client, err := manager.GetServiceClientManager().GetGrpc(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentName, 0), machineInfo.GetServerId())
	if err != nil {
		log.Error("获取agent客户端错误：%v", err)
	}
	return client
}

// CreateGameService 创建游戏服务,转发到合适的agent进行执行
func (m *DockerManager) CreateGameService(ctx context.Context, wg *sync.WaitGroup, request *message.CreateGameServiceRequest, response *message.CreateGameServiceResponse) {
	defer wg.Done()
	grpcClient := m.GetBestAgentClient()
	if grpcClient == nil {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    "无可用主机",
		}
		wg.Done()
		return
	}
	go func(ctx2 context.Context, wg2 *sync.WaitGroup, c *grpc.ClientConn) {
		defer wg2.Done()
		client := message.NewAgentServiceClient(c)
		r, err := client.CreateGameService(ctx2, request)
		if err != nil {
			response.Result = &message.MessageResult{
				Status: 500,
				Msg:    err.Error(),
			}
			log.Error("%v:请求Agent创建服务错误：%v", request, err)
		} else {
			response.Result = r.Result
			m.gameAgentGrpcClients[fmt.Sprintf("%v-%v", request.GetGameName(), request.GetGameId())] = c
		}
		log.Info("%v-%v:创建游戏服务：%v", request.GetGameName(), request.GetGameId(), response)
	}(ctx, wg, grpcClient)
}

// CloseGameService 创建游戏服务,转发到合适的agent进行执行
func (m *DockerManager) CloseGameService(ctx context.Context, wg *sync.WaitGroup, request *message.CloseGameServiceRequest, response *message.CloseGameServiceResponse) {
	defer wg.Done()
	grpcClient := m.gameAgentGrpcClients[fmt.Sprintf("%v-%v", request.GetGameName(), request.GetGameId())]

	if grpcClient == nil {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    "未找到创建服务的agent",
		}
		wg.Done()
		return
	}
	go func(ctx2 context.Context, wg2 *sync.WaitGroup, c *grpc.ClientConn) {
		defer wg2.Done()
		client := message.NewAgentServiceClient(c)
		r, err := client.CloseGameService(ctx2, request)
		if err != nil {
			response.Result = &message.MessageResult{
				Status: 500,
				Msg:    err.Error(),
			}
			log.Error("%v:请求Agent关闭服务错误：%v", request, err)
		} else {
			response.Result = r.Result
		}
		log.Info("%v-%v:关闭游戏服务：%v", request.GetGameName(), request.GetGameId(), response)
	}(ctx, wg, grpcClient)
}
