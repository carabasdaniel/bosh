package disk

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type DiskPartInterface interface {
	ExecuteDiskPartScript(script string) (string, error)
	GetPartitions(diskId int) (partitions []Partition, err error)
	GetDiskInfo(diskid int) (diskname, status string, size, free uint64)
	GetVolumes() (volumes map[int]string, err error)
}

type DiskPart struct {
}

func NewDiskPart() DiskPartInterface {
	return DiskPart{}
}

func (d DiskPart) ExecuteDiskPartScript(script string) (string, error) {
	file, err := os.Create("diskpart_script.txt")
	defer os.Remove("diskpart_script.txt")
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(file, script)
	if err != nil {
		return "", err
	}
	file.Close()

	output, err := exec.Command("diskpart.exe", "/s", "diskpart_script.txt").Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (d DiskPart) GetPartitions(diskId int) (partitions []Partition, err error) {
	script := fmt.Sprintf("Select disk %d\n list partition\nEXIT", diskId)
	output, err := d.ExecuteDiskPartScript(script)
	if err != nil {
		return nil, err
	}
	content := strings.Split(output, "\n")

	found := false
	partinfos := make(map[string][]string)
	for _, line := range content {
		if strings.Contains(line, "GB") {
			info := strings.Split(line, "  ")
			for _, data := range info {
				if len(strings.Trim(data, " ")) > 1 && !strings.EqualFold(data, info[0]) {
					partinfos[strings.TrimSpace(info[0])] = append(partinfos[strings.TrimSpace(info[0])], strings.Trim(data, " "))
				}
			}
			found = true
		}
	}

	for key := range partinfos {
		var part Partition

		size_asString := strings.TrimSpace(strings.Replace(partinfos[key][2], "GB", "", -1))
		size, err := strconv.ParseUint(size_asString, 10, 64)
		if err != nil {
			return nil, err
		}
		size = size * 1024

		part.SizeInMb = size
		part.Type = PartitionTypeWindows

		partitions = append(partitions, part)
	}

	if !found {
		return nil, errors.New(fmt.Sprintf("No partitions found on disk %d", diskId))
	}

	return partitions, nil
}

func (d DiskPart) GetDiskInfo(diskid int) (diskname, status string, size, free uint64) {

	key := "Disk " + strconv.Itoa(diskid)
	output, _ := d.ExecuteDiskPartScript("list disk")
	content := strings.Split(output, "\n")
	diskinfo := make(map[string][]string)

	for _, a := range content {
		if strings.Contains(a, "GB") {
			info := strings.Split(a, "    ")
			for _, b := range info {
				if len(strings.Trim(b, " ")) > 1 {
					diskinfo[strings.TrimSpace(info[0])] = append(diskinfo[strings.TrimSpace(info[0])], strings.Trim(b, " "))
				}
			}

		}
	}

	diskname = diskinfo[key][0]
	status = diskinfo[key][1]
	size_asString := strings.TrimSpace(strings.Replace(strings.Replace(diskinfo[key][2], "GB", "", -1), "B", "", -1))
	free_asString := strings.TrimSpace(strings.Replace(strings.Replace(diskinfo[key][3], "GB", "", -1), "B", "", -1))
	size, err := strconv.ParseUint(size_asString, 10, 64)
	if err != nil {
		panic(err)
	}
	size = size * 1024

	free, err = strconv.ParseUint(free_asString, 10, 64)
	if err != nil {
		panic(err)
	}
	free = free * 1024

	return diskname, status, size, free
}

func (d DiskPart) GetVolumes() (volumes map[int]string, err error) {
	script := fmt.Sprintf("list volume\nEXIT")
	output, err := d.ExecuteDiskPartScript(script)
	if err != nil {
		return nil, err
	}
	content := strings.Split(output, "\n")

	found := false
	volinfos := make(map[string][]string)

	lastkey := ""
	for _, line := range content {
		if strings.Contains(line, "Partition") {
			info := strings.Split(line, "   ")
			lastkey = strings.TrimSpace(info[0])
			for _, data := range info {
				if len(strings.Trim(data, " ")) >= 1 && !strings.EqualFold(data, info[0]) {
					volinfos[lastkey] = append(volinfos[lastkey], strings.Trim(strings.TrimSpace(data), "\n"))
				}
			}
			found = true
		} else {
			if strings.Contains(line, ":\\") {
				volinfos[lastkey] = append(volinfos[lastkey], strings.Trim(strings.TrimSpace(line), "\n"))
			}
		}
	}

	volumes = make(map[int]string)

	if found {
		for key := range volinfos {
			index, _ := strconv.Atoi(strings.Replace(key, "Volume ", "", -1))
			volumes[index] = strings.Join(volinfos[key], "-")
		}
	}

	return volumes, nil
}
