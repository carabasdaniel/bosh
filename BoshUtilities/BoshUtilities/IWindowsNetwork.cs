using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;

namespace BoshUtilities
{
     [ComVisible(true)]
    public interface IWindowsNetwork
    {
         void SetupNetwork(string MacAddress, string IP, string NetMask, string Gateway, string DNS);
         void SetupDhcp(string MacAddress);
    }
}
