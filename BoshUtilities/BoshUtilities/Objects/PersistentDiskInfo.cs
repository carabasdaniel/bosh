// -----------------------------------------------------------------------
// <copyright file="PersistentDiskInfo.cs" company="">
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
    public class PersistentDiskInfo
    {
        [JsonProperty("percent")]
        public string PersistentPercent { get { return persistentPercent; } set { persistentPercent = value; } }

        private string persistentPercent;
    }
}
