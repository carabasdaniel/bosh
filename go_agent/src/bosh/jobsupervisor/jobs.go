package jobsupervisor

import (
	bosherr "bosh/errors"
	boshsys "bosh/system"
	"encoding/json"
)

type Job struct {
	Name       string
	Index      int
	ConfigPath string
	State      string
	Status     string
}

type Jobs struct {
	Jobs []Job
}

func ReadJobs(fs boshsys.FileSystem) Jobs {
	var jobs_list Jobs
	in, err_read := fs.ReadFile(jobSupervisorPath)

	if err_read != nil {
		bosherr.WrapError(err_read, "Could not read the jobs file")
	}

	err_unmarshal := json.Unmarshal([]byte(in), &jobs_list.Jobs)

	if err_unmarshal != nil {
		bosherr.WrapError(err_unmarshal, "Could not unmarshal the jobs file")
	}

	return jobs_list
}
