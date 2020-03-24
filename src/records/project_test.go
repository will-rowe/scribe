package records

import "testing"

var (
	projectLabel = "test project"
	projectCID   = ""
)

// TestProject
func TestProject(t *testing.T) {

	// init a project
	project := InitProject(projectLabel)

	_ = project
}
