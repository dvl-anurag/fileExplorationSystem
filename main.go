package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	infoColor  = color.New(color.FgHiCyan).SprintFunc()
	errorColor = color.New(color.FgHiRed).SprintFunc()
)

const SOURCE_PATH = "Enter the source path:"
const RECURSIVE_BOOL = "Enter true/false for recursive:"

type FileOperator interface {
	Execute() error
}

type ListFiles struct {
	Path      string
	Recursive bool
}

type SearchFiles struct {
	Path      string
	Query     string
	Recursive bool
}

type CopyFile struct {
	Source      string
	Destination string
	Recursive   bool
}

type MoveFile struct {
	Source      string
	Destination string
	Recursive   bool
}

type DeleteFile struct {
	Path      string
	Recursive bool
}

func main() {
	for {
		fmt.Println("Select an operation:")
		fmt.Println("1. List files")
		fmt.Println("2. Search files")
		fmt.Println("3. Copy file")
		fmt.Println("4. Move file")
		fmt.Println("5. Delete file")
		fmt.Println("6. Exit")

		selectedOp := readUserInput("Enter the operation number: ")
		opNum, err := strconv.Atoi(selectedOp)
		if err != nil || opNum < 1 || opNum > 6 {
			fmt.Println("Invalid operation number. Please enter a number between 1 and 6.")
			continue
		}

		switch opNum {
		case 1:
			listFilesOperation()
		case 2:
			searchFilesOperation()
		case 3:
			copyFileOperation()
		case 4:
			moveFileOperation()
		case 5:
			deleteFileOperation()
		case 6:
			fmt.Println("Exiting program...")
			return
		}
	}
}

func listFilesOperation() {
	sourcePath := readUserInput(SOURCE_PATH)
	recursive := readBoolInput(RECURSIVE_BOOL)

	listFiles := ListFiles{Path: sourcePath, Recursive: recursive}
	if err := listFiles.Execute(); err != nil {
		logError(err)
	}
}

func searchFilesOperation() {
	sourcePath := readUserInput(SOURCE_PATH)
	query := readUserInput("Enter the search query: ")
	recursive := readBoolInput(RECURSIVE_BOOL)

	searchFiles := SearchFiles{Path: sourcePath, Query: query, Recursive: recursive}
	if err := searchFiles.Execute(); err != nil {
		logError(err)
	}
}

func copyFileOperation() {
	sourcePath := readUserInput(SOURCE_PATH)
	destinationPath := readUserInput("Enter the destination path: ")
	recursive := readBoolInput(RECURSIVE_BOOL)

	copyFile := CopyFile{Source: sourcePath, Destination: destinationPath, Recursive: recursive}
	if err := copyFile.Execute(); err != nil {
		logError(err)
	}
}

func moveFileOperation() {
	sourcePath := readUserInput(SOURCE_PATH)
	destinationPath := readUserInput("Enter the destination path: ")
	recursive := readBoolInput(RECURSIVE_BOOL)

	moveFile := MoveFile{Source: sourcePath, Destination: destinationPath, Recursive: recursive}
	if err := moveFile.Execute(); err != nil {
		logError(err)
	}
}

func deleteFileOperation() {
	sourcePath := readUserInput(SOURCE_PATH)
	recursive := readBoolInput(RECURSIVE_BOOL)

	deleteFile := DeleteFile{Path: sourcePath, Recursive: recursive}
	if err := deleteFile.Execute(); err != nil {
		logError(err)
	}
}

func readUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readBoolInput(prompt string) bool {
	input := readUserInput(prompt)
	b, err := strconv.ParseBool(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter true or false.")
		return readBoolInput(prompt)
	}
	return b
}

func (lf ListFiles) Execute() error {
	fmt.Printf("Listing files and directories in: %s\n", infoColor(lf.Path))

	err := filepath.Walk(lf.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logError(err)
			return nil
		}
		fmt.Println(path)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (sf SearchFiles) Execute() error {
	fmt.Printf("Searching for files matching '%s' in: %s\n", infoColor(sf.Query), infoColor(sf.Path))

	err := filepath.Walk(sf.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logError(err)
			return nil
		}
		if !info.IsDir() && (sf.Query == "" || containsNameOrExtension(info, sf.Query)) {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func containsNameOrExtension(info os.FileInfo, query string) bool {
	return containsIgnoreCase(info.Name(), query) || containsIgnoreCase(filepath.Ext(info.Name()), query)
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (cf CopyFile) Execute() error {
	fmt.Printf("Copying file from %s to %s\n", infoColor(cf.Source), infoColor(cf.Destination))

	err := copyFile(cf.Source, cf.Destination)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (mf MoveFile) Execute() error {
	fmt.Printf("Moving file from %s to %s\n", infoColor(mf.Source), infoColor(mf.Destination))

	err := moveFile(mf.Source, mf.Destination)
	if err != nil {
		return err
	}

	return nil
}

func moveFile(src, dest string) error {
	return os.Rename(src, dest)
}

func (df DeleteFile) Execute() error {
	fmt.Printf("Deleting file or directory: %s\n", infoColor(df.Path))

	err := deleteFile(df.Path, df.Recursive)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(path string, recursive bool) error {
	if recursive {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func logError(err error) {
	fmt.Println(errorColor(err.Error()))
}
