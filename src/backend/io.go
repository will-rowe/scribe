// Package backend interfaces with the Go IPFS API and enables pubsub for Scribe data. Inspiration taken from https://github.com/sahib/brig and https://github.com/planet-ethereum/relay-network
package backend

import (
	"bytes"
	"fmt"

	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs-api/options"
)

// Add will add the content to the IPFS, pinning it if instructed
func (node *Node) Add(content []byte, pin bool) (string, error) {
	if !node.IsOnline() {
		return "", ErrOffline
	}
	return node.sh.Add(bytes.NewReader(content), ipfs.Pin(pin), ipfs.Hash(MultiHash), ipfs.CidVersion(1))
}

// Cat will return the data for a given CID in the IPFS
func (node *Node) Cat(cid string) ([]byte, error) {
	if !node.IsOnline() {
		return nil, ErrOffline
	}
	resp, err := node.sh.Cat(cid)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	data := buf.Bytes()
	err = nil
	err = resp.Close()
	return data, err
}

// DagPut wraps the DagPut API call
func (node *Node) DagPut(data []byte, encoding, format string, pin bool) (string, error) {
	if !node.IsOnline() {
		return "", ErrOffline
	}

	// need to convert the pin bool to a string
	pinVal := "false"
	if pin {
		pinVal = "true"
	}

	// call DagPut with options
	return node.sh.DagPutWithOpts(data,
		options.Dag.Pin(pinVal),
		options.Dag.InputEnc(encoding),
		options.Dag.Kind(format),
		options.Dag.Hash(MultiHash),
	)
}

// DagGet wraps the DagGet API call
func (node *Node) DagGet(cid, field string, output interface{}) error {
	if !node.IsOnline() {
		return ErrOffline
	}
	var ref string
	if len(field) != 0 {
		ref = fmt.Sprintf("%v/%v", cid, field)
	} else {
		ref = cid
	}
	return node.sh.DagGet(ref, output)
}
