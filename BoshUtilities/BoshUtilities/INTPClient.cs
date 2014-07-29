using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;

namespace BoshUtilities
{
    [ComVisible(true)]
    public interface INTPClient
    {
        void Init(string host);

        void Connect(bool updateSystemTime);
        string ToString();
    }
}
