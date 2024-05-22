//go:build windows
// +build windows

package mft

import (
	"os"
	"syscall"
	"time"
)

// GetFileCreationTime gets the creation time of a file.
func (m *MFTUtils) GetFileCreationTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	stat := info.Sys().(*syscall.Win32FileAttributeData)
	creationTime := time.Unix(0, stat.CreationTime.Nanoseconds())
	return creationTime, nil
}
