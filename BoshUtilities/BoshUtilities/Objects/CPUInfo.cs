// -----------------------------------------------------------------------
// <copyright file="CPUInfo.cs" company="">
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
    public class CPUInfo
    {
        [JsonProperty("user")]
        public string User { get { return user; } set { user = value; } }

        private string user;

        [JsonProperty("sys")]
        public string Sys { get { return sys; } set { sys = value; } }

        private string sys;

        [JsonProperty("wait")]
        public string Wait { get { return wait; } set { wait = value; } }

        private string wait;
    }
}
