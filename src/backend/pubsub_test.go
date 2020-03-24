package backend

import (
	"context"
	"fmt"
	"testing"
	"time"

	ipfs "github.com/ipfs/go-ipfs-api"
)

func TestPubSub(t *testing.T) {

	// launch IPFS and then run the test
	WithIpfs(t, 1, func(t *testing.T, APIaddress string) {

		// init the node
		node, err := InitNode(APIaddress)
		if err != nil {
			t.Fatal(err)
		}

		// get the node's identity
		self, err := node.Identity()
		if err != nil {
			t.Fatal(err)
		}

		// subsribe the node
		ctx := context.Background()
		err = node.Subscribe(ctx, testProject)
		if err != nil {
			t.Fatal(err)
		}

		// collect all test errors via chan
		testErrs := make(chan error)

		// start the listener
		msgChan := make(chan *ipfs.Message)
		errChan := make(chan error, 1)
		sigChan := make(chan struct{})
		go node.Listen(msgChan, errChan, sigChan)

		// process any messages
		go func() {
			for {
				select {

				// collect any messages
				case msg := <-msgChan:

					// check the received message matches the sent one
					if testMessage != string(msg.Data) {
						testErrs <- fmt.Errorf("received message does not match the sent one")
					}

					// check the message source matches the sender address used
					if self.Addrs[0] != msg.From.Pretty() {
						testErrs <- fmt.Errorf("source address does not match sender address (%v vs %v)", self.Addrs, msg.From.Pretty())
					}

				// collect any errors
				case err := <-errChan:
					testErrs <- err
				}

			}
		}()

		// wait a second and then publish some test message to the network
		time.Sleep(1 * time.Second)
		go func() {
			if err := node.Publish(testMessage); err != nil {
				testErrs <- err
			}

			// one message is enough for test, signal the close
			close(sigChan)
			close(errChan)
		}()

		// check errors (channel closed after the test message is published)
		for err := range errChan {
			t.Fatal(err)
		}

		// unsubscribe the node
		if err := node.Unsubscribe(); err != nil {
			t.Fatal(err)
		}
	})
}
