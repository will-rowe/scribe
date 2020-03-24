// Package config is used to set up and manage the Scribe config file
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/will-rowe/scribe/src/helpers"
)

var (

	// DefaultName for the config file
	DefaultName = ".scribe.json"

	// DefaultLocation for the config file
	DefaultLocation = ""

	// DefaultType of file for the config file
	DefaultType = "json"

	// DefaultLicense for any file being created
	DefaultLicense = "MIT"

	// DefaultIpfsPath for storing IPFS files on this node
	DefaultIpfsPath = ""

	// DefaultStorageMax is the maximum storage available for the IPFS repo
	DefaultStorageMax = "1GB"

	// DefaultProject to operate on
	DefaultProject = "scribe-test-project"
)

// ScribeConfig is a struct to hold the config data
type ScribeConfig struct {
	FileName     string `json:"fileName"`
	FileLocation string `json:"fileLocation"`
	FileType     string `json:"fileType"`
	License      string `json:"license"`
	Private      bool   `json:"private"`
	IpfsPath     string `json:"ipfsPath"`
	StorageMax   string `json:"storageMax"`
	Pinning      bool   `json:"pinning"`
	RemoteCID    string `json:"remoteCID"`
	Project      string `json:"project"`
}

// init the default config filepaths
func init() {
	DefaultLocation, _ = homedir.Dir()
	DefaultIpfsPath = fmt.Sprintf("%v/.ipfs", DefaultLocation)
}

// GenerateDefault will generate the default config on disk
func GenerateDefault() error {

	// set up the default config data
	defaultConfig := &ScribeConfig{
		FileName:     DefaultName,
		FileLocation: DefaultLocation,
		FileType:     DefaultType,
		License:      DefaultLicense,
		Private:      false,
		IpfsPath:     DefaultIpfsPath,
		StorageMax:   DefaultStorageMax,
		Pinning:      false,
		RemoteCID:    "",
		Project:      DefaultProject,
	}

	// create the file
	fh, err := os.Create(fmt.Sprintf("%v/%v", DefaultLocation, DefaultName))
	defer fh.Close()

	// marshal the config
	d, err := json.MarshalIndent(defaultConfig, "", "\t")
	if err != nil {
		return err
	}

	// write the file
	_, err = fh.Write(d)
	return err
}

// ResetConfig will remove any existing config and replace it with the default one
// NOTE: the caller must reload the config into viper
func ResetConfig(configPath string) error {

	// remove the existing config if it exists
	if helpers.CheckFileExists(configPath) {
		if err := os.Remove(configPath); err != nil {
			return err
		}
	}

	// now generate the default and write it to disk
	return GenerateDefault()
}

// DumpConfig2Mem will unmarshall the config from Viper to a struct in memory
func DumpConfig2Mem() (*ScribeConfig, error) {
	c := &ScribeConfig{}
	err := viper.Unmarshal(c)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return c, nil
}

// DumpConfig2JSON will unmarshall the config from Viper to a JSON string
func DumpConfig2JSON() (string, error) {
	c, err := DumpConfig2Mem()
	if err != nil {
		return "", err
	}
	d, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// CheckConfig will check the fields of the in-memory Viper config
// it will attempt to make any directories that don't exist
func CheckConfig() error {

	// TODO: check for the config file? not actually needed

	// check the ipfs path exists, try making it if needed
	if err := helpers.CheckDirExists(viper.GetString("ipfsPath")); err != nil {
		if err := os.Mkdir(viper.GetString("ipfsPath"), 0755); err != nil {
			return fmt.Errorf("can't create new directory for IPFS (%v)", err)
		}
	}

	// TODO: add more checks as we work on the config

	return nil
}
