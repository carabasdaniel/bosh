package platform

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	bosherr "bosh/errors"
	boshdpresolv "bosh/infrastructure/devicepathresolver"
	boshlog "bosh/logger"
	boshcmd "bosh/platform/commands"
	boshstats "bosh/platform/stats"
	boshvitals "bosh/platform/vitals"
	boshsettings "bosh/settings"
	boshdir "bosh/settings/directories"
	boshdirs "bosh/settings/directories"
	boshsys "bosh/system"
)

type windowsPlatform struct {
	collector          boshstats.StatsCollector
	fs                 boshsys.FileSystem
	cmdRunner          boshsys.CmdRunner
	compressor         boshcmd.Compressor
	copier             boshcmd.Copier
	dirProvider        boshdirs.DirectoriesProvider
	vitalsService      boshvitals.Service
	devicePathResolver boshdpresolv.DevicePathResolver
	logger             boshlog.Logger
}

func NewWindowsPlatform(
	collector boshstats.StatsCollector,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	dirProvider boshdirs.DirectoriesProvider,
	logger boshlog.Logger,
) *windowsPlatform {
	return &windowsPlatform{
		fs:            fs,
		cmdRunner:     cmdRunner,
		collector:     collector,
		compressor:    boshcmd.NewTarballCompressor(cmdRunner, fs),
		copier:        boshcmd.NewCpCopier(cmdRunner, fs, logger),
		dirProvider:   boshdir.NewDirectoriesProvider("C:/"),
		vitalsService: boshvitals.NewService(collector, dirProvider),
	}
}

func (p windowsPlatform) GetFs() (fs boshsys.FileSystem) {
	fmt.Printf("Get FS\n")
	//fmt.Printf(p.dirProvider.BoshDir())
	GenerateFakeTask()
	return p.fs
}

func (p windowsPlatform) GetRunner() (runner boshsys.CmdRunner) {
	fmt.Printf("Get Runner\n")
	return p.cmdRunner
}

func (p windowsPlatform) GetCompressor() (compressor boshcmd.Compressor) {
	fmt.Printf("Get compressor\n")
	return p.compressor
}

func (p windowsPlatform) GetCopier() (copier boshcmd.Copier) {
	fmt.Printf("Get copier\n")
	return p.copier
}

func (p windowsPlatform) GetDirProvider() (dirProvider boshdir.DirectoriesProvider) {
	fmt.Printf("Get dir provider\n")
	return p.dirProvider
}

func (p windowsPlatform) GetVitalsService() (service boshvitals.Service) {
	fmt.Printf("Get vital service\n")
	return p.vitalsService
}

func (p windowsPlatform) GetDevicePathResolver() (devicePathResolver boshdpresolv.DevicePathResolver) {
	fmt.Printf("Get device path resolver\n")
	return p.devicePathResolver
}

func (p *windowsPlatform) SetDevicePathResolver(devicePathResolver boshdpresolv.DevicePathResolver) (err error) {
	fmt.Printf("Set device path resolver\n")
	p.devicePathResolver = devicePathResolver
	return
}

func (p windowsPlatform) SetupRuntimeConfiguration() (err error) {
	fmt.Printf("Setup runtime configuration\n")
	return
}

func (p windowsPlatform) CreateUser(username, password, basePath string) (err error) {
	fmt.Printf("Create user\n")
	return
}

func (p windowsPlatform) AddUserToGroups(username string, groups []string) (err error) {
	fmt.Printf("Add user to groups\n")
	return
}

func (p windowsPlatform) DeleteEphemeralUsersMatching(regex string) (err error) {
	fmt.Printf("Delete ephemeral users matching\n")
	return
}

func (p windowsPlatform) SetupSsh(publicKey, username string) (err error) {
	fmt.Printf("Setup ssh\n")
	return
}

func (p windowsPlatform) SetUserPassword(user, encryptedPwd string) (err error) {
	fmt.Printf("Set user password\n")
	return
}

func (p windowsPlatform) SetupHostname(hostname string) (err error) {
	fmt.Printf("Setup hostname\n")
	return
}

func (p windowsPlatform) SetupDhcp(networks boshsettings.Networks) (err error) {
	fmt.Printf("Setup dhcp\n")
	return
}

func (p windowsPlatform) SetupManualNetworking(networks boshsettings.Networks) (err error) {
	fmt.Printf("Setup manual networking\n")
	return
}

func (p windowsPlatform) SetupLogrotate(groupName, basePath, size string) (err error) {
	fmt.Printf("Setup Logrotate\n")
	return
}

func (p windowsPlatform) SetTimeWithNtpServers(servers []string) (err error) {
	fmt.Printf("Set time with ntp servers\n")
	return
}

func (p windowsPlatform) SetupEphemeralDiskWithPath(devicePath string) (err error) {
	fmt.Printf("Setup ephemeral disk with path\n")
	return
}

func (p windowsPlatform) SetupDataDir() error {
	fmt.Printf("Setup data dir\n")
	return nil
}

func (p windowsPlatform) SetupTmpDir() error {
	fmt.Printf("Setup tmp dir\n")
	return nil
}

func (p windowsPlatform) MountPersistentDisk(devicePath, mountPoint string) (err error) {
	fmt.Printf("Mount persistent disk\n")
	return
}

func (p windowsPlatform) UnmountPersistentDisk(devicePath string) (didUnmount bool, err error) {
	fmt.Printf("Unmount persistent disk\n")
	return
}

func (p windowsPlatform) NormalizeDiskPath(attachment string) (devicePath string, found bool) {
	fmt.Printf("Normalize disk path\n")
	return "/dev/sdb", true
}

func (p windowsPlatform) GetFileContentsFromCDROM(filePath string) (contents []byte, err error) {
	fmt.Printf("Get file contents from cd\n")
	return
}

func (p windowsPlatform) MigratePersistentDisk(fromMountPoint, toMountPoint string) (err error) {
	fmt.Printf("Migrate persistent disk\n")
	return
}

func (p windowsPlatform) IsMountPoint(path string) (result bool, err error) {
	fmt.Printf("Is mount point\n")
	return
}

func (p windowsPlatform) IsPersistentDiskMounted(path string) (result bool, err error) {
	fmt.Printf("Is persistent disk mounted\n")
	return
}

func (p windowsPlatform) StartMonit() (err error) {
	fmt.Printf("Start monit ...\n")
	return
}

func (p windowsPlatform) SetupMonitUser() (err error) {
	fmt.Printf("Setup monit user\n")
	return
}

func (p windowsPlatform) GetMonitCredentials() (username, password string, err error) {
	fmt.Printf("Get monit credentials\n")
	return
}

func (p windowsPlatform) PrepareForNetworkingChange() error {
	fmt.Printf("prepare for networking change\n")
	return nil
}

func (p windowsPlatform) GetDefaultNetwork() (boshsettings.Network, error) {
	fmt.Printf("Get default network\n")
	var network boshsettings.Network

	networkPath := filepath.Join(p.dirProvider.BoshDir(), "windows-default-network-settings.json")
	contents, err := p.fs.ReadFile(networkPath)
	if err != nil {
		return network, nil
	}

	err = json.Unmarshal([]byte(contents), &network)
	if err != nil {
		return network, bosherr.WrapError(err, "Unmarshal json settings")
	}

	return network, nil
}
