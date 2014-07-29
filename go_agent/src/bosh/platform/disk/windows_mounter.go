package disk

import (
	"fmt"
	"strconv"
	"time"

	bosherr "bosh/errors"
	boshsys "bosh/system"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

type windowsMounter struct {
	runner            boshsys.CmdRunner
	mountsSearcher    MountsSearcher
	maxUnmountRetries int
	unmountRetrySleep time.Duration
	dp                DiskPartInterface
}

func NewWindowsMounter(
	runner boshsys.CmdRunner,
	mountsSearcher MountsSearcher,
	unmountRetrySleep time.Duration,
) (mounter windowsMounter) {
	mounter.runner = runner
	mounter.maxUnmountRetries = 600
	mounter.unmountRetrySleep = unmountRetrySleep
	mounter.mountsSearcher = mountsSearcher
	mounter.dp = NewDiskPart()
	return
}

func NewFakeWindowsMounter(
	runner boshsys.CmdRunner,
	mountsSearcher MountsSearcher,
	unmountRetrySleep time.Duration,
	dp DiskPartInterface,
) (mounter windowsMounter) {
	mounter.runner = runner
	mounter.maxUnmountRetries = 600
	mounter.unmountRetrySleep = unmountRetrySleep
	mounter.dp = dp
	mounter.mountsSearcher = mountsSearcher
	return
}
func (m windowsMounter) CreatePrimaryPartition(DiskId int64, label string) error {
	diskIndex := m.GetDiskIndexForDiskId(DiskId)

	scriptfile := fmt.Sprintf("SELECT Disk %d\n ATTRIBUTE DISK CLEAR READONLY\nSELECT Disk %d\n ONLINE DISK NOERR\nSELECT Disk %d\nCREATE PARTITION PRIMARY\nSELECT PARTITION 1\nONLINE VOLUME NOERR\nFORMAT FS=NTFS LABEL=%s QUICK\nEXIT", diskIndex, diskIndex, diskIndex, label)

	_, err := m.dp.ExecuteDiskPartScript(scriptfile)
	if err != nil {
		for i := 1; i < 5; i++ {
			_, err = m.dp.ExecuteDiskPartScript(scriptfile)
			if err != nil {
				bosherr.WrapError(err, "Error executing diskpart scriptfile for primary partition... retrying %d", i)
			}
			return nil
		}
		bosherr.WrapError(err, "Error executing diskpart scriptfile for primary partition")
	}
	return nil
}

//partitionPath should be volumeid
//mountOptions can contain the partition number (can be found by using getPartition) that should be mounted to the path
func (m windowsMounter) Mount(partitionPath, mountPoint string, mountOptions ...string) error {

	volume, err := strconv.ParseInt(partitionPath, 10, 64)

	if err != nil {
		return bosherr.WrapError(err, "Error: bad partitionPath specified for windows mount, it should be an integer value representing the volume number")
	}

	scriptfile := fmt.Sprintf("SELECT VOLUME %d\n REMOVE ALL NOERR\n ASSIGN MOUNT=%s\n EXIT", volume, mountPoint)

	_, err = m.dp.ExecuteDiskPartScript(scriptfile)
	if err != nil {
		for i := 1; i < 5; i++ {
			_, err = m.dp.ExecuteDiskPartScript(scriptfile)
			if err != nil {
				bosherr.WrapError(err, "Error executing diskpart scriptfile for mount... retrying %d", i)
			}
			return nil
		}
		bosherr.WrapError(err, "Error executing diskpart scriptfile for mount")
	}
	return nil
}

func (m windowsMounter) RemountAsReadonly(mountPoint string) error {
	return m.Remount(mountPoint, mountPoint, "-o", "ro")
}

func (m windowsMounter) Remount(fromMountPoint, toMountPoint string, mountOptions ...string) error {
	partitionPath, found, err := m.findDeviceMatchingMountPoint(fromMountPoint)
	if err != nil || !found {
		return bosherr.WrapError(err, "Error finding device for mount point %s", fromMountPoint)
	}

	_, err = m.Unmount(partitionPath)
	if err != nil {
		return bosherr.WrapError(err, "Unmounting %s", fromMountPoint)
	}

	return m.Mount(partitionPath, toMountPoint, mountOptions...)
}

func (m windowsMounter) SwapOn(partitionPath string) (err error) {
	//???

	return nil
}

//partitionOrMountPoint should be the volume index
func (m windowsMounter) Unmount(partitionOrMountPoint string) (bool, error) {
	isMounted, err := m.IsMounted(partitionOrMountPoint)
	if err != nil || !isMounted {
		return false, err
	}

	volumeid, err := strconv.Atoi(partitionOrMountPoint)
	if err != nil {
		return false, err
	}
	scriptfile := fmt.Sprintf("SELECT VOLUME %d\nREMOVE ALL\nEXIT\n", volumeid)
	_, err = m.dp.ExecuteDiskPartScript(scriptfile)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (m windowsMounter) IsMountPoint(path string) (bool, error) {
	mounts, err := m.mountsSearcher.SearchMounts()
	if err != nil {
		return false, bosherr.WrapError(err, "Searching mounts")
	}

	for _, mount := range mounts {
		if mount.MountPoint == path {
			return true, nil
		}
	}

	return false, nil
}

func (m windowsMounter) findDeviceMatchingMountPoint(mountPoint string) (string, bool, error) {
	mounts, err := m.mountsSearcher.SearchMounts()
	if err != nil {
		return "", false, bosherr.WrapError(err, "Searching mounts")
	}

	for _, mount := range mounts {
		if mount.MountPoint == mountPoint {
			return mount.PartitionPath, true, nil
		}
	}

	return "", false, nil
}

func (m windowsMounter) IsMounted(partitionOrMountPoint string) (bool, error) {
	mounts, err := m.mountsSearcher.SearchMounts()
	if err != nil {
		return false, bosherr.WrapError(err, "Searching mounts")
	}

	for _, mount := range mounts {
		if mount.PartitionPath == partitionOrMountPoint || mount.MountPoint == partitionOrMountPoint {
			return true, nil
		}
	}

	return false, nil
}

func (m windowsMounter) shouldMount(partitionPath, mountPoint string) (bool, error) {
	mounts, err := m.mountsSearcher.SearchMounts()
	if err != nil {
		return false, bosherr.WrapError(err, "Searching mounts")
	}

	for _, mount := range mounts {
		switch {
		case mount.PartitionPath == partitionPath && mount.MountPoint == mountPoint:
			return false, nil
		case mount.PartitionPath == partitionPath && mount.MountPoint != mountPoint:
			return false, bosherr.New("Device %s is already mounted to %s, can't mount to %s",
				mount.PartitionPath, mount.MountPoint, mountPoint)
		case mount.MountPoint == mountPoint:
			return false, bosherr.New("Device %s is already mounted to %s, can't mount %s",
				mount.PartitionPath, mount.MountPoint, partitionPath)
		}
	}

	return true, nil
}

func (m windowsMounter) GetDiskIndexForDiskId(diskid int64) int64 {

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_DiskDrive")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()
		asString, _ := oleutil.GetProperty(item, "SCSITargetId")
		scsiTargetId, converr := asString.Value().(int64)
		if !converr {
			return -2
		}
		if scsiTargetId == diskid {
			index, _ := oleutil.GetProperty(item, "Index")
			return index.Value().(int64)
		}
	}
	return -1

}
