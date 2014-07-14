using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Result JSON for handler
    /// </summary>
    public class HandlerResult
    {
        /// <summary>
        /// Gets or sets the result.
        /// </summary>
        /// <value>
        /// The result.
        /// </value>
        public string Result
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the time.
        /// </summary>
        /// <value>
        /// The time.
        /// </value>
        public DateTime Time
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the agent task id.
        /// </summary>
        /// <value>
        /// The agent task id.
        /// </value>
        public string AgentTaskId
        {
            get;
            set;
        }
    }
}
