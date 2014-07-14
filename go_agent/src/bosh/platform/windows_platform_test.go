package platform_test

import (
	//"errors"
	//"os"
	//"path/filepath"
	//"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakedpresolv "bosh/infrastructure/devicepathresolver/fakes"
	boshlog "bosh/logger"
	. "bosh/platform"
	fakecd "bosh/platform/cdutil/fakes"
	boshcmd "bosh/platform/commands"
	//boshdisk "bosh/platform/disk"
	fakedisk "bosh/platform/disk/fakes"
	fakenet "bosh/platform/net/fakes"
	fakestats "bosh/platform/stats/fakes"
	boshvitals "bosh/platform/vitals"
	//boshsettings "bosh/settings"
	boshdirs "bosh/settings/directories"
	fakesys "bosh/system/fakes"
	//"bosh_platform_impl/platform"
)

var _ = Describe("WindowsPlatform", func() {

	var (
		collector          *fakestats.FakeStatsCollector
		fs                 *fakesys.FakeFileSystem
		cmdRunner          *fakesys.FakeCmdRunner
		diskManager        *fakedisk.FakeDiskManager
		dirProvider        boshdirs.DirectoriesProvider
		devicePathResolver *fakedpresolv.FakeDevicePathResolver
		platform           Platform
		cdutil             *fakecd.FakeCdUtil
		compressor         boshcmd.Compressor
		copier             boshcmd.Copier
		vitalsService      boshvitals.Service
		netManager         *fakenet.FakeNetManager
	)

	BeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)

		fs = fakesys.NewFakeFileSystem()
		cmdRunner = fakesys.NewFakeCmdRunner()
		collector = &fakestats.FakeStatsCollector{}
		diskManager = fakedisk.NewFakeDiskManager()
		dirProvider = boshdirs.NewDirectoriesProvider("\\fake-dir")
		cdutil = fakecd.NewFakeCdUtil()
		compressor = boshcmd.NewTarballCompressor(cmdRunner, fs)
		copier = boshcmd.NewCpCopier(cmdRunner, fs, logger)
		vitalsService = boshvitals.NewService(collector, dirProvider)
		netManager = &fakenet.FakeNetManager{}
		devicePathResolver = fakedpresolv.NewFakeDevicePathResolver()

		fs.SetGlob("/sys/bus/scsi/devices/*:0:0:0/block/*", []string{
			"/sys/bus/scsi/devices/0:0:0:0/block/sr0",
			"/sys/bus/scsi/devices/6:0:0:0/block/sdd",
			"/sys/bus/scsi/devices/fake-host-id:0:0:0/block/sda",
		})

		fs.SetGlob("/sys/bus/scsi/devices/fake-host-id:0:fake-disk-id:0/block/*", []string{
			"/sys/bus/scsi/devices/fake-host-id:0:fake-disk-id:0/block/sdf",
		})
	})

	JustBeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)

		platform = NewWindowsPlatform(
			collector,
			fs,
			cmdRunner,
			cdutil,
			dirProvider,
			diskManager,
			netManager,
			logger,
		)

		platform.SetDevicePathResolver(devicePathResolver)
	})

	Describe("SetupRuntimeConfiguration", func() {
		It("setups runtime configuration", func() {
			err := platform.SetupRuntimeConfiguration()
			Expect(err).NotTo(HaveOccurred())

			Expect(len(cmdRunner.RunCommands)).To(Equal(1))
			Expect(cmdRunner.RunCommands[0]).To(Equal([]string{"bosh-agent-rc"}))
		})
	})

	Describe("FAIL", func() {
		It("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE", func() {
			Expect(1).To(Equal(0))
		})
	})

})
