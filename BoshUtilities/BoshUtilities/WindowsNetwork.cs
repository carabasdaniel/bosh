using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Globalization;
using System.Management;
using System.Threading;
using Newtonsoft.Json;
using System.Runtime.InteropServices;

namespace BoshUtilities
{
    [ComVisible(true)]
    public class WindowsNetwork:IWindowsNetwork
    {        
        //DNS must be comma separated string values server ips
        public void SetupNetwork(string MacAddress,string IP,string NetMask,string Gateway,string DNS)
        {
            string macAddress = MacAddress;
                Collection<string> macAddreses = GetExistingMacAddresses();

                if (macAddreses.Contains(macAddress.ToUpperInvariant()))
                {   
                    string ip = IP;
                    string netmask = NetMask;
                    SetIP(ip, netmask, macAddress);

                    string gateway = Gateway;
                    SetGateway(gateway, macAddress);

                    SetDNS(DNS, macAddress);                  
                }           
        }

        public void SetupDhcp(string MacAddress)
        {            
            int retryCount = 30;
            bool done=false;
            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                ;
                while (retryCount > 0)
                {
                    ManagementObjectCollection objMOC = objMC.GetInstances();

                    foreach (ManagementObject objMO in objMOC)
                    {
                        if ((bool)objMO["IPEnabled"])
                        {
                            if (objMO["MACAddress"].ToString().ToUpperInvariant() == MacAddress.ToUpperInvariant())
                            {
                                var ndns = objMO.GetMethodParameters("SetDNSServerSearchOrder");
                                ndns["DNSServerSearchOrder"] = null;
                                var enableDhcp = objMO.InvokeMethod("EnableDHCP", null, null);
                                var setDns= objMO.InvokeMethod("SetDNSServerSearchOrder", ndns, null);
                                done=true;
                            }
                            retryCount = 0;
                        }
                    }
                    if (!done)
                    {
                        Thread.Sleep(5000);
                        retryCount--;
                    }
                }
            }            
        }


       private Collection<string> GetExistingMacAddresses()
        {          
            Collection<string> macAddresses = new Collection<string>();
            int retryCount = 30;

            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                while (retryCount > 0)
                {
                    ManagementObjectCollection objMOC = objMC.GetInstances();

                    foreach (ManagementObject objMO in objMOC)
                    {
                        if ((bool)objMO["IPEnabled"])
                        {
                            macAddresses.Add(objMO["MACAddress"].ToString().ToUpperInvariant());
                            retryCount = 0;
                        }
                    }
                    if (macAddresses.Count == 0)
                    {
                        Thread.Sleep(5000);
                        retryCount--;
                    }
                }
            }
            return macAddresses;
        }


        private void SetIP(string ipAddress, string subnetMask, string macAddress)
        {
            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                ManagementObjectCollection objMOC = objMC.GetInstances();

                foreach (ManagementObject objMO in objMOC)
                {
                    if ((bool)objMO["IPEnabled"])
                    {
                        if (objMO["MACAddress"].ToString().ToUpperInvariant().Equals(macAddress.ToUpperInvariant()))
                        {
                          
                            try
                            {
                                ManagementBaseObject newIP =
                                    objMO.GetMethodParameters("EnableStatic");

                                newIP["IPAddress"] = new string[] { ipAddress };
                                newIP["SubnetMask"] = new string[] { subnetMask };

                                objMO.InvokeMethod("EnableStatic", newIP, null);
                            }
                            catch (Exception)
                            {
                                throw;
                            }
                        }
                    }
                }
            }
        }

        /// <summary>
        /// Set's a new Gateway address of the local machine
        /// </summary>
        /// <param name="gateway">The Gateway IP Address</param>
        /// <remarks>Requires a reference to the System.Management namespace</remarks>
        private void SetGateway(string gateway, string macAddress)
        {
            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                ManagementObjectCollection objMOC = objMC.GetInstances();

                foreach (ManagementObject objMO in objMOC)
                {
                    if ((bool)objMO["IPEnabled"])
                    {
                        if (objMO["MACAddress"].ToString().ToUpperInvariant().Equals(macAddress.ToUpperInvariant()))
                        {
                            try
                            {
                                ManagementBaseObject newGateway =
                                    objMO.GetMethodParameters("SetGateways");

                                newGateway["DefaultIPGateway"] = new string[] { gateway };
                                newGateway["GatewayCostMetric"] = new int[] { 1 };

                                objMO.InvokeMethod("SetGateways", newGateway, null);
                            }
                            catch (Exception)
                            {
                                throw;
                            }
                        }
                    }
                }
            }
        }

        /// <summary>
        /// Set's the DNS Server of the local machine
        /// </summary>
        /// <param name="NIC">NIC address</param>
        /// <param name="DNS">DNS server address</param>
        /// <remarks>Requires a reference to the System.Management namespace</remarks>
        private void SetDNS(string DNS, string macAddress)
        {
            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                ManagementObjectCollection objMOC = objMC.GetInstances();

                foreach (ManagementObject objMO in objMOC)
                {
                    if ((bool)objMO["IPEnabled"])
                    {
                        if (objMO["MACAddress"].ToString().ToUpperInvariant().Equals(macAddress.ToUpperInvariant()))
                        {
                           try
                            {
                                ManagementBaseObject newDNS =
                                    objMO.GetMethodParameters("SetDNSServerSearchOrder");
                                newDNS["DNSServerSearchOrder"] = DNS.Split(',');
                                objMO.InvokeMethod("SetDNSServerSearchOrder", newDNS, null);
                            }
                            catch (Exception)
                            {
                                throw;
                            }
                        }
                    }
                }
            }
        }
        /// <summary>
        /// Set's WINS of the local machine
        /// </summary>
        /// <param name="NIC">NIC Address</param>
        /// <param name="priWINS">Primary WINS server address</param>
        /// <param name="secWINS">Secondary WINS server address</param>
        /// <remarks>Requires a reference to the System.Management namespace</remarks>
        private void SetWINS(string priWINS, string secWINS, string macAddress)
        {
            using (ManagementClass objMC = new ManagementClass("Win32_NetworkAdapterConfiguration"))
            {
                ManagementObjectCollection objMOC = objMC.GetInstances();

                foreach (ManagementObject objMO in objMOC)
                {
                    if ((bool)objMO["IPEnabled"])
                    {
                        if (objMO["MACAddress"].ToString().ToUpperInvariant().Equals(macAddress.ToUpperInvariant()))
                        {
                            try
                            {
                                ManagementBaseObject wins =
                                objMO.GetMethodParameters("SetWINSServer");
                                wins.SetPropertyValue("WINSPrimaryServer", priWINS);
                                wins.SetPropertyValue("WINSSecondaryServer", secWINS);

                                objMO.InvokeMethod("SetWINSServer", wins, null);
                            }
                            catch (Exception)
                            {
                                throw;
                            }
                        }
                    }
                }
            }
        }
    }
}
