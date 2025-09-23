package models

type ServerMetrics struct {
	CPU  int `json:"cpu" example:"67"`
	RAM  int `json:"ram" example:"3211"`
	Disk int `json:"disk" example:"32423"`
	//Network int `json:"network" example:"32423"`
}

func NewServerMetrics(cpu int, ram int, disk int) *ServerMetrics {
	return &ServerMetrics{
		CPU:  cpu,
		RAM:  ram,
		Disk: ram,
	}
}
