package disk

import (
	bosherr "bosh/errors"
	boshsys "bosh/system"

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
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Volume")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	mounts := make([]Mount, count)

	for i := 0; i < count; i++ {
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()

		NameasString, _ := oleutil.GetProperty(item, "Name")

		DevIDasString, _ := oleutil.GetProperty(item, "DeviceID")

		mounts[i].PartitionPath = DevIDasString.ToString()

		mounts[i].MountPoint = NameasString.ToString()
	}
	unknown.Release()

	return mounts, nil
}
