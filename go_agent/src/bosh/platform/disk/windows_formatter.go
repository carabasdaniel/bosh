package disk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"

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

func (f windowsFormatter) ChangeDiskPart(dp DiskPartInterface) {
	f.DPart = dp
	return
}

func (f windowsFormatter) Format(partitionPath string, fsType FileSystemType) (err error) {
	if f.partitionHasGivenType(partitionPath, fsType) {
		return
	}

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

func (f windowsFormatter) partitionHasGivenType(partitionPath string, fsType FileSystemType) bool {
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

	volumeid, err := strconv.Atoi(partitionPath)

	if err != nil {

		// result is a SWBemObjectSet
		resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Volume ")
		result := resultRaw.ToIDispatch()
		defer result.Release()

		countVar, _ := oleutil.GetProperty(result, "Count")
		count := int(countVar.Val)

		for i := 0; i < count; i++ {
			// item is a SWbemObject, but really a Win32_Process
			itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
			item := itemRaw.ToIDispatch()
			defer item.Release()
			asString, _ := oleutil.GetProperty(item, "Name")
			TypeAsString, _ := oleutil.GetProperty(item, "FileSystem")

			if strings.ToLower(asString.ToString()) == strings.ToLower(partitionPath) && string(fsType) == strings.ToLower(TypeAsString.ToString()) {
				return true
			}
		}
	} else {
		resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Volume WHERE DriveType='3'")
		result := resultRaw.ToIDispatch()
		defer result.Release()

		countVar, _ := oleutil.GetProperty(result, "Count")
		count := int(countVar.Val)

		removableRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Volume WHERE DriveType<>'3'")
		removable := removableRaw.ToIDispatch()
		defer removable.Release()

		initVar, _ := oleutil.GetProperty(removable, "Count")
		init := int(initVar.Val)

		for i := init; i < (count + init); i++ {
			// item is a SWbemObject, but really a Win32_Process
			//if you want to keep the volume index number the same as diskpart remember that diskpart takes removables first
			//count all removables and from there you have the fixed disks indexes in order
			itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i-init)
			item := itemRaw.ToIDispatch()
			defer item.Release()
			//asString, _ := oleutil.GetProperty(item, "Name")
			TypeAsString, _ := oleutil.GetProperty(item, "FileSystem")

			if i == volumeid && string(fsType) == strings.ToLower(TypeAsString.ToString()) {
				return true
			}
		}
	}

	return false
}
