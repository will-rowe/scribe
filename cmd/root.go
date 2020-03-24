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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/will-rowe/scribe/src/config"
	"github.com/will-rowe/scribe/src/helpers"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scribe",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init the basics
func init() {

	// init the config
	cobra.OnInitialize(initConfig)

	// persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scribe)")
	rootCmd.PersistentFlags().Bool("private", false, "run scribe in private mode")

	// bind flags to the config
	viper.BindPFlag("private", rootCmd.PersistentFlags().Lookup("private")) // update config to private mode

	// set up logrus
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {

	// check if using default or user supplied config
	if cfgFile != "" {

		// only use existing user supplied config, we're not making it for them
		if !helpers.CheckFileExists(cfgFile) {
			fmt.Printf("supplied config does not exist (%v)", cfgFile)
			os.Exit(1)
		}

		// load it
		viper.SetConfigFile(cfgFile)

	} else {

		// generate the default config if we can't find it
		if !helpers.CheckFileExists(fmt.Sprintf("%v/%v", config.DefaultLocation, config.DefaultName)) {
			if err := config.GenerateDefault(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// load it
		viper.SetConfigType(config.DefaultType)
		viper.SetConfigName(config.DefaultName)
		viper.AddConfigPath(config.DefaultLocation)
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// read in the config - there shouldn't be a not found error so catch any other issue
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
