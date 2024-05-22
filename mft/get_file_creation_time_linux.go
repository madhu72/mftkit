//go:build linux
// +build linux

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
	creationTime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	return creationTime, nil
}
