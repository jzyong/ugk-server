package manager

import (
	"context"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent/config"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/message/message"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"sync"
	"time"
)

// MachineManager 主机信息
type MachineManager struct {
	util.DefaultModule
}

var machineManager *MachineManager
var machineSingletonOnce sync.Once

func GetMachineManager() *MachineManager {
	machineSingletonOnce.Do(func() {
		machineManager = &MachineManager{}
	})
	return machineManager
}

func (m *MachineManager) Init() error {
	log.Info("MachineManager 初始化......")
	return nil
}
func (m *MachineManager) Run() {

	go m.run()
}

func (m *MachineManager) Stop() {
}

func (m *MachineManager) run() {
	uploadMachineTicker := time.Tick(time.Second * 3)
	for {
		select {
		case <-uploadMachineTicker:
			m.uploadMachineInfo()
		}
	}
}

// 上传主机信息
func (m *MachineManager) uploadMachineInfo() {
	//cpu
	perCpuPercents, _ := cpu.Percent(0, true)
	var cpuPercent float64
	for _, percent := range perCpuPercents {
		cpuPercent += percent
	}
	cpuPercent = cpuPercent / float64(len(perCpuPercents))
	log.Trace("CPU百分比：%.2f", cpuPercent)

	//内存
	virtualMemory, _ := mem.VirtualMemory()
	availableMemorySize := virtualMemory.Available / util.MB
	log.Trace("剩余内存：%vM", availableMemorySize)
	log.Trace("内存百分比：%.2f", virtualMemory.UsedPercent)

	//磁盘
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get Partitions failed, err:%v\n", err)
		return
	}
	var freeDiskSize uint64
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		freeDiskSize += diskInfo.Free
	}
	freeDiskSize = freeDiskSize / util.MB
	log.Trace("剩余磁盘：%vM", freeDiskSize)

	client, err := manager.GetServiceClientManager().GetGrpc(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentManagerName, 0), 0)
	if err != nil {
		log.Error("获取agent-manager失败：%v", err)
		return
	}
	stub := message.NewAgentControlServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	machineInfo := &message.HostMachineInfo{
		CpuPercent:          float32(cpuPercent),
		MemoryPercent:       float32(virtualMemory.UsedPercent),
		AvailableMemorySize: float32(availableMemorySize),
		AvailableDiskSize:   float32(freeDiskSize),
		ServerId:            config.BaseConfig.Id,
	}

	response, err := stub.HostMachineInfoUpload(ctx, &message.HostMachineInfoUploadRequest{HostMachineInfo: machineInfo})
	if err != nil {
		log.Error("上传主机信息失败：%v", err)
	}
	log.Trace("上传结果：%v", response)

}
