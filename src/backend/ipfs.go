// Package backend interfaces with the Go IPFS API and enables pubsub for Scribe data. Inspiration taken from https://github.com/sahib/brig and https://github.com/planet-ethereum/relay-network
package backend

import (
	"context"
	"errors"
	"fmt"
	"sync"

	ipfs "github.com/ipfs/go-ipfs-api"
)

// MultiHash is the multihash type used by IPFS
const MultiHash string = "sha2-256"

// error messages
var (
	// ErrOffline is returned by operations that need online support to work when a node is offline
	ErrOffline   = errors.New("node is offline")
	ErrNoProject = errors.New("node has no registered project")
)

// Node wraps the IPFS shell
type Node struct {
	sync.Mutex
	sh           *ipfs.Shell              // the IPFS shell
	allowNetwork bool                     // controls the node's network connection
	identity     string                   // the node's identifier
	subscription *ipfs.PubSubSubscription // the PubSub subscription for this node
	project      string                   // the PubSub project for this node
}

// InitNode will returns a HTTP based IPFS node
func InitNode(ipfsAPIendpoint string) (*Node, error) {

	// use the API endpoint at the specified port
	newNode := &Node{
		sh:           ipfs.NewShell(ipfsAPIendpoint),
		allowNetwork: true,
		identity:     "",
	}

	// check it's online
	if !newNode.IsOnline() {
		return nil, ErrOffline
	}

	return newNode, nil
}

// SetProject will register the node with a project
func (node *Node) SetProject(project string) {
	node.project = project
	return
}

// GetProject will return the registered project for the node
func (node *Node) GetProject() string {
	return node.project
}

// IsOnline returns true if the node is in online mode and the IPFS daemon is reachable
func (node *Node) IsOnline() bool {
	node.Lock()
	allowNetwork := node.allowNetwork
	node.Unlock()
	return node.sh.IsUp() && allowNetwork
}

// Connect allows the node to connect to the network
func (node *Node) Connect() error {
	node.Lock()
	defer node.Unlock()
	node.allowNetwork = true
	return nil
}

// Disconnect prevents the node from connecting to the network
func (node *Node) Disconnect() error {
	node.Lock()
	defer node.Unlock()
	node.allowNetwork = false
	return nil
}

// Subscribe will subscribe the node to a project
func (node *Node) Subscribe(ctx context.Context, project string) error {
	if !node.IsOnline() {
		return ErrOffline
	}
	sub, err := node.sh.PubSubSubscribe(project)
	if err != nil {
		return err
	}
	node.subscription = sub
	node.project = project
	return nil
}

// Unsubscribe will unsubscribe the node
func (node *Node) Unsubscribe() error {
	if !node.IsOnline() {
		return ErrOffline
	}
	if node.project == "" {
		return nil
	}
	err := node.subscription.Cancel()
	if err != nil {
		return err
	}
	node.subscription = nil
	node.project = ""
	return nil
}

// Publish will publish a message about the registered project
func (node *Node) Publish(message string) error {
	if !node.IsOnline() {
		return ErrOffline
	}
	if len(node.project) == 0 {
		return ErrNoProject
	}
	return node.sh.PubSubPublish(node.project, message)
}

// Listen will wait for messages on the PubSub subscription
//
// msgChan is to send back received messages to the caller
// errChan is to send back any errors to the the caller
// sigChan is to terminate the listener
func (node *Node) Listen(msgChan chan *ipfs.Message, errChan chan error, sigChan chan struct{}) {
	for {
		select {
		case <-sigChan:
			return
		default:
			message, err := node.subscription.Next()
			if err != nil {
				errChan <- fmt.Errorf("failed to wait for pubsub message: %v", err)
			} else {
				msgChan <- message
			}
		}
	}
}
