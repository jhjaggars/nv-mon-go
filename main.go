package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func main() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	device, ret := nvml.DeviceGetHandleByIndex(0)
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device at index %d: %v", 0, nvml.ErrorString(ret))
	}

	uuid, ret := device.GetUUID()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get uuid of device at index %d: %v", 0, nvml.ErrorString(ret))
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	for {
		<-c

		pclk, ret := device.GetClockInfo(nvml.CLOCK_VIDEO)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get video clock speed of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		mclk, ret := device.GetClockInfo(nvml.CLOCK_MEM)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get memory clock speed of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		temp, ret := device.GetTemperature(nvml.TEMPERATURE_GPU)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get gpu temperature of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		pwr, ret := device.GetPowerUsage()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get power usage of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		mem_info, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get memory information of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		reasons, ret := device.GetCurrentClocksThrottleReasons()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get throttle reasons of device at index %d: %v", 0, nvml.ErrorString(ret))
		}

		hw_break := reasons & nvml.ClocksThrottleReasonHwPowerBrakeSlowdown
		hw_slowdown := reasons & nvml.ClocksThrottleReasonHwThermalSlowdown
		hw_thermal := reasons & nvml.ClocksThrottleReasonHwThermalSlowdown
		sw_thermal := reasons & nvml.ClocksThrottleReasonSwThermalSlowdown

		fmt.Printf("gpu,hostname=%s,uuid=%s pclk=%d,mclk=%d,temp=%d,pwr=%d,mem_used=%d,mem_free=%d,mem_total=%d,hw_break=%d,hw_slowdown=%d,hw_thermal=%d,sw_thermal=%d %d\n",
			hostname, uuid,
			pclk, mclk, temp, pwr, mem_info.Used, mem_info.Free, mem_info.Total, hw_break, hw_slowdown, hw_thermal, sw_thermal,
			time.Now().UnixNano())
	}
}
