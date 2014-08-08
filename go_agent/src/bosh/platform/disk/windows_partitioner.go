package disk

import (
	"fmt"
	"strconv"

	bosherr "bosh/errors"
	boshlog "bosh/logger"
	boshsys "bosh/system"
)

type windowsPartitioner struct {
	logger    boshlog.Logger
	cmdRunner boshsys.CmdRunner
	logTag    string
	Dpart     DiskPartInterface
}

func NewWindowsPartitioner(logger boshlog.Logger, cmdRunner boshsys.CmdRunner) (partitioner windowsPartitioner) {
	partitioner.logger = logger
	partitioner.cmdRunner = cmdRunner
	partitioner.logTag = "WindowsPartitioner"
	partitioner.Dpart = NewDiskPart()
	return
}

func NewFakeWindowsPartitioner(logger boshlog.Logger, cmdRunner boshsys.CmdRunner, dp DiskPartInterface) (partitioner windowsPartitioner) {
	partitioner.logger = logger
	partitioner.cmdRunner = cmdRunner
	partitioner.logTag = "WindowsPartitioner"
	partitioner.Dpart = dp
	return
}

func (p windowsPartitioner) Partition(devicePath string, partitions []Partition) (err error) {
	diskId, err := strconv.Atoi(devicePath)
	if err != nil {
		return err
	}

	if p.diskMatchesPartitions(diskId, partitions) {
		return
	}

	script := fmt.Sprintf("SELECT DISK %d\n CLEAN\n", diskId)

	for _, a := range partitions {
		if a.SizeInMb == 0 {
			script = script + fmt.Sprintf("CREATE PARTITION PRIMARY\n")
		} else {
			freeSpace, _ := p.GetDeviceSizeInMb(devicePath)
			if a.SizeInMb > freeSpace {
				script = script + fmt.Sprintf("CREATE PARTITION PRIMARY\n")
			} else {
				script = script + fmt.Sprintf("CREATE PARTITION PRIMARY SIZE=%d\n", a.SizeInMb)
			}
		}
	}
	script = script + "EXIT\n"

	_, err = p.Dpart.ExecuteDiskPartScript(script)

	if err != nil {
		return err
	}

	return
}

func (p windowsPartitioner) GetDeviceSizeInMb(devicePath string) (size uint64, err error) {

	diskId, err := strconv.Atoi(devicePath)

	if err != nil {
		return 0, bosherr.WrapError(err, "Error: devicePath should be an integer value representing the physical disk drive index")
	}

	_, _, _, free := p.Dpart.GetDiskInfo(diskId)

	return free, nil
}

func (p windowsPartitioner) diskMatchesPartitions(diskId int, partitionsToMatch []Partition) (result bool) {
	existingPartitions, err := p.Dpart.GetPartitions(diskId)
	if err != nil {
		err = bosherr.WrapError(err, "Getting partitions for disk %d", diskId)
		return
	}

	if len(existingPartitions) < len(partitionsToMatch) {
		return
	}

	remainingDiskSpace, err := p.GetDeviceSizeInMb(strconv.Itoa(diskId))
	if err != nil {
		err = bosherr.WrapError(err, "Getting device size for disk %d", diskId)
		return
	}

	for index, partitionToMatch := range partitionsToMatch {
		if index == len(partitionsToMatch)-1 {
			partitionToMatch.SizeInMb = remainingDiskSpace
		}

		existingPartition := existingPartitions[index]
		switch {
		case existingPartition.Type != partitionToMatch.Type:
			return
		case notWithinDelta(existingPartition.SizeInMb, partitionToMatch.SizeInMb, 20):
			return
		}

		remainingDiskSpace = remainingDiskSpace - partitionToMatch.SizeInMb
	}

	return true
}
