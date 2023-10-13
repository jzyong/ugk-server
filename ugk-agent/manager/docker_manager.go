package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent/config"
	"github.com/jzyong/ugk/message/message"
	"os/exec"
	"sync"
)

// DockerManager 执行系统命令在docker中运行unity服务器
type DockerManager struct {
	util.DefaultModule
}

var dockerManager *DockerManager
var dockerSingletonOnce sync.Once

func GetDockerManager() *DockerManager {
	dockerSingletonOnce.Do(func() {
		dockerManager = &DockerManager{}
	})
	return dockerManager
}

func (m *DockerManager) Init() error {
	log.Info("DockerManager 初始化......")
	return nil
}
func (m *DockerManager) Run() {
	go m.showContainers()
}

// CreateGameServiceContainer 创建游戏服务容器
func (m *DockerManager) CreateGameServiceContainer(request *message.CreateGameServiceRequest) *message.CreateGameServiceResponse {
	//set UnityParam="grpcUrl=192.168.110.2:4000 serverId=1"
	//docker run -dit --name game-galactic-kittens-1 -e UnityParam=%UnityParam% game-galactic-kittens:develop
	containerName := fmt.Sprintf("%v-%d", request.GetGameName(), request.GetGameId())
	UnityParam := fmt.Sprintf("grpcUrl=%v serverId=%d", request.GetControlGrpcUrl(), request.GetGameId())
	UnityParam = fmt.Sprintf("UnityParam=%v", UnityParam)
	containerImage := fmt.Sprintf("%v:%v", request.GetGameName(), config.BaseConfig.Profile)
	cmd := exec.Command("docker", "run", "-dit", "--name", containerName, "-e", UnityParam, containerImage)
	output, err := cmd.CombinedOutput()
	response := &message.CreateGameServiceResponse{}
	if err != nil {
		log.Error("执行错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    fmt.Sprintf("create container fail:%v", err),
		}
	}
	log.Info("运行结果：%v", string(output))
	response.Result = &message.MessageResult{
		Status: 200,
		Msg:    string(output),
	}
	return response
}

// CloseGameServiceContainer 关闭游戏服务容器
func (m *DockerManager) CloseGameServiceContainer(request *message.CloseGameServiceRequest) *message.CloseGameServiceResponse {
	//docker stop game-galactic-kittens-1
	//docker rm  game-galactic-kittens-1
	containerName := fmt.Sprintf("%v-%d", request.GetGameName(), request.GetGameId())
	cmd := exec.Command("docker", "stop", containerName)
	output, err := cmd.CombinedOutput()
	response := &message.CloseGameServiceResponse{}
	if err != nil {
		log.Error("执行错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    fmt.Sprintf("stop container fail:%v", err),
		}
	}
	log.Info("运行结果：%v", string(output))

	cmd = exec.Command("docker", "rm", containerName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    fmt.Sprintf("rm container fail:%v", err),
		}
	}
	log.Info("运行结果：%v", string(output))

	response.Result = &message.MessageResult{
		Status: 200,
		Msg:    string(output),
	}
	return response
}

func (m *DockerManager) showContainers() {
	//执行docker ps
	cmd := exec.Command("docker", "ps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
		return
	}
	log.Info("本机Docker容器：\n%v", string(output))
}

func (m *DockerManager) Stop() {
}
