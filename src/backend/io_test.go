package backend

import (
	"encoding/json"
	"testing"
)

// TestIO
func TestIO(t *testing.T) {

	// launch IPFS and then run the test
	WithIpfs(t, 1, func(t *testing.T, APIaddress string) {

		// set up the node
		node, err := InitNode(APIaddress)
		if err != nil {
			t.Fatal(err)
		}

		// marshal test data to JSON
		testData := &testStruct{
			FieldA: fieldA,
			FieldB: fieldB,
		}
		data, err := json.Marshal(testData)
		if err != nil {
			t.Fatal(err)
		}

		// add to IPFS
		cid, err := node.Add(data, true)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("IPLD Explorer link: https://explore.ipld.io/#/explore/%s \n", string(cid+"\n"))

		// cat it from the IPFS
		retrievedData, err := node.Cat(cid)
		if err != nil {
			t.Fatal(err)
		}

		// unmarshal it into a new struct and check fields
		testCopy := &testStruct{}
		if err := json.Unmarshal(retrievedData, testCopy); err != nil {
			t.Fatal(err)
		}
		if (testCopy.FieldA != testData.FieldA) || (testCopy.FieldB != testData.FieldB) {
			t.Fatal("retrieved struct does not match original")
		}
	})
}
