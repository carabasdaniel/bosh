package platform_test

import (
	//"errors"
	//"os"
	//"path/filepath"
	"time"

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
	//boshsys "bosh/system"
	fakesys "bosh/system/fakes"
	//"bosh_platform_impl/platform"
)

var _ = Describe("WindowsPlatform", func() {

	var (
		collector          *fakestats.FakeStatsCollector
		fs                 *fakesys.FakeFileSystem
		cmdRunner          *fakesys.FakeCmdRunner
		dirProvider        boshdirs.DirectoriesProvider
		devicePathResolver *fakedpresolv.FakeDevicePathResolver
		diskManager        *fakedisk.FakeWindowsDiskManager
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
		dirProvider = boshdirs.NewDirectoriesProvider("C:\\fake-dir\\")
		diskManager = fakedisk.NewWindowsFakeDiskManager()
		cdutil = fakecd.NewFakeCdUtil()
		compressor = boshcmd.NewTarballCompressor(cmdRunner, fs)
		copier = boshcmd.NewCpCopier(cmdRunner, fs, logger)
		vitalsService = boshvitals.NewService(collector, dirProvider)
		netManager = &fakenet.FakeNetManager{}
		devicePathResolver = fakedpresolv.NewFakeDevicePathResolver()
	})

	JustBeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)

		platform = NewWindowsPlatform(
			fs,
			cmdRunner,
			collector,
			cdutil,
			dirProvider,
			diskManager,
			netManager,
			10*time.Second,
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

	Describe("CreateUser", func() {
		It("creates user", func() {
			err := platform.CreateUser("bosh_foo-user", "barpwd1234!", "c:\\userbase\\")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AddUserToGroups", func() {
		It("adds user to groups", func() {
			err := platform.AddUserToGroups("bosh_foo-user", []string{"group1", "group2"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SetUserPassword", func() {
		It("set user password", func() {
			err := platform.SetUserPassword("bosh_foo-user", "mypassword123!")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	//TO DO: For Windows this cannot be tested without restarting machine ?
	//Describe("SetupHostname", func() {
	//	It("sets up hostname", func() {
	//		err:=platform.SetupHostname("foobar.local")
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//})

	//TO DO: The fake fs does not go well with this test
	Describe("SetupLogrotate", func() {
		const expectedlogRotateWindowsTemplate = `<?xml version="1.0" encoding="utf-8" ?>
<logRotator poolInterval="900000">  
  <pattern action="rotate" dirPath="C:\testPath\\data\sys\log\*.log" filePattern="*.log" offset="00:01:00" subDirs="true" size="100000"/>
  <pattern action="delete" dirPath="C:\testPath\\data\sys\log\*.log" filePattern="*.gz" offset="00:10:00" deleteUnCompressed="false" subDirs="true"/>
</logRotator>`

		It("sets up logrotate", func() {
			err := platform.SetupLogrotate("fake-group-name", "C:\\testPath\\", "100000")

			logrotateFileContent, err := fs.ReadFileString("C:\\LogRotator\\LogRotator.xml")
			Expect(err).NotTo(HaveOccurred())
			Expect(logrotateFileContent).To(Equal(expectedlogRotateWindowsTemplate))
		})
	})

	Describe("SetTimeWithNtpServers", func() {
		It("sets time with ntp servers", func() {
			err := platform.SetTimeWithNtpServers([]string{"0.north-america.pool.ntp.org", "1.north-america.pool.ntp.org"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SetupDataDir", func() {
		It("creates sys/log and sys/run directories in data directory", func() {
			err := platform.SetupDataDir()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SetupTmpDir", func() {
		It("Sets up temporaty directory with permissions", func() {
			err := platform.SetupTmpDir()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetFileContentsFromCDROM", func() {
		It("delegates to cdutil", func() {
			cdutil.GetFileContentsContents = []byte("fake-contents")
			filename := "fake-env"
			contents, err := platform.GetFileContentsFromCDROM(filename)
			Expect(err).NotTo(HaveOccurred())
			Expect(cdutil.GetFileContentsFilename).To(Equal(filename))
			Expect(contents).To(Equal(cdutil.GetFileContentsContents))
		})
	})

	//Describe("NormalizeDiskPath", func() {
	//	Context("when real device path was resolved without an error", func() {
	//		It("returns real device path and true", func() {
	//			devicePathResolver.RegisterRealDevicePath("fake-device-path", "fake-real-device-path")
	//			realDevicePath, found := platform.NormalizeDiskPath("fake-device-path")
	//			Expect(realDevicePath).To(Equal("fake-real-device-path"))
	//			Expect(found).To(BeTrue())
	//		})
	//	})

	//	Context("when real device path was not resolved without an error", func() {
	//		It("returns real device path and true", func() {
	//			devicePathResolver.GetRealDevicePathErr = errors.New("fake-get-real-device-path-err")

	//			realDevicePath, found := platform.NormalizeDiskPath("fake-device-path")
	//			Expect(realDevicePath).To(Equal(""))
	//			Expect(found).To(BeFalse())
	//		})
	//	})
	//})

	Describe("MountPersistentDisk", func() {
		It("test windows fake mounter", func() {
			err := platform.MountPersistentDisk("3", "C:\\testttt\\")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("UnmountPersistentDisk", func() {
		It("test windows fake mounter", func() {
			_, err := platform.UnmountPersistentDisk("3")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("IsMountPoint", func() {
		It("test windows mounter ismountpoint", func() {
			ok, err := platform.IsMountPoint("C:\\")
			Expect(ok).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("MigratePersistentDisk", func() {
		It("test windows platform migration", func() {
			err := platform.MigratePersistentDisk("C:\\test\\", "C:\\mountP\\")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SetupSSH", func() {
		It("test windows platform ssh setup", func() {
			err := platform.SetupSsh("fake-key", "fake-user")
			authorizedkeycontent, err := fs.ReadFileString("C:\\cygwin\\admin_home\\.ssh\\authorized_keys")
			Expect(len(authorizedkeycontent)).NotTo(BeZero())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	//clean-up user tests
	Describe("DeleteEphemeralUsersMatching", func() {
		It("deletes users with prefix and regex", func() {
			//time.Sleep(5 * time.Second)
			err := platform.DeleteEphemeralUsersMatching("bosh_foo-user")
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
