package utils

import (
	"errors"
	"os"
)

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

// PadToBytes16 takes []byte and pads to [16]byte
func PadToBytes16(data []byte) ([]byte, error) {
	if len(data) > 16 {
		data = data[:16]
	} else if len(data) < 16 {
		padded := make([]byte, 16)
		copy(padded[16-len(data):], data)
		data = padded
	}

	if len(data) != 16 {
		return []byte{}, errors.New("padToBytes16(data) is not equal to 16 bytes")
	}

	return data, nil
}
