package fakes

import (
	boshlog "bosh/logger"
	boshdisk "bosh/platform/disk"
	fakesys "bosh/system/fakes"
	"time"
)

type FakeWindowsDiskManager struct {
	FakePartitioner boshdisk.Partitioner
	FakeFormatter   boshdisk.Formatter
	FakeMounter     boshdisk.Mounter
}

func NewWindowsFakeDiskManager() (manager *FakeWindowsDiskManager) {
	manager = &FakeWindowsDiskManager{}

	fakeRunner := fakesys.NewFakeCmdRunner()
	fakeFs := fakesys.NewFakeFileSystem()
	logger := boshlog.NewLogger(boshlog.LevelNone)

	manager.FakePartitioner = boshdisk.NewFakeWindowsPartitioner(logger, fakeRunner, NewFakeDiskPart())
	manager.FakeFormatter = boshdisk.NewFakeWindowsFormatter(fakeRunner, fakeFs, NewFakeDiskPart())

	searcher := boshdisk.NewWindowsMountsSearcher(fakeFs)
	manager.FakeMounter = boshdisk.NewFakeWindowsMounter(fakeRunner, searcher, 1*time.Second, NewFakeDiskPart())
	return
}

func (m FakeWindowsDiskManager) GetPartitioner() boshdisk.Partitioner {
	return m.FakePartitioner
}

func (m FakeWindowsDiskManager) GetFormatter() boshdisk.Formatter {
	return m.FakeFormatter
}

func (m FakeWindowsDiskManager) GetMounter() boshdisk.Mounter {
	return m.FakeMounter
}
