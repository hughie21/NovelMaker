package logging

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemMssage struct {
	OS     string
	Kernel string
	CPU    string
	Memory uint64
}

func NewSystem() *SystemMssage {
	hostInfo, _ := host.Info()
	cpuinfo, _ := cpu.Info()
	memoryInfo, _ := mem.VirtualMemory()
	return &SystemMssage{
		OS:     hostInfo.Platform,
		Kernel: hostInfo.KernelVersion,
		CPU:    cpuinfo[0].ModelName,
		Memory: memoryInfo.Total,
	}
}

func (s *SystemMssage) String() string {
	return fmt.Sprintf(`
===================System Information====================
SystemOS: %s
Version: %s
CPU: %s
Memory: %d
====================================================
	`, s.OS, s.Kernel, s.CPU, s.Memory)
}

func (s *SystemMssage) getLevel() Level {
	return InfoLevel
}
