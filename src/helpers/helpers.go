// Package helpers contains some useful helper functions for scribe
package helpers

import (
	"fmt"
	"os"
)

// CheckDirExists checks a directory exists
func CheckDirExists(dirPath string) error {
	fh, err := os.Stat(dirPath)
	switch {
	case err != nil:
		return err
	case fh.IsDir():
		return nil
	default:
		return fmt.Errorf("not a directory (%v)", dirPath)
	}
}

// CheckFileExists checks a returns true if a file exists and is not a directory
func CheckFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
