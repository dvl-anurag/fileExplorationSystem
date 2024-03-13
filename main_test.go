package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestListFiles_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create some files and directories inside the temporary directory
	createTempFiles(t, tempDir)

	// Create a ListFiles instance
	listFiles := ListFiles{Path: tempDir}

	// Execute the ListFiles operation
	err := listFiles.Execute()
	if err != nil {
		t.Errorf("ListFiles.Execute() returned an error: %v", err)
	}
}

func TestSearchFiles_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create some files and directories inside the temporary directory
	createTempFiles(t, tempDir)

	// Create a SearchFiles instance
	searchFiles := SearchFiles{Path: tempDir, Query: "test"}

	// Execute the SearchFiles operation
	err := searchFiles.Execute()
	if err != nil {
		t.Errorf("SearchFiles.Execute() returned an error: %v", err)
	}
}

func TestCopyFile_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	createTempFile(t, sourceFile)

	// Define the destination file path
	destFile := filepath.Join(tempDir, "destination.txt")

	// Create a CopyFile instance
	copyFile := CopyFile{Source: sourceFile, Destination: destFile}

	// Execute the CopyFile operation
	err := copyFile.Execute()
	if err != nil {
		t.Errorf("CopyFile.Execute() returned an error: %v", err)
	}

	// Check if the destination file exists
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Errorf("CopyFile.Execute() failed to copy the file to the destination")
	}
}

func TestMoveFile_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	createTempFile(t, sourceFile)

	// Define the destination file path
	destFile := filepath.Join(tempDir, "destination.txt")

	// Create a MoveFile instance
	moveFile := MoveFile{Source: sourceFile, Destination: destFile}

	// Execute the MoveFile operation
	err := moveFile.Execute()
	if err != nil {
		t.Errorf("MoveFile.Execute() returned an error: %v", err)
	}

	// Check if the source file exists
	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Errorf("MoveFile.Execute() failed to move the file from the source")
	}

	// Check if the destination file exists
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Errorf("MoveFile.Execute() failed to move the file to the destination")
	}
}

func TestDeleteFile_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a file
	file := filepath.Join(tempDir, "test.txt")
	createTempFile(t, file)

	// Create a DeleteFile instance
	deleteFile := DeleteFile{Path: file}

	// Execute the DeleteFile operation
	err := deleteFile.Execute()
	if err != nil {
		t.Errorf("DeleteFile.Execute() returned an error: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Errorf("DeleteFile.Execute() failed to delete the file")
	}
}

// Helper functions to create temporary files and directories for testing
func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "file_operator_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return tempDir
}

func createTempFile(t *testing.T, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer file.Close()
}

func createTempFiles(t *testing.T, dirPath string) {
	for i := 1; i <= 3; i++ {
		filePath := filepath.Join(dirPath, fmt.Sprintf("file%d.txt", i))
		createTempFile(t, filePath)
	}
}
