// Package backend interfaces with the Go IPFS API and enables pubsub for Scribe data. Inspiration taken from https://github.com/sahib/brig and https://github.com/planet-ethereum/relay-network
package backend

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// WithIpfs starts a new IPFS instance and calls `fn` with the API port to it.
// `portOff` is the offset to add on all standard ports.
func WithIpfs(t *testing.T, portOff int, fn func(t *testing.T, APIaddress string)) {

	// setup the IPFS
	ipfsPath, err := ioutil.TempDir("./", "test-scribe-backend")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(ipfsPath)

	// add some ports
	tgatewayPort := gatewayPort + portOff
	tswarmPort := swarmPort + portOff
	tapiPort := apiPort + portOff

	APIaddress := fmt.Sprintf("%s:%d", apiAddress, tapiPort)

	// run the IPFS daemon
	os.Setenv("IPFS_PATH", ipfsPath)
	script := [][]string{
		{"ipfs", "init"},
		{"ipfs", "config", "--json", "Experimental.Libp2pStreamMounting", "true"},
		{"ipfs", "config", "Addresses.API", fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", tapiPort)},
		{"ipfs", "config", "--json", "Addresses.Swarm", fmt.Sprintf("[\"/ip4/127.0.0.1/tcp/%d\"]", tswarmPort)},
		{"ipfs", "config", "Addresses.Gateway", fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", tgatewayPort)},
	}
	for _, line := range script {
		cmd := exec.Command(line[0], line[1:]...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("IPFS_PATH=%s", ipfsPath))
		err := cmd.Run()
		if err != nil {
			t.Fatalf("%v: %v", err, strings.Join(line, " "))
		}
	}
	daemonCmd := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment")
	daemonCmd.Env = append(daemonCmd.Env, fmt.Sprintf("IPFS_PATH=%s", ipfsPath))
	if err := daemonCmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := daemonCmd.Process.Kill(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait until the daemon actually offers the API interface
	connected := false
	for tries := 0; tries < 200; tries++ {
		conn, err := net.Dial("tcp", APIaddress)
		if err == nil {
			conn.Close()
			connected = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !connected {
		t.Fatal("could not connect to api")
	}

	// run the test
	fn(t, APIaddress)
}
