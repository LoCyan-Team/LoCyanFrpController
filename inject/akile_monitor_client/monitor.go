package akile_monitor_client

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
	"lcf-controller/inject/akile_monitor_client/model"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	"runtime"
	"strconv"
	"time"
)

var cfg = config.ReadCfg().MonitorConfig

func GetState() *model.HostState {
	var ret model.HostState
	cp, err := cpu.Percent(0, false)
	if err != nil || len(cp) == 0 {
		logger.Warn("cpu.Percent error", zap.Error(err))
	} else {
		ret.CPU = cp[0]
	}

	loadStat, err := load.Avg()
	if err != nil {
		logger.Warn("load.Avg error", zap.Error(err))
	} else {
		ret.Load1 = Decimal(loadStat.Load1)
		ret.Load5 = Decimal(loadStat.Load5)
		ret.Load15 = Decimal(loadStat.Load15)

	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		logger.Warn("mem.VirtualMemory error", zap.Error(err))
	} else {
		ret.MemUsed = vm.Total - vm.Available
	}

	uptime, err := host.Uptime()
	if err != nil {
		logger.Warn("host.Uptime error", zap.Error(err))
	} else {
		ret.Uptime = uptime
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		logger.Warn("mem.SwapMemory error", zap.Error(err))
	} else {
		ret.SwapUsed = swap.Used
	}

	ret.NetInTransfer, ret.NetOutTransfer = netInTransfer, netOutTransfer
	ret.NetInSpeed, ret.NetOutSpeed = netInSpeed, netOutSpeed

	return &ret

}

func GetHost() *model.Host {
	var ret model.Host
	ret.Name = cfg.Name
	var cpuType string
	hi, err := host.Info()
	if err != nil {
		logger.Warn("host.Info error", zap.Error(err))
	}
	cpuType = "Virtual"
	ret.Platform = hi.Platform
	ret.PlatformVersion = hi.PlatformVersion
	ret.Arch = hi.KernelArch
	ret.Virtualization = hi.VirtualizationSystem
	ret.BootTime = hi.BootTime

	// 检查是否在 Docker 环境中
	if ret.Virtualization == "docker" {
		// 创建 Docker 客户端
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			logger.Warn("failed to create Docker client", zap.Error(err))
		} else {
			defer cli.Close()

			// 获取 Docker 信息
			info, err := cli.Info(context.Background())
			if err != nil {
				logger.Warn("Failed to get Docker info", zap.Error(err))
			} else {
				// 更新宿主机信息
				ret.Platform = info.OperatingSystem
				ret.PlatformVersion = ""
				ret.Arch = info.Architecture
			}
		}
	}

	ci, err := cpu.Info()
	if err != nil {
		logger.Warn("cpu.Info error", zap.Error(err))
	}
	ret.CPU = append(ret.CPU, fmt.Sprintf("%s %d %s Core", ci[0].ModelName, runtime.NumCPU(), cpuType))
	vm, err := mem.VirtualMemory()
	if err != nil {
		logger.Warn("mem.VirtualMemory error", zap.Error(err))
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		logger.Warn("mem.SwapMemory error", zap.Error(err))
	}

	ret.MemTotal = vm.Total
	ret.SwapTotal = swap.Total
	return &ret

}

var (
	netInSpeed, netOutSpeed, netInTransfer, netOutTransfer, lastUpdateNetStats uint64
)

// TrackNetworkSpeed NIC监控，统计流量与速度
func TrackNetworkSpeed() {
	var innerNetInTransfer, innerNetOutTransfer uint64
	nc, err := net.IOCounters(true)
	if err == nil {
		for _, v := range nc {
			if v.Name == cfg.NetworkDevice {
				innerNetInTransfer += v.BytesRecv
				innerNetOutTransfer += v.BytesSent
			}
		}
		now := uint64(time.Now().Unix())
		diff := now - lastUpdateNetStats
		if diff > 0 {
			netInSpeed = (innerNetInTransfer - netInTransfer) / diff
			netOutSpeed = (innerNetOutTransfer - netOutTransfer) / diff
		}
		netInTransfer = innerNetInTransfer
		netOutTransfer = innerNetOutTransfer
		lastUpdateNetStats = now

	}
}

// 保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
