// Package backend interfaces with the Go IPFS API and enables pubsub for Scribe data. Inspiration taken from https://github.com/sahib/brig and https://github.com/planet-ethereum/relay-network
package backend

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
)

// PublishName will announce `name` to the network and make the node discoverable
func (node *Node) PublishName(name string) error {
	if !node.IsOnline() {
		return ErrOffline
	}
	fullName := "scribe:" + string(name)
	key, err := node.sh.BlockPut([]byte(fullName), "v0", MultiHash, -1)
	log.Debugf("published name: »%s« (key %s)", name, key)
	return err
}

// Identity gets the node's identity. It will cache the identity after the initial request
func (node *Node) Identity() (ipfs.PeerInfo, error) {
	node.Lock()

	// check the cached identity first
	if node.identity != "" {
		defer node.Unlock()
		return ipfs.PeerInfo{
			Addrs: []string{node.identity},
			ID:    node.identity,
		}, nil
	}

	// don't hold the lock during network operations
	node.Unlock()
	id, err := node.sh.ID()
	if err != nil {
		return ipfs.PeerInfo{}, err
	}
	node.Lock()
	node.identity = id.ID
	node.Unlock()
	return ipfs.PeerInfo{
		Addrs: []string{id.ID},
		ID:    id.ID,
	}, nil
}
