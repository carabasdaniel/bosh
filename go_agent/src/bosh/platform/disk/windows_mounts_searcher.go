package disk

import (
	bosherr "bosh/errors"
	boshsys "bosh/system"
	"strconv"
	"strings"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

type windowsMountsSearcher struct {
	fs boshsys.FileSystem
}

func NewWindowsMountsSearcher(fs boshsys.FileSystem) MountsSearcher {
	return windowsMountsSearcher{fs}
}

func (m windowsMountsSearcher) SearchMounts() ([]Mount, error) {
	var errCallMethod error
	var property *ole.VARIANT
	ole.CoInitialize(0)
	defer ole.CoUninitialize()
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")

	if err != nil {
		return nil, bosherr.WrapError(err, "Error IUnknown interface init")
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, bosherr.WrapError(err, "WMI query interface mounts")
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, errConnectServer := oleutil.CallMethod(wmi, "ConnectServer")
	if errConnectServer != nil {
		return nil, bosherr.WrapError(errConnectServer, "WMI Connect Server error")
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, errExecQuery := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Volume")
	if errExecQuery != nil {
		return nil, bosherr.WrapError(errExecQuery, "WMI query error")
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, errGetProperty := oleutil.GetProperty(result, "Count")
	if errGetProperty != nil {
		return nil, bosherr.WrapError(errGetProperty, "WMI Get Property error")
	}
	count := int(countVar.Val)

	mounts := make([]Mount, count)

	for i := 0; i < count; i++ {
		property, errCallMethod = oleutil.CallMethod(result, "ItemIndex", i)
		if errCallMethod != nil {
			return nil, bosherr.WrapError(errCallMethod, "ItemIndex Call method")
		}
		item := property.ToIDispatch()

		property, errCallMethod = oleutil.GetProperty(item, "Name")
		if errCallMethod != nil {
			return nil, bosherr.WrapError(errCallMethod, "Get Name")
		}
		mounts[i].MountPoint = property.ToString()

		property, errCallMethod = oleutil.GetProperty(item, "DeviceID")
		if errCallMethod != nil {
			return nil, bosherr.WrapError(errCallMethod, "Get Device Id ")
		}
		mounts[i].PartitionPath = property.ToString()
		item.Release()
	}

	vols, errDisk := DiskPart{}.GetVolumes(" ")
	if errDisk != nil {
		return nil, bosherr.WrapError(errDisk, "Disk Part")
	}
	for k := range mounts {
		for n, i := range vols {
			if len(mounts[k].MountPoint) == 3 {
				letter := strings.Replace(mounts[k].MountPoint, ":\\", "", -1)
				data := strings.Split(i, "-")
				if data[0] == letter {
					mounts[k].PartitionPath = strconv.Itoa(n)
				}
			} else {
				if strings.Contains(i, mounts[k].MountPoint) {
					mounts[k].PartitionPath = strconv.Itoa(n)
				}
			}
		}
	}
	return mounts, nil
}
