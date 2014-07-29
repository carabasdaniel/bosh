using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Text;
using Uhuru.BOSH.Agent.Objects;

namespace BoshUtilities
{    
        [JsonObject("vitals")]
        public class Vitals
        {
            [JsonProperty("load")]
            public Collection<string> Load { get; private set; }

            [JsonProperty("cpu")]
            public CPUInfo CPU { get { return cpu; } set { cpu = value; } }

            private CPUInfo cpu;

            [JsonProperty("mem")]
            public MemoryInfo Memory { get { return memory; } set { memory = value; } }

            private MemoryInfo memory;

            [JsonProperty("disk")]
            public DiskInfo Disk { get { return disk; } set { disk = value; } }

            private DiskInfo disk;

            public Vitals()
            {
                Load = new Collection<string>();
            }

            public void LoadAdd(string value)
            {
                Load.Add(value);
            }
        }    
}
