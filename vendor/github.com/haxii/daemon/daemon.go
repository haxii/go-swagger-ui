package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	service "github.com/takama/daemon"
)

// Daemon wrapper of takama's Daemon
type Daemon struct {
	srv service.Daemon

	usageName string
	usageMsg  string
}

const (
	// UsageMessage help to show the daemon's usage message
	UsageMessage string = "Send signal to a master process: install, remove, start, stop, status"
	// UsageDefaultName default signal is status
	UsageDefaultName string = "status"
)

// handle the command like
// exec -s install | remove | start | stop | status commands
func (d *Daemon) handle(entry func()) (string, error) {
	if len(os.Args) >= 2 && strings.EqualFold(os.Args[1], d.usageName) {
		// test the command, set to `status` as default
		command := UsageDefaultName
		if len(os.Args) > 2 {
			command = os.Args[2]
		}

		switch command {
		case "install":
			if len(os.Args) > 3 {
				return d.srv.Install(os.Args[3:]...)
			}
			return d.srv.Install()
		case "remove":
			return d.srv.Remove()
		case "start":
			return d.srv.Start()
		case "stop":
			return d.srv.Stop()
		case "status":
			return d.srv.Status()
		default:
			return d.usageMsg, nil
		}
	}

	entry()

	return d.usageMsg, nil
}

// Run run the daemon by specifing the entry point
func (d *Daemon) Run(entry func()) {
	// invoke the original entry point
	status, err := d.handle(entry)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		// show the running result
		fmt.Println(status)
	}
	os.Exit(1)
}

// Make make a new daemon.
//
// usageName is the sub command name of the service control for daemon and `-s`
// is recommended, serviceName & description is the name and description of the
// daemon, and dependencies are the services used by daemon
func Make(usageName string, serviceName,
	description string, dependencies ...string) *Daemon {
	// make a takama‘s daemon instance
	srv, err := service.New(serviceName, description, dependencies...)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	d := &Daemon{
		srv:       srv,
		usageName: usageName,
		usageMsg: "Usage: " + filepath.Base(os.Args[0]) + " " +
			usageName + " COMMAND" + "\n\n" + UsageMessage,
	}
	return d
}
