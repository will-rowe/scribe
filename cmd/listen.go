/*
Copyright © 2020 Will Rowe <w.p.m.rowe@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	ipfs "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/will-rowe/scribe/src/backend"
	"github.com/will-rowe/scribe/src/config"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for project updates that are being pushed to the network",
	Long: `Listen for project updates that are being pushed to the network.
	
This command uses the pubsub protocol, which is currently an experimental IPFS
feature.`,
	Run: func(cmd *cobra.Command, args []string) {
		runListen()
	},
}

// init the subcommand
func init() {
	rootCmd.AddCommand(listenCmd)
}

// runListen is the main block for the listen subcommand
func runListen() {

	// run the config checker to make sure we've got everything
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start the subcommand
	log.Info("------------SCRIBE------------")
	log.Info("starting the listen subcommand...")
	log.Infof("\tconfig file: %v", viper.ConfigFileUsed())
	log.Info("setting config watcher...")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %v", e.Name)
	})

	// configure the daemon
	log.Info("configuring IPFS daemon...")
	config, err := config.DumpConfig2Mem()
	if err != nil {
		log.Fatal(err)
	}
	if err := backend.ConfigureDaemon(config); err != nil {
		log.Fatal(err)
	}
	log.Infof("\tupdated IPFS daemon using the scribe config")

	// check if the IPFS daemon is running (launch if not)
	tmpShell := ipfs.NewShell(backend.GetAPI())
	if !tmpShell.IsUp() {
		log.Infof("\tlaunching daemon")
		if err := backend.LaunchDaemon(config); err != nil {
			log.Fatal(err)
		}
	}
	log.Infof("\tdaemon is running")

	// init the node
	log.Info("initialising the node...")
	node, err := backend.InitNode(backend.GetAPI())
	if err != nil {
		log.Fatal(err)
	}
	nodeIdentity, err := node.Identity()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("\tAPI server listening on: %s", backend.GetAPI())
	log.Info("\tswarm and gateway info to go here...")
	log.Infof("\tnode identity: %v", nodeIdentity.ID)

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// subscribe to the requested project
	log.Info("subscribing the node...")
	if err := node.Subscribe(ctx, viper.GetString("project")); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := node.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Infof("\tlistening for: %v", viper.GetString("project"))

	// setup the pubsub listener
	msgChan := make(chan *ipfs.Message)
	errChan := make(chan error, 1)
	sigChan := make(chan struct{})
	go node.Listen(msgChan, errChan, sigChan)

	// catch the os interupt for graceful close down
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {

		// wait for interupt
		<-c
		log.Info("interrupt received - shutting down")

		// quit the PubSub listener
		close(sigChan)
		os.Exit(0)
	}()

	// process incoming messages
	for {
		select {

		// collect any messages
		case msg := <-msgChan:

			// check the message over
			log.Infof("\tmessage received from: %v", msg.From.Pretty())
			log.Infof("\tcontent: %v", msg.Data)

			// handle it
			//var doc map[string]interface{}
			//json.Unmarshal([]byte(s), &doc)
			//context, hasContext := doc["@context"]

		// collect any errors from the PubSub
		case err := <-errChan:
			log.Warn(err)
		}

	}
}
