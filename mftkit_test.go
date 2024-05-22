package main

import (
	"fmt"
	"github.com/madhu72/mftkit/mft"
	"testing"
)

func TestMFTUtils(t *testing.T) {
	utils := mft.MFT{}

	// Test CopyFile
	err := utils.CopyFile("sample.txt", "test.txt")
	if err != nil {
		t.Errorf("Error copying file: %v", err)
	}

	// Test ReadFile
	content, err := utils.ReadFile("test.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	} else {
		expectedContent := "This is for testing"
		if content != expectedContent {
			t.Errorf("Expected file content: %v, got: %v", expectedContent, content)
		}
	}

	creationTime, err := utils.GetFileCreationTime("sample.txt")
	if err != nil {
		fmt.Println("Error getting file creation time:", err)
	} else {
		fmt.Println("File creation time:", creationTime)
	}

	// Add more tests for other operations...
}
