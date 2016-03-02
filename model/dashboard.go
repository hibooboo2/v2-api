package model

import "github.com/rancher/go-rancher/client"

type DashBoard struct {
	Hosts      *Hosts `json:"hosts"`
	Stacks     *Stacks
	Containers *Containers
	Services   *Services
	Processes  *Processes
	AuditLogs  []client.AuditLog
}

type BucketOfIDs struct {
	RangeStart float64
	RangeEnd   float64
	IDs        []string
}

type HostsBucket struct {
	Buckets []*BucketOfIDs
}

func (b *HostsBucket) AddValue(val float64, ID string) bool {
	for _, bucket := range b.Buckets {
		if bucket.AddValue(val, ID) {
			return true
		}
	}
	return false
}

func (b *BucketOfIDs) AddValue(val float64, ID string) bool {
	if b.RangeStart <= val && val < b.RangeEnd {
		b.IDs = append(b.IDs, ID)
		return true
	}
	return false
}

type Health struct {
	Type string
	ID   string
}

type Hosts struct {
	Health  []Health
	CPU     *HostsBucket
	Memory  *HostsBucket
	Disk    *HostsBucket
	Network *HostsBucket
}

type Stacks struct {
	Health []Health
}

type Containers struct {
	Health []Health
}

type Services struct {
	Health []Health
}

type Processes struct {
	CurrentCount          int64
	LongRunning           int64
	LongRunningTime       int64
	NumRecalled           int64
	TimeToRecallProcesses int64
}

type Host struct {
	Data  string `db:"data"`
	State string `db:"state"`
	ID    int64  `db:"id"`
}

type HostData struct {
	ID          int64 `json:"id"`
	FormattedID string
	Fields      HostDataFields `json:"fields"`
}

type HostDataFields struct {
	Info HostInfo `json:"info"`
}

type HostInfo struct {
	CPUInfo    CPUInfo    `json:"cpuInfo"`
	DiskInfo   DiskInfo   `json:"diskInfo"`
	MemoryInfo MemoryInfo `json:"memoryInfo"`
}

type CPUInfo struct {
	CPUCoresPercentages []float64 `json:"cpuCoresPercentages"`
}

type DiskInfo struct {
	MountPoints map[string]MountPoint
}

type MemoryInfo struct {
	SwapFree       float64
	MemFree        float64
	MemTotal       float64
	Buffers        float64
	Cached         float64
	UsedPercentage float64
}

type MountPoint struct {
	PercentUsed float64
}
