package backend

import (
	"testing"
)

// testStruct
type testStruct struct {
	FieldA string `json:"fieldA"`
	FieldB int    `json:"fieldB"`
}

var (
	testProject = "test project"
	testMessage = "test message for the pubsub unit test"
	fieldA      = "test field"
	fieldB      = 666
)

// TestNode will test the node init
func TestNode(t *testing.T) {

	// launch IPFS and then run the test
	WithIpfs(t, 1, func(t *testing.T, APIaddress string) {
		node, err := InitNode(APIaddress)
		if err != nil {
			t.Fatal(err)
		}
		node.Disconnect()
		if node.IsOnline() {
			t.Fatal("node registered as online after disconnect")
		}
		node.Connect()
		if !node.IsOnline() {
			t.Fatal("node registered as offline after connect")
		}
	})
}
