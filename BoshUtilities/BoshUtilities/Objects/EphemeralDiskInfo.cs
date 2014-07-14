// -----------------------------------------------------------------------
// <copyright file="EphemeralDiskInfo.cs" company="">
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
    public class EphemeralDiskInfo
    {
        [JsonProperty("percent")]
        public string Percent { get { return percent; } set { percent = value; } }

        private string percent;
    }
}
