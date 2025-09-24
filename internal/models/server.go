package models

import "time"

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
