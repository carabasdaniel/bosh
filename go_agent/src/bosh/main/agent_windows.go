// +build windows

package main

import (
	"code.google.com/p/winsvc/svc"
	"os"
)

type WindowsService struct {
}

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue}

	go runAgent()

loop:
	for {
		select {
		case change := <-r:
			switch change.Cmd {
			case svc.Interrogate:
				s <- change.CurrentStatus
			case svc.Stop, svc.Shutdown:
				{
					break loop
				}
			case svc.Pause:
				s <- svc.Status{State: svc.Paused, Accepts: svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue}
			case svc.Continue:
				s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue}
			default:
				{
					break loop
				}
			}
		}
	}
	s <- svc.Status{State: svc.StopPending}
	return
}

func main() {

	if contains(os.Args, "console") {
		runAgent()
	} else {
		var err error
		//setting stdout and stderr
		os.Stdout, err = os.Create("c:\\BoshGo\\log\\stdout.log")
		if err != nil {
			panic(err.Error())
		}

		os.Stderr, err = os.Create("c:\\BoshGo\\log\\stderr.log")
		if err != nil {
			panic(err.Error())
		}
		ws := WindowsService{}
		run := svc.Run

		err = run("boshagent", &ws)
		if err != nil {
			os.Exit(1)
		}
	}
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
