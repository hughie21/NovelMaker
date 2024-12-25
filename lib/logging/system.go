// Description: Get system hardware information, such as memory, GPU, etc.
// Author: Hughie21
// Date: 2024-12-20
// license that can be found in the LICENSE file.
package logging

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// SystemMssage is a struct that contains the system information
type SystemMssage struct {
	OS     string
	Kernel string
	CPU    string
	Memory uint64
}

// NewSystem creates a new SystemMssage object and returns a pointer to it.
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

// Output system information
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

// Get the log level of the system information
func (s *SystemMssage) getLevel() Level {
	return InfoLevel
}
