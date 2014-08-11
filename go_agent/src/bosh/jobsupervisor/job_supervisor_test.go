package jobsupervisor_test

import (
	//"bytes"
	//"errors"
	//"fmt"
	//"net/smtp"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//boshalert "bosh/agent/alert"
	. "bosh/jobsupervisor"
	boshlog "bosh/logger"
	boshdir "bosh/settings/directories"
	fakesys "bosh/system/fakes"
)

var _ = Describe("jobSupervisor_windows", func() {
	var (
		fs                    *fakesys.FakeFileSystem
		runner                *fakesys.FakeCmdRunner
		logger                boshlog.Logger
		dirProvider           boshdir.DirectoriesProvider
		jobFailuresServerPort int
		jobsupervisor         JobSupervisor
	)

	//var jobFailureServerPort = 5000

	//getJobFailureServerPort := func() int {
	//	jobFailureServerPort++
	//	return jobFailureServerPort
	//}

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
		dirProvider = boshdir.NewDirectoriesProvider("C:\\fake-dir\\")

		jobsupervisor = NewJobSupervisor(
			fs,
			runner,
			logger,
			dirProvider,
			jobFailuresServerPort,
			ReloadOptions{
				MaxTries:               3,
				MaxCheckTries:          10,
				DelayBetweenCheckTries: 0 * time.Millisecond,
			},
		)
	})

	//doJobFailureEmail := func(email string, port int) error {
	//	conn, err := smtp.Dial(fmt.Sprintf("localhost:%d", port))
	//	for err != nil {
	//		conn, err = smtp.Dial(fmt.Sprintf("localhost:%d", port))
	//	}

	//	conn.Mail("sender@example.org")
	//	conn.Rcpt("recipient@example.net")
	//	writeCloser, err := conn.Data()
	//	if err != nil {
	//		return err
	//	}

	//	defer writeCloser.Close()

	//	buf := bytes.NewBufferString(fmt.Sprintf("%s\r\n", email))
	//	_, err = buf.WriteTo(writeCloser)
	//	if err != nil {
	//		return err
	//	}

	//	return nil
	//}

	//TO DO: Fake winsvc package for testing jobsupervisor start
	//Describe("Start", func() {
	//	It("start starts each monit service in group vcap", func() {
	//		//err := jobsupervisor.Start()
	//		//Expect(err).ToNot(HaveOccurred())
	//	})
	//})

	//TO DO: Fake winsvc package for testing jobsupervisor stop
	//Describe("Stop", func() {
	//	It("stop stops each monit service in group vcap", func() {
	//		//err := jobsupervisor.Stop()
	//		//Expect(err).ToNot(HaveOccurred())
	//	})
	//})

	Describe("AddJob", func() {
		fake := "[{\"Name\":\"LogRotator\",\"Index\":0,\"ConfigPath\":\"C:\\\\boshexecs\\\\logrotator\",\"State\":\"stopped\",\"Status\":\"monitored\"}]"

		It("test add job functionality", func() {
			err := jobsupervisor.AddJob("LogRotator", 0, "C:\\boshexecs\\logrotator")
			Expect(err).ToNot(HaveOccurred())

			writtenConfig, err := fs.ReadFileString("C:\\sc_jobs\\jobs3.json")

			Expect(err).ToNot(HaveOccurred())
			Expect(writtenConfig).To(Equal(fake))
		})
	})

	Describe("Status", func() {
		It("status returns unknown when no jobs", func() {
			status := jobsupervisor.Status()
			Expect("unknown").To(Equal(status))
		})

		It("status returns failing when service is not running", func() {
			jobsupervisor.AddJob("LogRotator", 0, "C:\\boshexecs\\logrotator")
			status := jobsupervisor.Status()
			Expect("failing").To(Equal(status))
		})
	})

	Describe("RemoveAllJobs", func() {
		Context("when jobs directory removal succeeds", func() {
			It("does not return error because all jobs are removed", func() {
				err := jobsupervisor.RemoveAllJobs()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Unmonitor", func() {
		fake := "[{\"Name\":\"LogRotator\",\"Index\":0,\"ConfigPath\":\"C:\\\\boshexecs\\\\logrotator\",\"State\":\"stopped\",\"Status\":\"not_monitored\"}]"
		It("testing setting jobs to unmonitored status", func() {
			jobsupervisor.AddJob("LogRotator", 0, "C:\\boshexecs\\logrotator")
			err := jobsupervisor.Unmonitor()
			writtenConfig, err := fs.ReadFileString("C:\\sc_jobs\\jobs3.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(writtenConfig).To(Equal(fake))
		})
	})

	Describe("MonitorJobFailures", func() {
		//It("monitor job failures", func() {
		//	var handledAlert boshalert.MonitAlert

		//	failureHandler := func(alert boshalert.MonitAlert) (err error) {
		//		handledAlert = alert
		//		return
		//	}

		//	go jobsupervisor.MonitorJobFailures(failureHandler)

		//	msg := `Message-id: <1304319946.0@localhost>
		//	Service: nats
		//	Event: does not exist
		//	Action: restart
		//	Date: Sun, 22 May 2011 20:07:41 +0500
		//	Description: process is not running`

		//	err := doJobFailureEmail(msg, jobFailuresServerPort)
		//	Expect(err).ToNot(HaveOccurred())

		//	Expect(handledAlert).To(Equal(boshalert.MonitAlert{
		//		ID:          "1304319946.0@localhost",
		//		Service:     "nats",
		//		Event:       "does not exist",
		//		Action:      "restart",
		//		Date:        "Sun, 22 May 2011 20:07:41 +0500",
		//		Description: "process is not running",
		//	}))
		//})

		//It("ignores other emails", func() {
		//	var didHandleAlert bool

		//	failureHandler := func(alert boshalert.MonitAlert) (err error) {
		//		didHandleAlert = true
		//		return
		//	}

		//	go jobsupervisor.MonitorJobFailures(failureHandler)

		//	//err := doJobFailureEmail(`fake-other-email`, jobFailuresServerPort)
		//	//Expect(err).ToNot(HaveOccurred())
		//	Expect(didHandleAlert).To(BeFalse())
		//})
	})
})
