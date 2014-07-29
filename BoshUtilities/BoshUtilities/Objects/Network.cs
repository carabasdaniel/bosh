using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Collections.ObjectModel;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Represents a network described in the state yaml
    /// </summary>
    public class Network
    {
        /// <summary>
        /// Gets or sets the ip.
        /// </summary>
        /// <value>
        /// The iP.
        /// </value>
        public string IP { get; set; }

        /// <summary>
        /// Gets or sets the name.
        /// </summary>
        /// <value>
        /// The name.
        /// </value>
        public string Name { get; set; }
    }
}
