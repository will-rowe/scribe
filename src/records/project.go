// Package records interfaces the protobuf definitions with Scribe
package records

// InitProject will init a project struct with the minimum required values
func InitProject(label string) *Project {

	// create the project
	project := &Project{
		Label: label,
	}

	// return pointer to the project
	return project
}

// Register will register a project on the IPFS
func (project *Project) Register() error {

	return nil
}
