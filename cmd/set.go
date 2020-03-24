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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/will-rowe/scribe/src/config"
)

// set up the flags
var (
	reset      *bool
	echo       *bool
	ipfsPath   *string
	storageMax *string
	remoteCID  *string
	pinning    *bool
	project    *string
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config fields",
	Long: `This subcommand is used to set the config fields programmatically.

You can edit the config directly but this subcommand offers some checks for the fields being set.
An attempt will be made to create any directories that don't exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSet()
	},
}

// init the subcommand
func init() {
	rootCmd.AddCommand(setCmd)

	// local flags
	reset = setCmd.Flags().Bool("reset", false, "Reset the config to default by replacing any existing config file (any other command flags will be set afterwards)")
	echo = setCmd.Flags().Bool("echo", false, "Print the config to screen after setting values")
	ipfsPath = setCmd.Flags().String("ipfsPath", config.DefaultIpfsPath, "Path to the IPFS repository on this node")
	storageMax = setCmd.Flags().String("storageMax", config.DefaultStorageMax, "Maximum storage available for the IPFS repository")
	remoteCID = setCmd.Flags().String("remoteCID", "", "The CID of the remote project database")
	pinning = setCmd.Flags().Bool("pinning", true, "Pin IPFS objects (which will prevent local garabage collection)")
	project = setCmd.Flags().String("project", config.DefaultProject, "Project to operate on (add|update|listen)")

	// bind local flags to the config
	viper.BindPFlag("ipfsPath", setCmd.LocalFlags().Lookup("ipfsPath"))
	viper.BindPFlag("storageMax", setCmd.LocalFlags().Lookup("storageMax"))
	viper.BindPFlag("remoteCID", setCmd.LocalFlags().Lookup("remoteCID"))
	viper.BindPFlag("Pinning", setCmd.LocalFlags().Lookup("pinning"))
	viper.BindPFlag("project", setCmd.LocalFlags().Lookup("project"))
}

// runSet is the main block for the set subcommand
func runSet() {

	// reset takes precedence - it will replace any config currently on disk
	if *reset {
		configPath := fmt.Sprintf("%v/%v", viper.GetString("fileLocation"), viper.GetString("fileName"))
		if err := config.ResetConfig(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// re-init the Viper config
		initConfig()
	}

	// the in-memory Viper config will have picked up any set flags
	// run the config checker to make sure they are legit
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// now safe to update the config on disk
	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// echo if requested
	if *echo {
		dump, err := config.DumpConfig2JSON()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(dump)
	}
}
