package models

import "time"

type ServerInfo struct {
	ModelName     string
	RamToTal      uint64
	DiskTotal     uint64
	OS            string
	Platform      string
	KernelVersion string
	Uptime        uint64
}

func NewServerInfo(modelName, os, platform, kernelVersion string, ramTotal, diskTotal, uptime uint64) *ServerInfo {
	return &ServerInfo{
		ModelName:     modelName,
		RamToTal:      ramTotal,
		DiskTotal:     diskTotal,
		OS:            os,
		Platform:      platform,
		KernelVersion: kernelVersion,
		Uptime:        uptime,
	}
}

type ServerMetrics struct {
	CPU  int       `json:"cpu" example:"67"`
	RAM  int       `json:"ram" example:"3211"`
	Disk int       `json:"disk" example:"32423"`
	Date time.Time `json:"date" example:""`
}

func NewServerMetrics(cpu, ram, disk int, date time.Time) *ServerMetrics {
	return &ServerMetrics{
		CPU:  cpu,
		RAM:  ram,
		Disk: disk,
		Date: date,
	}
}
