package records

import (
	"testing"

	"github.com/will-rowe/scribe/src/backend"
)

// TestDAG
func TestDAG(t *testing.T) {

	// launch IPFS and then run the test
	backend.WithIpfs(t, 1, func(t *testing.T, APIaddress string) {

		// set up the node
		node, err := backend.InitNode(APIaddress)
		if err != nil {
			t.Fatal(err)
		}

		// init the db
		db := InitDB()

		// add in a project record
		testProject := InitProject("test project 1")
		if err := db.AddProject(testProject); err != nil {
			t.Fatal(err)
		}

		// push the db to the IPFS
		cid, err := db.Push(node)
		if err != nil {
			t.Fatal(err)
		}

		// here's the explorer link for debug
		t.Logf("IPLD Explorer link: https://explore.ipld.io/#/explore/%s \n", string(cid+"\n"))

		// pull the db from the IPFS
		db2 := InitDB()
		if err := db2.Pull(node, cid); err != nil {
			t.Fatal(err)
		}

		// check that the project has been retrieved
		retreivedProject, err := db2.GetProject("test project 1")
		if err != nil {
			t.Fatal(err)
		}
		if retreivedProject.GetLabel() != testProject.GetLabel() {
			t.Fatalf("mismatch between retrieved label and original: %s vs %s", retreivedProject.GetLabel(), testProject.GetLabel())
		}
	})
}
