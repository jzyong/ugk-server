package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
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
	log.Debug("CPU百分比：%.2f", cpuPercent)

	//内存
	virtualMemory, _ := mem.VirtualMemory()
	log.Debug("剩余内存：%vM", virtualMemory.Available/util.MB)
	log.Debug("内存百分比：%.2f", virtualMemory.UsedPercent)

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
	log.Debug("剩余磁盘：%vM", freeDiskSize/util.MB)
}
