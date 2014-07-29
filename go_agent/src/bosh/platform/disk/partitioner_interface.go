package disk

type PartitionType string

const (
	PartitionTypeSwap    PartitionType = "swap"
	PartitionTypeLinux   PartitionType = "linux"
	PartitionTypeEmpty   PartitionType = "empty"
	PartitionTypeWindows PartitionType = "windows"
)

type Partition struct {
	SizeInMb uint64
	Type     PartitionType
}

type Partitioner interface {
	Partition(devicePath string, partitions []Partition) (err error)
	GetDeviceSizeInMb(devicePath string) (size uint64, err error)
}
