//go:build darwin
// +build darwin

package mft

import (
	"os"
	"syscall"
	"time"
)

// GetFileCreationTime gets the creation time of a file.
func (m *MFT) GetFileCreationTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	stat := info.Sys().(*syscall.Stat_t)
	creationTime := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	return creationTime, nil
}
