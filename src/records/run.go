// Package records interfaces the protobuf definitions with Scribe
package records

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// InitRun will init a run struct with the minimum required values
func InitRun(label, outputDir, fast5Dir, fastqDir string) *Run {

	// create the run
	run := &Run{
		Created:              ptypes.TimestampNow(),
		Label:                label,
		History:              []*Comment{},
		Status:               1,
		Tags:                 make(map[string]bool),
		RequestOrder:         []string{},
		OutputDirectory:      outputDir,
		Fast5OutputDirectory: fast5Dir,
		FastqOutputDirectory: fastqDir,
	}

	// create the history
	run.AddComment("run created.")

	// return pointer to the run
	return run
}

// AddComment adds a comment to the run history
func (run *Run) AddComment(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("no comment provided")
	}
	comment := &Comment{
		Timestamp: ptypes.TimestampNow(),
		Text:      text,
	}
	run.History = append(run.History, comment)
	return nil
}

// Sync will store a run in the IPFS and update the parent project's record
func (run *Run) Sync() error {

	/*
		// check there is a registered parent project for this run
		if len(run.GetParentProjectCID()) == 0 {
			return fmt.Errorf("orphan run can't be synced - needs a parent project")
		}
	*/

	return nil
}
