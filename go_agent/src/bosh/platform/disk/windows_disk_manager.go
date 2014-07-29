package disk

import (
	"time"

	boshlog "bosh/logger"
	boshsys "bosh/system"
)

type windowsDiskManager struct {
	partitioner Partitioner
	formatter   Formatter
	mounter     Mounter
}

func NewWindowsDiskManager(
	logger boshlog.Logger,
	runner boshsys.CmdRunner,
	fs boshsys.FileSystem,
	bindMount bool,
) (manager Manager) {
	var mounter Mounter
	var mountsSearcher MountsSearcher

	mountsSearcher = NewWindowsMountsSearcher(fs)

	mounter = NewWindowsMounter(runner, mountsSearcher, 1*time.Second)

	return windowsDiskManager{
		partitioner: NewWindowsPartitioner(logger, runner),
		formatter:   NewWindowsFormatter(runner, fs),
		mounter:     mounter,
	}
}

func (m windowsDiskManager) GetPartitioner() Partitioner { return m.partitioner }
func (m windowsDiskManager) GetFormatter() Formatter     { return m.formatter }
func (m windowsDiskManager) GetMounter() Mounter         { return m.mounter }
