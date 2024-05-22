package mft

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CopyDirectory recursively copies a directory.
func (m *MFT) CopyDirectory(src, dst string) error {
	src = strings.TrimSuffix(src, "/")
	dst = strings.TrimSuffix(dst, "/")
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, src)
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		return m.CopyFile(path, dstPath)
	})
}

// DeleteFile deletes a file.
func (m *MFT) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// ReadFile reads the content of a file.
func (m *MFT) ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file.
func (m *MFT) WriteFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// ListFilesWithExtension lists all files with a specific extension in a directory.
func (m *MFT) ListFilesWithExtension(directoryPath, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// MoveFile moves a file from src to dst.
func (m *MFT) MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// CompareFiles compares the content of two files.
func (m *MFT) CompareFiles(file1, file2 string) (bool, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}
	return string(content1) == string(content2), nil
}

// GetFileModificationTime returns the last modification time of a file.
func (m *MFT) GetFileModificationTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

// SetFileModificationTime sets the last modification time of a file.
func (m *MFT) SetFileModificationTime(filePath string, modTime time.Time) error {
	return os.Chtimes(filePath, modTime, modTime)
}

// CalculateFileSize returns the size of a file in bytes.
func (m *MFT) CalculateFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// CheckFileExists checks if a file exists.
func (m *MFT) CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetAbsolutePath returns the absolute path of a file.
func (m *MFT) GetAbsolutePath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

// ListSubdirectories lists all subdirectories in a directory.
func (m *MFT) ListSubdirectories(directoryPath string) ([]string, error) {
	var dirs []string
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}

// GetFileExtension returns the file extension.
func (m *MFT) GetFileExtension(filePath string) string {
	return filepath.Ext(filePath)
}

// ChangeFileExtension changes the file extension.
func (m *MFT) ChangeFileExtension(filePath, newExtension string) (string, error) {
	newPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + newExtension
	err := os.Rename(filePath, newPath)
	return newPath, err
}

// ReadFileAsBytes reads the content of a file as bytes.
func (m *MFT) ReadFileAsBytes(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFileFromBytes writes bytes to a file.
func (m *MFT) WriteFileFromBytes(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

// CreateTempDirectory creates a temporary directory.
func (m *MFT) CreateTempDirectory(prefix string) (string, error) {
	return ioutil.TempDir("", prefix)
}

// MoveDirectory moves a directory from src to dst.
func (m *MFT) MoveDirectory(src, dst string) error {
	return os.Rename(src, dst)
}
