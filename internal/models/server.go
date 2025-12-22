package models

import "time"

type ServerInfo struct {
	ModelName     string `json:"model_name"`
	RAMToTal      uint64 `json:"ram_total"`
	DiskTotal     uint64 `json:"disk_total"`
	OS            string `json:"os"`
	Platform      string `json:"platform"`
	KernelVersion string `json:"kernel_version"`
	Uptime        uint64 `json:"uptime"`
}

func NewServerInfo(modelName, os, platform, kernelVersion string, ramTotal, diskTotal, uptime uint64) *ServerInfo {
	return &ServerInfo{
		ModelName:     modelName,
		RAMToTal:      ramTotal,
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
