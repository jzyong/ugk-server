package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"testing"
	"time"
)

// 查询主机信息
func TestQueryMachineInfo(t *testing.T) {

	for i := 0; i < 100; i++ {
		//cpu
		perCpuPercents, _ := cpu.Percent(0, true)
		var cpuPercent float64
		for _, percent := range perCpuPercents {
			cpuPercent += percent
		}
		cpuPercent = cpuPercent / float64(len(perCpuPercents))
		log.Info("CPU百分比：%.2f", cpuPercent)

		//内存
		virtualMemory, _ := mem.VirtualMemory()
		log.Info("剩余内存：%vM", virtualMemory.Available/util.MB)
		log.Info("内存百分比：%.2f", virtualMemory.UsedPercent)

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
		log.Info("剩余磁盘：%vM", freeDiskSize/util.MB)

		time.Sleep(time.Second * 3)
	}

}
