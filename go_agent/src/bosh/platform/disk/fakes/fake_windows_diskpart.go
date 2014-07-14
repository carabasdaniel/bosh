package fakes

import (
	boshdisk "bosh/platform/disk"
	"io"
	"os"
)

type FakeDiskPart struct {
}

func NewFakeDiskPart() boshdisk.DiskPartInterface {
	return FakeDiskPart{}
}

func (d FakeDiskPart) ExecuteDiskPartScript(script string) (string, error) {
	file, err := os.Create("test_diskpart_script.txt")
	defer os.Remove("test_diskpart_script.txt")
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(file, script)
	if err != nil {
		return "", err
	}
	file.Close()
	return script, nil
}

func (d FakeDiskPart) GetPartitions(diskId int) (partitions []boshdisk.Partition, err error) {
	return nil, nil
}
func (d FakeDiskPart) GetDiskInfo(diskid int) (diskname, status string, size, free uint64) {
	return "nn", "OK", 123, 121
}
func (d FakeDiskPart) GetVolumes() (volumes map[int]string, err error) {
	return nil, nil
}
