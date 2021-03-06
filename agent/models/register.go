package models

type RegisterMsg struct {
	UUID          string         `json:"uuid"`
	Hostname      string         `json:"hostname"`
	OutBoundIP    NullString `json:"outBoundIp"`
	ClusterIP     NullString `json:"clusterIp"`
	IPs           string         `json:"ips"`
	OS            string         `json:"os"`
	Arch          string         `json:"arch"`
	CpuCount      int32          `json:"cpuCount"`
	CpuUsePercent float64        `json:"cpuUsePercent"`
	RamTotal      int64          `json:"ramTotal"` // MB
	RamUsed       int64          `json:"ramUsed"`  // MB
	RamUsePercent float64        `json:"ramPercent"`
	Disks         string         `json:"disks"`
	AvgLoad       string         `json:"avgLoad"`
	BootTime      NullInt64  `json:"bootTime"`
}

type Disk struct {
	Total      int64   `json:"total"`
	Used       int64   `json:"used"`
	UsePercent float64 `json:"usePercent"`
}
