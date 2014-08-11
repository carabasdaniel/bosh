package jobsupervisor

import (
	//	boshalert "bosh/agent/alert"
	bosherr "bosh/errors"
	boshlog "bosh/logger"
	boshdir "bosh/settings/directories"
	boshsys "bosh/system"

	"code.google.com/p/winsvc/mgr"
	"code.google.com/p/winsvc/svc"
	"encoding/json"
	"fmt"
	"github.com/pivotal/go-smtpd/smtpd"
	"time"
)

const jobSupervisorLogTag = "jobSupervisor"
const jobSupervisorPath = "C:\\sc_jobs\\jobs.json"

var serviceArguments = []string{}

type jobSupervisor struct {
	fs          boshsys.FileSystem
	runner      boshsys.CmdRunner
	logger      boshlog.Logger
	dirProvider boshdir.DirectoriesProvider

	jobFailuresServerPort int

	reloadOptions ReloadOptions
}

type ReloadOptions struct {
	// Number of times reload will be executed
	MaxTries int

	// Number of times monit incarnation will be checked
	// for difference after executing `monit reload`
	MaxCheckTries int

	// Length of time between checking for incarnation difference
	DelayBetweenCheckTries time.Duration
}

func NewJobSupervisor(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
	dirProvider boshdir.DirectoriesProvider,
	jobFailuresServerPort int,
	reloadOptions ReloadOptions,
) (js jobSupervisor) {
	return jobSupervisor{
		fs:          fs,
		runner:      runner,
		logger:      logger,
		dirProvider: dirProvider,

		jobFailuresServerPort: jobFailuresServerPort,

		reloadOptions: reloadOptions,
	}
}

func (js jobSupervisor) Reload() error {
	//this method was used for reloading monit
	return nil
}

func (js jobSupervisor) Start() error {
	jobs := ReadJobs(js.fs)

	for counter := 0; counter < len(jobs.Jobs); counter++ {
		name := jobs.Jobs[counter].Name
		preScript := jobs.Jobs[counter].PreStart

		if len(preScript) > 0 {
			_, stderr, exitcode, err := js.runner.RunCommand(preScript)

			if err != nil || exitcode != 0 {
				return bosherr.WrapError(err, fmt.Sprintf("Exit code: %d - Error output %s", exitcode, stderr))
			}
		}

		m, err := mgr.Connect()
		if err != nil {
			return bosherr.WrapError(err, "Connection error")
		}
		defer m.Disconnect()
		s, err := m.OpenService(name)
		if err != nil {
			return bosherr.WrapError(err, "Could not access service")
		}
		defer s.Close()
		err = s.Start(serviceArguments)
		if err != nil {
			return bosherr.WrapError(err, "could not start service")
		}
	}
	return nil
}

func (js jobSupervisor) Stop() error {
	jobs := ReadJobs(js.fs)

	for counter := 0; counter < len(jobs.Jobs); counter++ {
		name := jobs.Jobs[counter].Name

		preScript := jobs.Jobs[counter].PreStop

		if len(preScript) > 0 {
			_, stderr, exitcode, err := js.runner.RunCommand(preScript)

			if err != nil || exitcode != 0 {
				return bosherr.WrapError(err, fmt.Sprintf("Exit code: %d - Error output %s", exitcode, stderr))
			}
		}

		c := svc.Stop
		to := svc.Stopped
		m, err := mgr.Connect()
		if err != nil {
			return err
		}
		defer m.Disconnect()
		s, err := m.OpenService(name)
		if err != nil {
			return bosherr.WrapError(err, "could not access service: %v")
		}
		defer s.Close()
		status, err := s.Control(c)
		if err != nil {
			return bosherr.WrapError(err, "could not send stop: %v")
		}
		timeout := time.Now().Add(10 * time.Second)
		for status.State != to {
			if timeout.Before(time.Now()) {
				//"timeout waiting for service to go to state"
			}
			time.Sleep(300 * time.Millisecond)
			status, err = s.Query()
			if err != nil {
				return bosherr.WrapError(err, "could not retrieve service status: %v")
			}
		}
	}
	return nil
}

//Desired implementation
func (js jobSupervisor) Status() (status string) {
	jobs := ReadJobs(js.fs)

	status = "runnning"
	for counter := 0; counter < len(jobs.Jobs); counter++ {
		name := jobs.Jobs[counter].Name

		if jobs.Jobs[counter].Status != "monitored" {
			status = "failing"
		}

		m, err := mgr.Connect()
		if err != nil {
			status = "failing"
			return
		}
		defer m.Disconnect()
		s, err := m.OpenService(name)
		if err != nil {
			status = "unknown"
			return //"could not access service"
		}
		defer s.Close()
	}
	return status
}

func (js jobSupervisor) Unmonitor() error {
	jobs := ReadJobs(js.fs)

	for counter := 0; counter < len(jobs.Jobs); counter++ {
		jobs.Jobs[counter].Status = "not_monitored"
	}

	bytes, _ := json.Marshal(jobs.Jobs)
	js.fs.WriteFile(jobSupervisorPath, bytes)

	return nil
}

func (js jobSupervisor) AddPreStart(name, preStart string) {
	jobs := ReadJobs(js.fs)

	for counter := 0; counter < len(jobs.Jobs); counter++ {
		if jobs.Jobs[counter].Name == name {
			jobs.Jobs[counter].PreStart = preStart
			break
		}
	}

	bytes, _ := json.Marshal(jobs.Jobs)
	js.fs.WriteFile(jobSupervisorPath, bytes)
}

func (js jobSupervisor) AddPreStop(name, preStop string) {
	jobs := ReadJobs(js.fs)

	for counter := 0; counter < len(jobs.Jobs); counter++ {
		if jobs.Jobs[counter].Name == name {
			jobs.Jobs[counter].PreStop = preStop
			break
		}
	}

	bytes, _ := json.Marshal(jobs.Jobs)
	js.fs.WriteFile(jobSupervisorPath, bytes)
}

//TO DO: configPath treated as binPath
func (js jobSupervisor) AddJob(jobName string, jobIndex int, configPath string) error {
	jobs_list := ReadJobs(js.fs)
	//PreStart script path and PreStop script path can't be specified here
	newjob := Job{jobName, jobIndex, "monitored", "", configPath, ""}

	add_job := append(jobs_list, newjob)
	bytes, _ := json.Marshal(add_job)
	js.fs.WriteFile(jobSupervisorPath, bytes)

	return nil
}

func (js jobSupervisor) RemoveAllJobs() error {
	bytes, _ := json.Marshal(nil)
	js.fs.WriteFile(jobSupervisorPath, bytes)

	return nil
}

func (js jobSupervisor) MonitorJobFailures(handler JobFailureHandler) (err error) {
	alertHandler := func(smtpd.Connection, smtpd.MailAddress) (env smtpd.Envelope, err error) {
		//env = &alertEnvelope{
		//	new(smtpd.BasicEnvelope),
		//	handler,
		//	new(boshalert.MonitAlert),
		//}
		return
	}

	serv := &smtpd.Server{
		Addr:      fmt.Sprintf(":%d", js.jobFailuresServerPort),
		OnNewMail: alertHandler,
	}

	err = serv.ListenAndServe()
	if err != nil {
		err = bosherr.WrapError(err, "Listen for SMTP")
	}
	return
}
