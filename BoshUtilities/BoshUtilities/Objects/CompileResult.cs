using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Newtonsoft.Json;

namespace Uhuru.BOSH.Agent.Objects
{
    public class CompileResult
    {
        /// <summary>
        /// Gets or sets the upload result
        /// </summary>
        [JsonProperty("result")]
        public UploadResult Result
        {
            get;
            set;
        }

    }

    /// <summary>
    /// Upload Result JSon
    /// </summary>
    public class UploadResult
    {
        /// <summary>
        /// Gets or sets the sha1.
        /// </summary>
        /// <value>
        /// The sha1.
        /// </value>
        [JsonProperty("sha1")]
        public string Sha1 { get; set; }

        /// <summary>
        /// Gets or sets the blobstore id.
        /// </summary>
        /// <value>
        /// The blobstore id.
        /// </value>
        [JsonProperty("blobstore_id")]
        public string BlobstoreId { get; set; }

        /// <summary>
        /// Gets or sets the compile log id.
        /// </summary>
        /// <value>
        /// The compile log id.
        /// </value>
        [JsonProperty("compile_log")]
        public string CompileLogId { get; set; }
    }
}
