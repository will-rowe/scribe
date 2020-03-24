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
	"fmt"
	"os"

	ipfs "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/will-rowe/scribe/src/backend"
	"github.com/will-rowe/scribe/src/config"
	"github.com/will-rowe/scribe/src/records"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <run|library>",
	Short: "Add a run or library to an existing project",
	Long: `Add a run or library to an existing project.
	
	This will collect the project (registered using scribe set --project XXX) and then add
specified record to the project, before committing it back to the IPFS.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAdd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// runAdd is the main block for the add subcommand
func runAdd(arg string) {

	// check the arg is a known record type
	switch arg {
	case "run", "library":
		break
	default:
		fmt.Printf("unrecognised argument (%v), use either run|library\n", arg)
		os.Exit(1)
	}

	// run the config checker to make sure we've got everything
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start the subcommand
	log.Info("------------SCRIBE------------")
	log.Info("starting the add subcommand...")
	log.Infof("\tconfig file: %v", viper.ConfigFileUsed())

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
		log.Infof("\tlaunching daemon...")
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
	log.Infof("\tnode identity: %v", nodeIdentity.ID)
	node.SetProject(config.Project)
	log.Infof("\tregistered node with proejct: %v", node.GetProject())

	// init the local project database
	log.Info("initialising the local project database...")
	db := records.InitDB()
	if len(config.RemoteCID) != 0 {
		log.Infof("\tremote CID found: %s", config.RemoteCID)
		log.Info("\tpulling project database from IPFS...")
		if err := db.Pull(node, config.RemoteCID); err != nil {
			log.Fatal(err)
		}
		log.Infof("\tnumber of projects added to local database: %d", db.GetNumProjects())
	} else {
		log.Info("\tno existing CID found")
	}

	// check for the required project
	log.Info("checking the local project database...")
	proj, err := db.GetProject(config.Project)
	switch err {
	case nil:
		log.Infof("\tproject found for %v", config.Project)

	case records.ErrNotFound:
		log.Infof("\tproject not found for %v", config.Project)
		log.Info("\tcreating project...")
		proj = records.InitProject(config.Project)
		log.Info("\tadding to the database...")
		if err := db.AddProject(proj); err != nil {
			log.Fatal(err)
		}
		log.Info("\tpushing database changes to IPFS...")
		cid, err := db.Push(node)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("\tupdating CID...")
		config.RemoteCID = cid
		viper.Set("remoteCID", cid)
		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
		log.Infof("\tCID updated: %s", config.RemoteCID)
		log.Infof("\tview on: %v", fmt.Sprintf("https://explore.ipld.io/#/explore/%s", config.RemoteCID))

	default:
		log.Fatal(err)
	}

	// work with the project
	log.Infof("\tproject loaded: %v", proj.GetLabel())

	// test message
	if err := node.Publish("just loaded the project over here..."); err != nil {
		log.Fatal(err)
	}

}
