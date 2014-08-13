package cdrom

import (
	boshdisk "bosh/platform/disk"
	boshsys "bosh/system"
	"fmt"
)

type WindowsCdrom struct {
	diskManager boshdisk.Manager
	runner      boshsys.CmdRunner
}

func NewWindowsCdrom(diskManager boshdisk.Manager, runner boshsys.CmdRunner) (cdrom WindowsCdrom) {
	cdrom = WindowsCdrom{
		diskManager: diskManager,
		runner:      runner,
	}
	return
}

func getCdromVolumeId() (id string, err error) {

	diskpart := boshdisk.NewDiskPart()
	cdromDevices, err := diskpart.GetVolumes("DVD-ROM")
	fmt.Println("CD rom devices:", cdromDevices)
	if err != nil {
		return "-1", err
	}

	for key, _ := range cdromDevices {
		id = string(key)
		break
	}

	return id, nil
}

func (cdrom WindowsCdrom) WaitForMedia() (err error) {
	return nil
}
func (cdrom WindowsCdrom) Mount(mountPath string) (err error) {
	cdromId, err := getCdromVolumeId()
	if err != nil {
		return err
	}

	err = cdrom.diskManager.GetMounter().Mount(cdromId, mountPath)
	if err != nil {
		return err
	}

	return nil
}
func (cdrom WindowsCdrom) Unmount() (err error) {
	cdromId, err := getCdromVolumeId()
	if err != nil {
		return err
	}
	_, err = cdrom.diskManager.GetMounter().Unmount(cdromId)
	if err != nil {
		return err
	}

	return nil
}
func (cdrom WindowsCdrom) Eject() (err error) {
	return nil
}
