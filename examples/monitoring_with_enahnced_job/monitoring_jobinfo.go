package main

import (
	"drmaa2"
	"fmt"
)

func main() {
	// A Session Manager is always required
	var sm drmaa2.SessionManager

	// Open Monitoring Session and check for an error
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		fmt.Printf("Failed when creating monitoring session: %s\n", err)
	}
	// We need to close the Monitoring Session at exit
	defer ms.CloseMonitoringSession()

	// Create a JobInfo as filter with UNSET values. This is required
	// otherwise the filter will not work correctly.
	ji := drmaa2.CreateJobInfo()
	// We want to list all jobs submitted by that user.
	ji.JobOwner = "cbrunner"

	// The Univa Grid Engine DRMAA2 implementation is event based
	// meaning that GetAllJobs() can be called as often as required,
	// it does not trigger addional communication to the qmaster
	// process.
	if jobs, err := ms.GetAllJobs(&ji); err != nil {
		fmt.Printf("Error during GetAllJobs() call: %s\n", err)
	} else {
		for _, j := range jobs {
			fmt.Printf("Job id [%s] Job Name [%s]\n", j.GetId(), j.GetName())
		}

	}
}
