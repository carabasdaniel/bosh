using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Job sincronized from the state yaml
    /// </summary>
    public class Job
    {
        /// <summary>
        /// Gets or sets the name of the job.
        /// </summary>
        /// <value>
        /// The name.
        /// </value>
        public string Name
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the version.
        /// </summary>
        /// <value>
        /// The version.
        /// </value>
        public string Version
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the sha1.
        /// </summary>
        /// <value>
        /// The sha1.
        /// </value>
        public string SHA1
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the template.
        /// </summary>
        /// <value>
        /// The template.
        /// </value>
        public string Template
        {
            get;
            set;
        }

        /// <summary>
        /// Gets or sets the blobstore_id.
        /// </summary>
        /// <value>
        /// The blobstore_id.
        /// </value>
        public string BlobstoreId
        {
            get;
            set;
        }
    }
}
