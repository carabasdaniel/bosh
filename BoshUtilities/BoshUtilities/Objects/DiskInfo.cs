// -----------------------------------------------------------------------
// <copyright file="DiskInfo.cs" company="">
// TODO: Update copyright text.
// </copyright>
// -----------------------------------------------------------------------

namespace Uhuru.BOSH.Agent.Objects
{
    using System;
    using System.Collections.Generic;
    using System.Linq;
    using System.Text;
    using Newtonsoft.Json;

    /// <summary>
    /// TODO: Update summary.
    /// </summary>
    public class DiskInfo
    {
        [JsonProperty("system")]
        public SystemDiskInfo SystemDisk { get { return systemDisk; } set { systemDisk = value; } }

        private SystemDiskInfo systemDisk;

        [JsonProperty("ephemeral")]
        public EphemeralDiskInfo EphemeralDisk { get { return ephemeralDisk; } set { ephemeralDisk = value; } }

        private EphemeralDiskInfo ephemeralDisk;

        [JsonProperty("persistent")]
        public PersistentDiskInfo PersistentDisk { get { return persistentDisk; } set { persistentDisk = value; } }

        private PersistentDiskInfo persistentDisk;
    }
}
