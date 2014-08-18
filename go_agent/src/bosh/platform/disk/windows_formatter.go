package disk

import (
	"fmt"
	"strconv"

	bosherr "bosh/errors"
	boshsys "bosh/system"
)

type windowsFormatter struct {
	runner boshsys.CmdRunner
	fs     boshsys.FileSystem
	DPart  DiskPartInterface
}

func NewWindowsFormatter(runner boshsys.CmdRunner, fs boshsys.FileSystem) (formatter windowsFormatter) {
	formatter.runner = runner
	formatter.fs = fs
	formatter.DPart = NewDiskPart()
	return formatter
}

func NewFakeWindowsFormatter(runner boshsys.CmdRunner, fs boshsys.FileSystem, dp DiskPartInterface) (formatter windowsFormatter) {
	formatter.runner = runner
	formatter.fs = fs
	formatter.DPart = dp
	return formatter
}

func (f windowsFormatter) ChangeDiskPart(dp DiskPartInterface) {
	f.DPart = dp
	return
}

func (f windowsFormatter) Format(partitionPath string, fsType FileSystemType) (err error) {

	if fsType == FileSystemNtfs {

		volumeid, converr := strconv.Atoi(partitionPath)

		if converr != nil {
			_, _, _, err = f.runner.RunCommand("format", partitionPath, "/q", "/FS:NTFS", "/Y")

			if err != nil {
				err = bosherr.WrapError(err, "Shelling out format")

			}
		} else {
			script := fmt.Sprintf("SELECT VOLUME %d\n FORMAT FS=NTFS QUICK\n EXIT", volumeid)
			_, err := f.DPart.ExecuteDiskPartScript(script)
			if err != nil {
				err = bosherr.WrapError(err, "Shelling out diskpart formatting")
			}
		}
	}
	return
}
