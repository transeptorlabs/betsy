package utils

import "os"

// RemoveFile removes a file from the filesystem
func RemoveFile(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}
	return nil
}

// RemoveDevWallets removes the temp wallets from the filesystem
func RemoveDevWallets() error {
	if err := os.RemoveAll("./wallet/tmp"); err != nil {
		return err
	}
	return nil
}
