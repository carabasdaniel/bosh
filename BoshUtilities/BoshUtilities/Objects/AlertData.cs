using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Newtonsoft.Json;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Class used for sending alert data through nats
    /// </summary>
    [JsonObject(MemberSerialization.OptIn)] 
    public class AlertData
    {
        /// <summary>
        /// Gets or sets the id.
        /// </summary>
        /// <value>
        /// The id.
        /// </value>
        [JsonProperty("id")]
        public string Id { get; set; }

        /// <summary>
        /// Gets or sets the severity.
        /// </summary>
        /// <value>
        /// The serverity.
        /// </value>
        [JsonProperty("severity")]
        public int? Severity { get; set; }

        /// <summary>
        /// Gets or sets the title.
        /// </summary>
        /// <value>
        /// The title.
        /// </value>
        [JsonProperty("title")]
        public string Title { get; set; }

        /// <summary>
        /// Gets or sets the summary.
        /// </summary>
        /// <value>
        /// The summary.
        /// </value>
        [JsonProperty("summary")]
        public string Summary { get; set; }

        /// <summary>
        /// Gets or sets the created at.
        /// </summary>
        /// <value>
        /// The created at.
        /// </value>
        [JsonProperty("created_at")]
        public DateTime? CreatedAt { get; set; }

    }
}
