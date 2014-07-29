using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Newtonsoft.Json;
using BoshUtilities;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Return JSON for a heartbeat message
    /// </summary>
    public class HeartbeatMessage
    {
        /// <summary>
        /// Gets or sets the job.
        /// </summary>
        /// <value>
        /// The job.
        /// </value>
        [JsonProperty("job")]
        public string Job { get; set; }

        /// <summary>
        /// Gets or sets the index.
        /// </summary>
        /// <value>
        /// The index.
        /// </value>
        [JsonProperty("index")]
        public int Index { get; set; }

        /// <summary>
        /// Gets or sets the state of the job.
        /// </summary>
        /// <value>
        /// The state of the job.
        /// </value>
        [JsonProperty("job_state")]
        public string JobState { get; set; }

        /// <summary>
        /// Gets or sets the vitals.
        /// </summary>
        /// <value>
        /// The vitals.
        /// </value>
        [JsonProperty("vitals")]
        public Vitals Vitals { get; set; }

        /// <summary>
        /// Gets or sets the NTP MSG.
        /// </summary>
        /// <value>
        /// The NTP MSG.
        /// </value>
        [JsonProperty("ntp")]
        public NtpMessage NtpMsg { get; set; }

       
    }
    /// <summary>
    /// Ntp message JSON
    /// </summary>
    public class NtpMessage
    {
        /// <summary>
        /// Gets or sets the offset.
        /// </summary>
        /// <value>
        /// The offset.
        /// </value>
        [JsonProperty("offset", NullValueHandling = NullValueHandling.Ignore)]
        public string Offset { get; set; }

        /// <summary>
        /// Gets or sets the timestamp.
        /// </summary>
        /// <value>
        /// The timestamp.
        /// </value>
        [JsonProperty("timestamp", NullValueHandling = NullValueHandling.Ignore)]
        public string Timestamp { get; set; }

        /// <summary>
        /// Gets or sets the message.
        /// </summary>
        /// <value>
        /// The message.
        /// </value>
        [JsonProperty("message", NullValueHandling=NullValueHandling.Ignore)]
        public string Message { get; set; }
    }
}
