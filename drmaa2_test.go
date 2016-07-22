// +build unit

package drmaa2

import (
	"fmt"
	"testing"
)

// Tests if MonitoringSession can be opened and closed.
// Requires the libdrmaa2.so in $LD_LIBRARY_PATH.
func TestOpenMonitoringSession(t *testing.T) {
	// Simple test for open and closing as MonitoringSession
	var sm SessionManager
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		t.Errorf("Couldn't open Monitoring session. %s", err)
		if ms != nil {
			t.Errorf("MonitoringSession needs to be nil in case of error")
		}
		return
	}
	t.Log("OpenMonitoringSession() created a MonitoringSession succesfully")
	if err := ms.CloseMonitoringSession(); err != nil {
		t.Errorf("CloseMonitoringSession() returned error: %s", err)
	}
}

func TestMonitoringSessionGetAllMachines(t *testing.T) {
	var sm SessionManager
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		t.Errorf("Couldn't open Monitoring session: %s", err)
		if ms != nil {
			t.Errorf("MonitoringSession needs to be nil in case of error")
		}
		return
	}
	// get all machines
	machine, err := ms.GetAllMachines(nil)
	if err != nil {
		t.Errorf("Error during GetAllMachines(nil): %s", err)
		return
	}
	amount := len(machine)
	if amount < 1 {
		t.Errorf("Error: No machine returned in GetAllMachines(nil)")
	}
	for _, m := range machine {
		t.Logf("Machine: %+v\n", m.Name)
	}
	// get a single machine
	var names []string
	names = append(names, machine[0].Name)
	if machine2, err := ms.GetAllMachines(names); err != nil {
		t.Errorf("Error in GetAllMachines(string): %s", err)
	} else {
		if len(machine2) != 1 {
			t.Error("Filter for machines in GetAllMachines([]string) seems not to work")
			return
		}
	}

	return
}

// TestReapJob tests job reaping from internal lists by calling the job's
// Reap() method.
func TestReapJob(t *testing.T) {
	var sm SessionManager
	var err error
	var js *JobSession

	// create or open job session
	if js, err = sm.OpenJobSession("TestReapJob"); err != nil {
		if js, err = sm.CreateJobSession("TestReapJob", ""); err != nil {
			t.Fatal(err)
		}
	}
	defer sm.DestroyJobSession("TestReapJob")

	var jt JobTemplate
	jt.RemoteCommand = "/bin/sleep"
	jt.Args = []string{"1"}

	jt.OutputPath = "/dev/null"
	jt.JoinFiles = true

	job, errRun := js.RunJob(jt)
	if errRun != nil {
		t.Fatal(errRun)
	}

	// wait until sleep is finished
	job.WaitTerminated(InfiniteTime)

	// it finsihed jobs appear in all job lists (GetJobs / monitoring session GetAllJobs)
	jl, errJL := js.GetJobs(nil)
	if errJL != nil {
		t.Fatalf("Error during GetJobs(): %s\n", errJL)
	}

	if len(jl) != 1 {
		t.Logf("Job list must be 1 but it is %d\n", len(jl))
	}

	if errReap := job.Reap(); errReap != nil {
		t.Fatalf("Reaping of job caused an error: %s\n", errReap)
	}
	t.Log("Reaping of job successful")

	jl, errJL = js.GetJobs(nil)

	if len(jl) != 0 {
		t.Fatalf("Job list still contains reaped jobs: %d\n", len(jl))
	}
}

func TestGetJobTemplate(t *testing.T) {
	var jt JobTemplate
	var sm SessionManager
	var js *JobSession
	var err error

	if js, err = sm.OpenJobSession("TestGetJobTemplate"); err != nil {
		if js, err = sm.CreateJobSession("TestGetJobTemplate", ""); err != nil {
			t.Errorf("Failed when creating job session: %s\n", err)
			return
		}
	}
	defer sm.DestroyJobSession("TestGetJobTemplate")

	jt.JobEnvironment = make(map[string]string, 0)
	jt.JobEnvironment["one"] = "1"
	jt.JobEnvironment["two"] = "2"
	jt.JobEnvironment["tree"] = "3"

	jt.RemoteCommand = "/bin/sleep"
	jt.Args = []string{"0"}

	job, _ := js.RunJob(jt)
	template, _ := job.GetJobTemplate()

	env := template.JobEnvironment

	for k, v := range jt.JobEnvironment {
		if env[k] != v {
			t.Errorf("JobEnvironment is not correctly recovered (%s != %s)\n", env[k], v)
		} else {
			fmt.Println("Found environment variable.")
		}
	}
}

// TODO add more :)
