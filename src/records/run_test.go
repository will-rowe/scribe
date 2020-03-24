package records

import (
	"testing"
)

var (
	runLabel  = "test label"
	outputDir = "test output directory"
	fast5Dir  = "fast5s"
	fastqDir  = "fastqs"
)

// TestSync
func TestSync(t *testing.T) {
	run := InitRun(runLabel, outputDir, fast5Dir, fastqDir)

	// make sure you can't sync an orphan run
	//if err := run.Sync(); err == nil {
	//	t.Fatal("orphan run was synced")
	//}

	// add a parent project
	run.ParentProjectCID = projectCID

	// try syncing again
	if err := run.Sync(); err != nil {
		t.Fatal(err)
	}

}
