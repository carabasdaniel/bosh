using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Uhuru.BOSH.Agent.Objects
{
    /// <summary>
    /// Job manifest
    /// </summary>
    public class JobManifest
    {

        private Dictionary<string, string> templates = new Dictionary<string,string>() ;
        private ICollection<string> packages = new List<string>();

        /// <summary>
        /// Initializes a new instance of the <see cref="JobManifest"/> class.
        /// </summary>
        public JobManifest()
        {
        }

        /// <summary>
        /// Gets or sets the name.
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
        /// Gets or sets the templates.
        /// </summary>
        /// <value>
        /// The templates.
        /// </value>
        public Dictionary<string, string> Templates
        {
            get
            {
                return templates;
            }        
        }

        /// <summary>
        /// Adds the template.
        /// </summary>
        /// <param name="templateName">Name of the template.</param>
        /// <param name="templateValue">The template value.</param>
        public void AddTemplate(string templateName, string templateValue)
        {
            this.templates.Add(templateName, templateValue);
        }

        /// <summary>
        /// Gets or sets the packages.
        /// </summary>
        /// <value>
        /// The packages.
        /// </value>
        public ICollection<string> Packages
        {
            get
            {
                return packages;
            }
        }

        /// <summary>
        /// Adds the package.
        /// </summary>
        /// <param name="package">The package.</param>
        public void AddPackage(string package)
        {
            this.packages.Add(package);
        }
    }
}
