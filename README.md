# MFT Library Documentation

# MFTKIT

MFTKIT is a Go package providing a set of utility functions for file and directory management. This package includes functions for copying, moving, deleting, reading, and writing files and directories, as well as other file-related operations.

## Installation

To use the MFTKIT package, you need to install it using `go get`:

```sh
go get github.com/madhu72/mftkit
```

Then import it in your Go code:

```go
import "github.com/madhu72/mftkit/mft"
```

## Overview
The MFT library provides various functionalities for file management, including encryption, compression, splitting, merging, uploading, downloading, monitoring, and more. This documentation covers all the available methods and their usage.

## Methods

### EncryptFile
Encrypts a file using AES.
```go
func (m *MFT) EncryptFile(inputPath, outputPath, key string) error
```


### DecryptFile
Decrypts a file using AES.
```go
func (m *MFT) DecryptFile(inputPath, outputPath, key string) error
```

### CompressFile
Compresses a file using gzip.
```go
func (m *MFT) CompressFile(inputPath, outputPath string) error
```

### DecompressFile
Decompresses a gzip file.
```go
func (m *MFT) DecompressFile(inputPath, outputPath string) error
```

### CalculateChecksum
Calculates the SHA-256 checksum of a file.
```go
func (m *MFT) CalculateChecksum(filePath string) (string, error)
```

### SplitFile
Splits a file into parts of the specified size.
```go
func (m *MFT) SplitFile(filePath string, partSize int) ([]string, error)
```

### MergeFiles
Merges multiple file parts into a single file.
```go
func (m *MFT) MergeFiles(parts []string, outputPath string) error
```

### UploadFile
Uploads a file to a remote server.
```go
func (m *MFT) UploadFile(server, filePath, destinationPath string) error
```

### DownloadFile
Downloads a file from a remote server.
```go
func (m *MFT) DownloadFile(server, filePath, destinationPath string) error
```

### LogTransfer
Logs the file transfer action.
```go
func (m *MFT) LogTransfer(action, fileName, status string) error
```

### ValidateFile
Checks if a file exists and matches the given checksum.
```go
func (m *MFT) ValidateFile(filePath, checksum string) (bool, error)
```

### MonitorDirectory
Monitors a directory for changes and calls the callback function on each event.
```go
func (m *MFT) MonitorDirectory(directoryPath string, callback func(event FileEvent)) error
```

### SecureDelete
Securely deletes a file by overwriting its content.
```go
func (m *MFT) SecureDelete(filePath string) error
```

### CreateTempFile
Creates a temporary file with the given prefix and suffix.
```go
func (m *MFT) CreateTempFile(prefix, suffix string) (string, error)
```

### CleanUpTempFiles
Removes all temporary files created during the session.
```go
func (m *MFT) CleanUpTempFiles() error
```

### TrackTransferProgress
Tracks the progress of file transfer.
```go
func (m *MFT) TrackTransferProgress(filePath string, callback func(progress float64)) error
```

### SendTransferNotification
Sends a transfer notification.
```go
func (m *MFT) SendTransferNotification(event TransferEvent) error
```

### RegisterNotificationHandler
Registers a handler for transfer notifications.
```go
func (m *MFT) RegisterNotificationHandler(handler func(event TransferEvent))
```

### ScheduleFileTransfer
Schedules a file transfer.
```go
func (m *MFT) ScheduleFileTransfer(schedule Schedule) error
```

### SetTransferRateLimit
Sets the transfer rate limit.
```go
func (m *MFT) SetTransferRateLimit(bytesPerSecond int) error
```

### GetTransferRateLimit
Gets the current transfer rate limit.
```go
func (m *MFT) GetTransferRateLimit() (int, error)
```

### LockFile
Locks a file.
```go
func (m *MFT) LockFile(filePath string) error
```

### UnlockFile
Unlocks a file.
```go
func (m *MFT) UnlockFile(filePath string) error
```

### LogError
Logs an error with context information.
```go
func (m *MFT) LogError(err error, context string) error
```

### NormalizeFilePath
Normalizes a file path.
```go
func (m *MFT) NormalizeFilePath(filePath string) (string, error)
```

### ResolveConflict
Resolves a conflict between existing and new files.
```go
func (m *MFT) ResolveConflict(existingFilePath, newFilePath string) error
```

### AddFileDependency
Adds a dependency for a file.
```go
func (m *MFT) AddFileDependency(filePath, dependencyPath string) error
```

### RemoveFileDependency
Removes a dependency for a file.
```go
func (m *MFT) RemoveFileDependency(filePath, dependencyPath string) error
```

### RetrieveFileByContent
Retrieves a file by its content hash.
```go
func (m *MFT) RetrieveFileByContent(contentHash, searchDirectory string) (string, error)
```

### SanitizeFileData
Sanitizes the content of a file.
```go
func (m *MFT) SanitizeFileData(filePath string, sanitizationRules []SanitizationRule) error
```

### LogFileTransfer
Logs the file transfer action.
```go
func (m *MFT) LogFileTransfer(sourcePath, destinationPath, status string) error
```

### LoadConfiguration
Loads configuration from a file.
```go
func (m *MFT) LoadConfiguration(configFilePath string) (map[string]interface{}, error)
```

### SaveConfiguration
Saves configuration to a file.
```go
func (m *MFT) SaveConfiguration(configFilePath string, config map[string]interface{}) error
```

### GetCrossPlatformPath
Converts a file path to a cross-platform path.
```go
func (m *MFT) GetCrossPlatformPath(filePath string) (string, error)
```

### ChangeFileOwner
Changes the owner of a file.
```go
func (m *MFT) ChangeFileOwner(filePath, owner string) error
```

### GetFileOwner
Gets the owner of a file.
```go
func (m *MFT) GetFileOwner(filePath string) (string, error)
```

### MultiThreadedUpload
Uploads a file to a remote server using multiple threads.
```go
func (m *MFT) MultiThreadedUpload(filePath, destination string, numThreads int) error
```

### MultiThreadedDownload
Downloads a file from a remote server using multiple threads.
```go
func (m *MFT) MultiThreadedDownload(filePath, destination string, numThreads int) error
```

### TransformFile
Transforms a file's content using a provided function.
```go
func (m *MFT) TransformFile(inputPath, outputPath string, transformer func(data []byte) []byte) error
```

### GetFileAccessTime
Gets the last access time of a file.
```go
func (m *MFT) GetFileAccessTime(filePath string) (time.Time, error)
```

### SetFileAccessTime
Sets the last access time of a file.
```go
func (m *MFT) SetFileAccessTime(filePath string, accessTime time.Time) error
```

### RenameFile
Renames a file.
```go
func (m *MFT) RenameFile(oldPath, newPath string) error
```

### GetFileSize
Gets the size of a file.
```go
func (m *MFT) GetFileSize(filePath string) (int64, error)
```

### CreateDirectory
Creates a new directory.
```go
func (m *MFT) CreateDirectory(dirPath string) error
```

### DeleteDirectory
Deletes a directory.
```go
func (m *MFT) DeleteDirectory(dirPath string) error
```

### ListFiles
Lists all files in a directory.
```go
func (m *MFT) ListFiles(directoryPath string) ([]string, error)
```

### ListDirectories
Lists all directories in a directory.
```go
func (m *MFT) ListDirectories(directoryPath string) ([]string, error)
```

### ValidateFilePath
Checks if the file path exists and is accessible.
```go
func (m *MFT) ValidateFilePath(filePath string) (bool, error)
```

### Version
Gets the version of the MFT library.
```go
func (m *MFT) Version() string
```

### GetFileMetadata
Retrieves metadata for a file.
```go
func (m *MFT) GetFileMetadata(filePath string) (FileMetadata, error)
```

### SetFileMetadata
Sets metadata for a file.
```go
func (m *MFT) SetFileMetadata(filePath string, metadata FileMetadata) error
```

### SaveFileVersion
Saves a version of a file.
```go
func (m *MFT) SaveFileVersion(filePath string) (string, error)
```

### RevertToFileVersion
Reverts to a saved version of a file.
```go
func (m *MFT) RevertToFileVersion(filePath, versionID string) error
```

### SyncDirectories
Synchronizes the content of two directories.
```go
func (m *MFT) SyncDirectories(sourceDir, targetDir string) error
```

### CopyFile
Copies

a file from source to destination.
```go
func (m *MFT) CopyFile(src, dst string) error
```

### VerifyDataIntegrity
Verifies the integrity of a file using the specified hash type.
```go
func (m *MFT) VerifyDataIntegrity(filePath, hashType, expectedHash string) (bool, error)
```

### AddCustomProtocolHandler
Adds a handler for a custom protocol.
```go
func (m *MFT) AddCustomProtocolHandler(protocol string, handler func(request Request) (Response, error)) error
```

### HandleCustomProtocolRequest
Handles a custom protocol request.
```go
func (m *MFT) HandleCustomProtocolRequest(request Request) (Response, error)
```

### ArchiveFiles
Creates a zip archive from a list of files.
```go
func (m *MFT) ArchiveFiles(filePaths []string, archivePath string) error
```

### UnarchiveFile
Extracts a zip archive to the specified destination.
```go
func (m *MFT) UnarchiveFile(archivePath, destinationDir string) error
```

### SetFilePermissions
Sets the permissions for a file.
```go
func (m *MFT) SetFilePermissions(filePath string, permissions os.FileMode) error
```

### GetFilePermissions
Gets the permissions of a file.
```go
func (m *MFT) GetFilePermissions(filePath string) (os.FileMode, error)
```

### AddFileDependency2
Adds a dependency for a file.
```go
func (m *MFT) AddFileDependency2(filePath, dependencyPath string) error
```

### RemoveFileDependency2
Removes a dependency for a file.
```go
func (m *MFT) RemoveFileDependency2(filePath, dependencyPath string) error
```

### TrackTransferProgress2
Tracks the progress of a file transfer.
```go
func (m *MFT) TrackTransferProgress2(filePath string, callback func(progress float64)) error
```

## Structs

### FileEvent
Represents a file system event.
```go
type FileEvent struct {
	Op   fsnotify.Op
	Name string
}
```

### TransferEvent
Represents a file transfer event.
```go
type TransferEvent struct {
	Action   string
	FileName string
	Status   string
}
```

### Schedule
Represents a file transfer schedule.
```go
type Schedule struct {
	Time        time.Time
	FilePath    string
	Destination string
}
```

### FileMetadata
Represents metadata of a file.
```go
type FileMetadata struct {
	Size        int64
	Permissions os.FileMode
	ModTime     time.Time
	IsDir       bool
}
```

### Request
Represents a custom protocol request.
```go
type Request struct {
	Protocol string
	Data     []byte
}
```

### Response
Represents a custom protocol response.
```go
type Response struct {
	Status string
	Data   []byte
}
```

### SanitizationRule
Represents a rule for data sanitization.
```go
type SanitizationRule struct {
	Search  string
	Replace string
}
```

## Functions

### File Operations

- `CopyFile(src, dst string) error`
    - Copies a file from `src` to `dst`.

- `DeleteFile(filePath string) error`
    - Deletes a file at `filePath`.

- `ReadFile(filePath string) (string, error)`
    - Reads the content of a file and returns it as a string.

- `WriteFile(filePath, content string) error`
    - Writes the provided content to a file.

- `ListFilesWithExtension(directoryPath, extension string) ([]string, error)`
    - Lists all files with the specified extension in a directory.

- `MoveFile(src, dst string) error`
    - Moves a file from `src` to `dst`.

- `CompareFiles(file1, file2 string) (bool, error)`
    - Compares the content of two files and returns whether they are the same.

- `GetFileModificationTime(filePath string) (time.Time, error)`
    - Returns the last modification time of a file.

- `SetFileModificationTime(filePath string, modTime time.Time) error`
    - Sets the last modification time of a file.

- `CalculateFileSize(filePath string) (int64, error)`
    - Returns the size of a file in bytes.

- `CheckFileExists(filePath string) bool`
    - Checks if a file exists.

- `GetAbsolutePath(filePath string) (string, error)`
    - Returns the absolute path of a file.

### Directory Operations

- `CopyDirectory(src, dst string) error`
    - Recursively copies a directory from `src` to `dst`.

- `ListSubdirectories(directoryPath string) ([]string, error)`
    - Lists all subdirectories in a directory.

- `MoveDirectory(src, dst string) error`
    - Moves a directory from `src` to `dst`.

- `CreateTempDirectory(prefix string) (string, error)`
    - Creates a temporary directory with the given prefix.

- `ListAllFiles(directoryPath string) ([]string, error)`
  - Lists all files in a directory and its subdirectories.

- `CalculateDirectorySize(directoryPath string) (int64, error)`
  - Calculates the total size of a directory.

- `FindFilesByName(directoryPath, fileName string) ([]string, error)`
  - Finds files with a specific name in a directory and its subdirectories.

- `IsEmptyDirectory(directoryPath string) (bool, error)`
  - Checks if a directory is empty.

- `WatchDirectory(directoryPath string, callback func(event fsnotify.Event)) error`
  - Watches a directory for changes.
  - 
### File Extension Operations

- `GetFileExtension(filePath string) string`
    - Returns the file extension.

- `ChangeFileExtension(filePath, newExtension string) (string, error)`
    - Changes the file extension.

### File Content Operations

- `ReadFileAsBytes(filePath string) ([]byte, error)`
    - Reads the content of a file as bytes.

- `WriteFileFromBytes(filePath string, data []byte) error`
    - Writes bytes to a file.


## Usage Example

```go
package main

import (
    "fmt"
    "github.com/madhu72/mftkit/mft"
)

func main() {
    utils := mft.MFT{}
    
    // Copy a file
    err := utils.CopyFile("source.txt", "destination.txt")
    if err != nil {
        fmt.Println("Error copying file:", err)
    }

    // Read a file
    content, err := utils.ReadFile("destination.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
    } else {
        fmt.Println("File content:", content)
    }
    
    // Other operations...
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Author

[Venkateswara Rao T](https://github.com/madhu72)
```
