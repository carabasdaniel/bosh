package net

import (

	//bosherr "bosh/errors"
	boshlog "bosh/logger"
	bosharp "bosh/platform/net/arp"
	boship "bosh/platform/net/ip"
	boshsettings "bosh/settings"
	boshsys "bosh/system"
	"strings"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

const windowsNetManagerLogTag = "windowsNetManager"

type windowsNetManager struct {
	DefaultNetworkResolver

	cmdRunner          boshsys.CmdRunner
	fs                 boshsys.FileSystem
	ipResolver         boship.IPResolver
	addressBroadcaster bosharp.AddressBroadcaster
	logger             boshlog.Logger
}

func NewWindowsNetManager(
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	defaultNetworkResolver DefaultNetworkResolver,
	ipResolver boship.IPResolver,
	addressBroadcaster bosharp.AddressBroadcaster,
	logger boshlog.Logger,
) windowsNetManager {
	return windowsNetManager{
		DefaultNetworkResolver: defaultNetworkResolver,
		cmdRunner:              cmdRunner,
		fs:                     fs,
		ipResolver:             ipResolver,
		addressBroadcaster:     addressBroadcaster,
		logger:                 logger,
	}
}

func (net windowsNetManager) getDNSServers(networks boshsettings.Networks) []string {
	dnsNetwork, _ := networks.DefaultNetworkFor("dns")
	return dnsNetwork.DNS
}

func (net windowsNetManager) SetupDhcp(networks boshsettings.Networks, errCh chan error) error {
	var isError bool
	for _, network := range networks {
		err := net.SetupDHCP(network.Mac)
		if err != nil {
			errCh <- err
		}
	}
	if isError {
		close(errCh)
	}
	return nil
}

func (net windowsNetManager) SetupManualNetworking(networks boshsettings.Networks, errCh chan error) error {
	var isError bool
	for _, network := range networks {
		dns := strings.Join(network.DNS, ",")
		err := net.SetupNetwork(network.Mac, network.IP, network.Netmask, network.Gateway, dns)
		if err != nil {
			isError = true
			errCh <- err
		}
	}
	if isError {
		close(errCh)
	}

	return nil
}

func (net windowsNetManager) SetupNetwork(macaddress, ip, netmask, gateway, dns string) error {

	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	defer ole.CoUninitialize()

	if err != nil {
		return err
	}

	unknown1, err := oleutil.CreateObject("BoshUtilities.WindowsNetwork")
	if err != nil {
		return err
	}

	cons1, err := unknown1.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}

	_, err = oleutil.CallMethod(cons1, "SetupNetwork", macaddress, ip, netmask, gateway, dns)

	if err != nil {
		return err
	}

	return nil
}

func (net windowsNetManager) SetupDHCP(macaddress string) error {

	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	defer ole.CoUninitialize()

	if err != nil {
		return err
	}

	unknown1, err := oleutil.CreateObject("BoshUtilities.WindowsNetwork")
	if err != nil {
		return err
	}

	cons1, err := unknown1.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}

	_, err = oleutil.CallMethod(cons1, "SetupDhcp", macaddress)

	if err != nil {
		return err
	}

	return nil

}
