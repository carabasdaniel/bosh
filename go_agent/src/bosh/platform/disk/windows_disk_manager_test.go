package disk_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "bosh/logger"
	. "bosh/platform/disk"
	fakesys "bosh/system/fakes"
)

func init() {
	Describe("windowsDiskManager", func() {
		var (
			runner *fakesys.FakeCmdRunner
			fs     *fakesys.FakeFileSystem
			logger boshlog.Logger
		)

		It("returns windows disk manager configured", func() {
			runner = fakesys.NewFakeCmdRunner()
			fs = fakesys.NewFakeFileSystem()
			logger = boshlog.NewLogger(boshlog.LevelNone)

			expectedMountsSearcher := NewWindowsMountsSearcher(fs)
			expectedMounter := NewWindowsMounter(runner, expectedMountsSearcher, 1*time.Second)

			diskManager := NewWindowsDiskManager(logger, runner, fs, false)
			Expect(diskManager.GetMounter()).To(Equal(expectedMounter))
		})

	})
}
