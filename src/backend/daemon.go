// Package backend interfaces with the Go IPFS API and enables pubsub for Scribe data. Inspiration taken from https://github.com/sahib/brig and https://github.com/planet-ethereum/relay-network
package backend

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/will-rowe/scribe/src/config"
)

var (
	apiAddress  = "127.0.0.1"
	apiPort     = 5001
	swarmPort   = 4001
	gatewayPort = 8081
)

// GetAPI returns the API address the daemon will serve the API from
func GetAPI() string {
	return fmt.Sprintf("%s:%d", apiAddress, apiPort)
}

// ConfigureDaemon will configure the IPFS daemon
func ConfigureDaemon(conf *config.ScribeConfig) error {

	// point IPFS to the repo location
	os.Setenv("IPFS_PATH", conf.IpfsPath)

	// init the daemon
	// this will return an error if daemon already inited, so ignore that for now
	cmd := exec.Command("ipfs", "init")
	cmd.Env = append(cmd.Env, fmt.Sprintf("IPFS_PATH=%s", conf.IpfsPath))
	cmd.Run()

	// set up the configure commands
	script := [][]string{
		{"ipfs", "config", "--json", "Experimental.Libp2pStreamMounting", "true"},
		{"ipfs", "config", "Datastore.StorageMax", conf.StorageMax},
		{"ipfs", "config", "Addresses.API", fmt.Sprintf("/ip4/%s/tcp/%d", apiAddress, apiPort)},
		{"ipfs", "config", "--json", "Addresses.Swarm", fmt.Sprintf("[\"/ip4/%s/tcp/%d\"]", apiAddress, swarmPort)},
		{"ipfs", "config", "Addresses.Gateway", fmt.Sprintf("/ip4/%s/tcp/%d", apiAddress, gatewayPort)},
	}

	// run the commands
	for _, line := range script {
		cmd := exec.Command(line[0], line[1:]...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("IPFS_PATH=%s", conf.IpfsPath))
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("could not configure the IPFS daemon (%v: %v)", err, strings.Join(line, " "))
		}
	}
	return nil
}

// LaunchDaemon will attempt to launch the IPFS daemon
// NOTE: daemon is left running and must be terminated by user
func LaunchDaemon(conf *config.ScribeConfig) error {

	// launch the daemon if it's not already running
	daemonCmd := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment")
	daemonCmd.Env = append(daemonCmd.Env, fmt.Sprintf("IPFS_PATH=%s", conf.IpfsPath))
	if err := daemonCmd.Start(); err != nil {
		return err
	}

	// wait until the daemon actually offers the API interface
	connected := false
	for tries := 0; tries < 200; tries++ {
		conn, err := net.Dial("tcp", GetAPI())
		if err == nil {
			conn.Close()
			connected = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !connected {
		return fmt.Errorf("IPFS daemon is not offering the API interface")
	}
	return nil
}
