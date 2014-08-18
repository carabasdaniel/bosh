package cdrom

import (
	boshlog "bosh/logger"
	boshdisk "bosh/platform/disk"
	boshsys "bosh/system"
	"strconv"
)

type WindowsCdrom struct {
	diskManager boshdisk.Manager
	runner      boshsys.CmdRunner
	logger      boshlog.Logger
}

func NewWindowsCdrom(diskManager boshdisk.Manager, runner boshsys.CmdRunner, logger boshlog.Logger) (cdrom WindowsCdrom) {
	cdrom = WindowsCdrom{
		diskManager: diskManager,
		runner:      runner,
		logger:      logger,
	}
	return
}

func getCdromVolumeId() (id string, err error) {

	diskpart := boshdisk.NewDiskPart()
	cdromDevices, err := diskpart.GetVolumes("DVD-ROM")
	if err != nil {
		return "-1", err
	}

	for key, _ := range cdromDevices {
		id = strconv.Itoa(key)
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

	cdrom.logger.Debug("CdRom", "Unmounting cdrom")
	cdromId, err := getCdromVolumeId()
	if err != nil {
		return err
	}
	cdrom.logger.Debug("CdRom", "Found cdrom volume ID %s", cdromId)
	_, err = cdrom.diskManager.GetMounter().Unmount(cdromId)
	if err != nil {
		return err
	}

	cdrom.logger.Debug("CdRom", "Unmount Done")
	return nil
}
func (cdrom WindowsCdrom) Eject() (err error) {
	return nil
}
