package mft

import (
	"archive/zip"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type MFT struct {
}

func NewMFT() *MFT {
	return &MFT{}
}

// EncryptFile encrypts a file using AES.
func (m *MFT) EncryptFile(inputPath, outputPath, key string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := outputFile.Write(iv); err != nil {
		return err
	}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

// DecryptFile decrypts a file using AES.
func (m *MFT) DecryptFile(inputPath, outputPath, key string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}

// CompressFile compresses a file using gzip.
func (m *MFT) CompressFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := gzip.NewWriter(outputFile)
	defer writer.Close()

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

// DecompressFile decompresses a gzip file.
func (m *MFT) DecompressFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader, err := gzip.NewReader(inputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}

func (m *MFT) CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SplitFile splits a file into parts of the specified size.
func (m *MFT) SplitFile(filePath string, partSize int) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	var parts []string
	buffer := make([]byte, partSize)
	for i := 0; ; i++ {
		bytesRead, err := inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}

		partPath := filePath + ".part" + string(i)
		outputFile, err := os.Create(partPath)
		if err != nil {
			return nil, err
		}
		defer outputFile.Close()

		if _, err := outputFile.Write(buffer[:bytesRead]); err != nil {
			return nil, err
		}

		parts = append(parts, partPath)
	}

	return parts, nil
}

// MergeFiles merges multiple file parts into a single file.
func (m *MFT) MergeFiles(parts []string, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for _, partPath := range parts {
		inputFile, err := os.Open(partPath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		if _, err := io.Copy(outputFile, inputFile); err != nil {
			return err
		}
	}

	return nil
}

// UploadFile uploads a file to a remote server.
func (m *MFT) UploadFile(server, filePath, destinationPath string) error {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(conn, file); err != nil {
		return err
	}

	return nil
}

// DownloadFile downloads a file from a remote server.
func (m *MFT) DownloadFile(server, filePath, destinationPath string) error {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}
	defer conn.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, conn); err != nil {
		return err
	}

	return nil
}

// LogTransfer logs the file transfer action.
func (m *MFT) LogTransfer(action, fileName, status string) error {
	file, err := os.OpenFile("transfer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("%s: %s - %s\n", action, fileName, status)

	return nil
}

// ValidateFile checks if a file exists and matches the given checksum.
func (m *MFT) ValidateFile(filePath, checksum string) (bool, error) {
	calculatedChecksum, err := m.CalculateChecksum(filePath)
	if err != nil {
		return false, err
	}
	return calculatedChecksum == checksum, nil
}

// FileEvent represents a file system event.
type FileEvent struct {
	Op   fsnotify.Op
	Name string
}

// MonitorDirectory monitors a directory for changes and calls the callback function on each event.
func (m *MFT) MonitorDirectory(directoryPath string, callback func(event FileEvent)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(directoryPath)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			callback(FileEvent{Op: event.Op, Name: event.Name})
		case err := <-watcher.Errors:
			return err
		}
	}
}

// SecureDelete securely deletes a file by overwriting its content.
func (m *MFT) SecureDelete(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	data := make([]byte, info.Size())
	if _, err := file.Write(data); err != nil {
		return err
	}
	return os.Remove(filePath)
}

// CreateTempFile creates a temporary file with the given prefix and suffix.
func (m *MFT) CreateTempFile(prefix, suffix string) (string, error) {
	tempFile, err := ioutil.TempFile("", prefix+"*"+suffix)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	return tempFile.Name(), nil
}

// CleanUpTempFiles removes all temporary files created during the session.
func (m *MFT) CleanUpTempFiles() error {
	tempDir := os.TempDir()
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(filepath.Join(tempDir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

// TrackTransferProgress tracks the progress of file transfer.
func (m *MFT) TrackTransferProgress(filePath string, callback func(progress float64)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	totalSize := fileInfo.Size()
	buffer := make([]byte, 1024)
	var bytesRead int64

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		bytesRead += int64(n)
		progress := float64(bytesRead) / float64(totalSize) * 100
		callback(progress)
	}

	return nil
}

// TransferEvent represents a file transfer event.
type TransferEvent struct {
	Action   string
	FileName string
	Status   string
}

// SendTransferNotification sends a transfer notification.
func (m *MFT) SendTransferNotification(event TransferEvent) error {
	// Replace this with your actual notification logic
	fmt.Printf("Notification: %s - %s (%s)\n", event.Action, event.FileName, event.Status)
	return nil
}

// RegisterNotificationHandler registers a handler for transfer notifications.
func (m *MFT) RegisterNotificationHandler(handler func(event TransferEvent)) {
	// Replace this with your actual registration logic
	// This example simply calls the handler directly
	handler(TransferEvent{Action: "Upload", FileName: "example.txt", Status: "Success"})
}

// Schedule represents a file transfer schedule.
type Schedule struct {
	Time        time.Time
	FilePath    string
	Destination string
}

// ScheduleFileTransfer schedules a file transfer.
func (m *MFT) ScheduleFileTransfer(schedule Schedule) error {
	duration := time.Until(schedule.Time)
	time.AfterFunc(duration, func() {
		err := m.UploadFile("server_address", schedule.FilePath, schedule.Destination)
		if err != nil {
			fmt.Println("Scheduled transfer failed:", err)
		} else {
			fmt.Println("Scheduled transfer succeeded")
		}
	})
	return nil
}

var transferRateLimit int = 1024 // Default limit in bytes per second

// SetTransferRateLimit sets the transfer rate limit.
func (m *MFT) SetTransferRateLimit(bytesPerSecond int) error {
	transferRateLimit = bytesPerSecond
	return nil
}

// GetTransferRateLimit gets the current transfer rate limit.
func (m *MFT) GetTransferRateLimit() (int, error) {
	return transferRateLimit, nil
}

// LockFile locks a file.
func (m *MFT) LockFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}

// UnlockFile unlocks a file.
func (m *MFT) UnlockFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

// LogError logs an error with context information.
func (m *MFT) LogError(err error, context string) error {
	file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("Error: %s - %v\n", context, err)

	return nil
}

// NormalizeFilePath normalizes a file path.
func (m *MFT) NormalizeFilePath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

// ResolveConflict resolves a conflict between existing and new files.
func (m *MFT) ResolveConflict(existingFilePath, newFilePath string) error {
	return os.Rename(newFilePath, existingFilePath)
}

// ResolveConflict resolves a file conflict by renaming the new file.
func (m *MFT) ResolveConflict2(existingFilePath, newFilePath string) error {
	backupPath := existingFilePath + ".bak"
	if err := os.Rename(existingFilePath, backupPath); err != nil {
		return err
	}
	return os.Rename(newFilePath, existingFilePath)
}

var fileDependencies = make(map[string][]string)

// AddFileDependency adds a dependency for a file.
func (m *MFT) AddFileDependency(filePath, dependencyPath string) error {
	fileDependencies[filePath] = append(fileDependencies[filePath], dependencyPath)
	return nil
}

// RemoveFileDependency removes a dependency for a file.
func (m *MFT) RemoveFileDependency(filePath, dependencyPath string) error {
	dependencies := fileDependencies[filePath]
	for i, dep := range dependencies {
		if dep == dependencyPath {
			fileDependencies[filePath] = append(dependencies[:i], dependencies[i+1:]...)
			return nil
		}
	}
	return errors.New("dependency not found")
}

// RetrieveFileByContent retrieves a file by its content hash.
func (m *MFT) RetrieveFileByContent(contentHash, searchDirectory string) (string, error) {
	files, err := ioutil.ReadDir(searchDirectory)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(searchDirectory, file.Name())
			hash, err := m.CalculateChecksum(filePath)
			if err != nil {
				return "", err
			}
			if hash == contentHash {
				return filePath, nil
			}
		}
	}

	return "", os.ErrNotExist
}

// SanitizationRule represents a rule for data sanitization.
type SanitizationRule struct {
	Search  string
	Replace string
}

// SanitizeFileData sanitizes the content of a file.
func (m *MFT) SanitizeFileData(filePath string, sanitizationRules []SanitizationRule) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)
	for _, rule := range sanitizationRules {
		content = strings.ReplaceAll(content, rule.Search, rule.Replace)
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}

func (m *MFT) LogFileTransfer(sourcePath, destinationPath, status string) error {
	file, err := os.OpenFile("file_transfer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("Transfer: %s -> %s (%s)\n", sourcePath, destinationPath, status)

	return nil
}

// LoadConfiguration loads configuration from a file.
func (m *MFT) LoadConfiguration(configFilePath string) (map[string]interface{}, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfiguration saves configuration to a file.
func (m *MFT) SaveConfiguration(configFilePath string, config map[string]interface{}) error {
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&config); err != nil {
		return err
	}
	return nil
}

func (m *MFT) GetCrossPlatformPath(filePath string) (string, error) {
	return filepath.ToSlash(filePath), nil
}

// ChangeFileOwner changes the owner of a file.
func (m *MFT) ChangeFileOwner(filePath, owner string) error {
	usr, err := user.Lookup(owner)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return err
	}

	return os.Chown(filePath, uid, gid)
}

// GetFileOwner gets the owner of a file.
func (m *MFT) GetFileOwner(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	sys := fileInfo.Sys()
	stat, ok := sys.(*syscall.Stat_t)
	if !ok {
		return "", nil
	}

	uid := strconv.Itoa(int(stat.Uid))
	user, err := user.LookupId(uid)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

// MultiThreadedUpload uploads a file to a remote server using multiple threads.
func (m *MFT) MultiThreadedUpload(filePath, destination string, numThreads int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	partSize := fileSize / int64(numThreads)
	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			partStart := int64(i) * partSize
			partEnd := partStart + partSize
			if i == numThreads-1 {
				partEnd = fileSize
			}
			partFilePath := filePath + ".part" + strconv.Itoa(i)
			partFile, err := os.Create(partFilePath)
			if err != nil {
				fmt.Printf("Error occurred:%v", err)
				return
			}
			defer partFile.Close()
			_, err = file.Seek(partStart, io.SeekStart)
			if err != nil {
				fmt.Printf("Error occurred:%v", err)
				return
			}
			_, err = io.CopyN(partFile, file, partEnd-partStart)
			if err != nil && err != io.EOF {
				fmt.Printf("Error occurred:%v", err)
				return
			}
			err = m.UploadFile(destination, partFilePath, destination+"/"+filepath.Base(partFilePath))
			if err != nil {
				fmt.Printf("Error occurred:%v", err)
				return
			}
		}(i)
	}
	wg.Wait()

	return nil
}

// MultiThreadedDownload downloads a file from a remote server using multiple threads.
func (m *MFT) MultiThreadedDownload(filePath, destination string, numThreads int) error {
	// Implement multi-threaded download logic
	return nil
}

func (m *MFT) TransformFile(inputPath, outputPath string, transformer func(data []byte) []byte) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	transformedData := transformer(data)

	return os.WriteFile(outputPath, transformedData, 0644)
}

func (m *MFT) GetFileAccessTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

// SetFileAccessTime sets the last access time of a file.
func (m *MFT) SetFileAccessTime(filePath string, accessTime time.Time) error {
	return os.Chtimes(filePath, accessTime, accessTime)
}

// RenameFile renames a file.
func (m *MFT) RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// GetFileSize gets the size of a file.
func (m *MFT) GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// CreateDirectory creates a new directory.
func (m *MFT) CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

// DeleteDirectory deletes a directory.
func (m *MFT) DeleteDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}

// ListFiles lists all files in a directory.
func (m *MFT) ListFiles(directoryPath string) ([]string, error) {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList, nil
}

// ListDirectories lists all directories in a directory.
func (m *MFT) ListDirectories(directoryPath string) ([]string, error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	var dirList []string
	for _, file := range files {
		if file.IsDir() {
			dirList = append(dirList, file.Name())
		}
	}

	return dirList, nil
}

// ValidateFilePath checks if the file path exists and is accessible.
func (m *MFT) ValidateFilePath(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MFT) Version() string {
	return "1.0.0"
}

// FileMetadata represents metadata of a file.
type FileMetadata struct {
	Size        int64
	Permissions os.FileMode
	ModTime     time.Time
	IsDir       bool
}

// GetFileMetadata retrieves metadata for a file.
func (m *MFT) GetFileMetadata(filePath string) (FileMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileMetadata{}, err
	}

	return FileMetadata{
		Size:        fileInfo.Size(),
		Permissions: fileInfo.Mode(),
		ModTime:     fileInfo.ModTime(),
		IsDir:       fileInfo.IsDir(),
	}, nil
}

// SetFileMetadata sets metadata for a file.
func (m *MFT) SetFileMetadata(filePath string, metadata FileMetadata) error {
	err := os.Chmod(filePath, metadata.Permissions)
	if err != nil {
		return err
	}

	err = os.Chtimes(filePath, time.Now(), metadata.ModTime)
	if err != nil {
		return err
	}

	if !metadata.IsDir {
		return errors.New("cannot set metadata for non-directory")
	}

	return nil
}

// SaveFileVersion saves a version of a file.
func (m *MFT) SaveFileVersion(filePath string) (string, error) {
	versionPath := filePath + "." + strconv.FormatInt(time.Now().Unix(), 10)
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(versionPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return "", err
	}

	return versionPath, nil
}

// RevertToFileVersion reverts to a saved version of a file.
func (m *MFT) RevertToFileVersion(filePath, versionID string) error {
	versionPath := filePath + "." + versionID
	sourceFile, err := os.Open(versionPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// SyncDirectories synchronizes the content of two directories.
func (m *MFT) SyncDirectories(sourceDir, targetDir string) error {
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return m.CopyFile(path, targetPath)
	})
	return err
}

// copyFile copies a file from src to dst.
func (m *MFT) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return nil
}

// VerifyDataIntegrity verifies the integrity of a file using the specified hash type.
func (m *MFT) VerifyDataIntegrity(filePath, hashType, expectedHash string) (bool, error) {
	var hasher hash.Hash
	switch hashType {
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	default:
		return false, errors.New("unsupported hash type")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return false, err
	}

	calculatedHash := hex.EncodeToString(hasher.Sum(nil))
	return calculatedHash == expectedHash, nil
}

var customProtocolHandlers = make(map[string]func(request Request) (Response, error))

// AddCustomProtocolHandler adds a handler for a custom protocol.
func (m *MFT) AddCustomProtocolHandler(protocol string, handler func(request Request) (Response, error)) error {
	if protocol == "" {
		return errors.New("protocol name cannot be empty")
	}
	customProtocolHandlers[protocol] = handler
	return nil
}

// Request represents a custom protocol request.
type Request struct {
	Protocol string
	Data     []byte
}

// Response represents a custom protocol response.
type Response struct {
	Status string
	Data   []byte
}

// HandleCustomProtocolRequest handles a custom protocol request.
func (m *MFT) HandleCustomProtocolRequest(request Request) (Response, error) {
	handler, exists := customProtocolHandlers[request.Protocol]
	if !exists {
		return Response{}, errors.New("no handler for protocol")
	}
	return handler(request)
}

// ArchiveFiles creates a zip archive from a list of files.
func (m *MFT) ArchiveFiles(filePaths []string, archivePath string) error {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		err := m.AddFileToZip(zipWriter, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddFileToZip adds a file to a zip archive.
func (m *MFT) AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// UnarchiveFile extracts a zip archive to the specified destination.
func (m *MFT) UnarchiveFile(archivePath, destinationDir string) error {
	archiveFile, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	for _, file := range archiveFile.File {
		filePath := filepath.Join(destinationDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		err := m.ExtractFileFromZip(file, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExtractFileFromZip extracts a file from a zip archive.
func (m *MFT) ExtractFileFromZip(file *zip.File, filePath string) error {
	archiveFile, err := file.Open()
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	destFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, archiveFile)
	return err
}

// SetFilePermissions sets the permissions for a file.
func (m *MFT) SetFilePermissions(filePath string, permissions os.FileMode) error {
	return os.Chmod(filePath, permissions)
}

// GetFilePermissions gets the permissions of a file.
func (m *MFT) GetFilePermissions(filePath string) (os.FileMode, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Mode(), nil
}

var fileDependencies2 = make(map[string][]string)

// AddFileDependency2 adds a dependency for a file.
func (m *MFT) AddFileDependency2(filePath, dependencyPath string) error {
	fileDependencies2[filePath] = append(fileDependencies2[filePath], dependencyPath)
	return nil
}

// RemoveFileDependency2 removes a dependency for a file.
func (m *MFT) RemoveFileDependency2(filePath, dependencyPath string) error {
	dependencies, exists := fileDependencies2[filePath]
	if !exists {
		return errors.New("file not found")
	}

	for i, dep := range dependencies {
		if dep == dependencyPath {
			fileDependencies[filePath] = append(dependencies[:i], dependencies[i+1:]...)
			return nil
		}
	}

	return errors.New("dependency not found")
}

// TrackTransferProgress2 tracks the progress of a file transfer.
func (m *MFT) TrackTransferProgress2(filePath string, callback func(progress float64)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	totalSize := fileInfo.Size()
	buffer := make([]byte, 1024)
	var bytesRead int64

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		bytesRead += int64(n)
		progress := float64(bytesRead) / float64(totalSize) * 100
		callback(progress)
	}

	return nil
}
