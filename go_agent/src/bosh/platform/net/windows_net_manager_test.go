package net_test

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"

	boshlog "bosh/logger"
	. "bosh/platform/net"
	fakearp "bosh/platform/net/arp/fakes"
	fakenet "bosh/platform/net/fakes"
	"fmt"
	//boship "bosh/platform/net/ip"
	fakeip "bosh/platform/net/ip/fakes"
	boshsettings "bosh/settings"
	fakesys "bosh/system/fakes"
)

//TO DO: these are integration tests and should be used carefully that's why they are commented out
/func init() {
/*	Describe("windowsNetManager", func() {
		var (
			fs                     *fakesys.FakeFileSystem
			cmdRunner              *fakesys.FakeCmdRunner
			defaultNetworkResolver *fakenet.FakeDefaultNetworkResolver
			ipResolver             *fakeip.FakeIPResolver
			addressBroadcaster     *fakearp.FakeAddressBroadcaster
			netManager             NetManager
		)

		BeforeEach(func() {
			fs = fakesys.NewFakeFileSystem()
			cmdRunner = fakesys.NewFakeCmdRunner()
			defaultNetworkResolver = &fakenet.FakeDefaultNetworkResolver{}
			ipResolver = &fakeip.FakeIPResolver{}
			addressBroadcaster = &fakearp.FakeAddressBroadcaster{}
			logger := boshlog.NewLogger(boshlog.LevelNone)
			netManager = NewWindowsNetManager(
				fs,
				cmdRunner,
				defaultNetworkResolver,
				ipResolver,
				addressBroadcaster,
				logger,
			)
		})

		Describe("SetupDhcp", func() {
			It("Integration test of a DHCP setup call", func() {
				networks := boshsettings.Networks{
					"bosh": boshsettings.Network{
						Mac: "00:0C:29:04:B5:3B",
					},
				}
				errChan := make(chan error)
				netManager.SetupDhcp(networks, errChan)
				select {
				case err, ok := <-errChan:
					if !ok {
						fmt.Println("No error occured, error channel closed")
					} else {
						fmt.Println(err)
					}
				}
			})
		})

		Describe("SetupManualNetworking", func() {
			It("Integration test of a manual setup call", func() {
				networks := boshsettings.Networks{
					"bosh": boshsettings.Network{
						Mac:     "00:0C:29:04:B5:3B",
						IP:      "192.168.1.196",
						Netmask: "255.255.255.0",
						Gateway: "192.168.1.1",
						DNS:     []string{"8.8.8.8", "213.154.124.1", "193.231.252.1"},
					},
				}
				errChan := make(chan error)
				netManager.SetupManualNetworking(networks, errChan)
				select {
				case err, ok := <-errChan:
					if !ok {
						fmt.Println("No error occured, error channel closed")
					} else {
						fmt.Println(err)
					}
				}
			})
		})

	})
*/
}
