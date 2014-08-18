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
	var exists bool
	diskId, err := strconv.Atoi(devicePath)
	if err != nil {
		return err
	}

	exists, err = p.diskMatchesPartitions(diskId, partitions[0])
	if err != nil {
		return err
	}
	if exists == true {
		return nil
	}

	errDiskPart := p.executeDiskPart(fmt.Sprintf("ONLINE DISK\n ATTRIBUTES DISK CLEAR READONLY\n"), diskId)
	if errDiskPart != nil {
		p.logger.Debug("Partitioner", fmt.Sprintf("Disk already online %s", errDiskPart))
	}

	for _, a := range partitions {
		if a.SizeInMb == 0 {
			return
		} else {
			freeSpace, _ := p.GetDeviceSizeInMb(devicePath)
			if a.SizeInMb > freeSpace {
				err = p.executeDiskPart(fmt.Sprintf("CREATE PARTITION PRIMARY\n"), diskId)

			} else {
				err = p.executeDiskPart(fmt.Sprintf("CREATE PARTITION PRIMARY SIZE=%d\n", a.SizeInMb), diskId)

			}
		}
		if err != nil {
			return bosherr.WrapError(err, fmt.Sprintf("Error creation partition with size %d for %d", a.SizeInMb, diskId))
		}
	}

	return
}

func (p windowsPartitioner) executeDiskPart(command string, diskId int) (err error) {
	script := fmt.Sprintf("SELECT DISK %d\n %s\n EXIT\n", diskId, command)
	_, err = p.Dpart.ExecuteDiskPartScript(script)
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

func (p windowsPartitioner) diskMatchesPartitions(diskId int, partitionToMatch Partition) (exists bool, err error) {
	existingPartitions, err := p.Dpart.GetPartitions(diskId)
	if err != nil {
		err = bosherr.WrapError(err, "Getting partitions for disk %s", diskId)
		return
	}

	for _, partition := range existingPartitions {
		if partition.Type == partitionToMatch.Type &&
			partition.SizeInMb == partitionToMatch.SizeInMb {
			exists = true
			return
		}
	}
	exists = false
	return

}
