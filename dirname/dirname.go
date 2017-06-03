// Copyright © 2017 shoarai

package dirname

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RenameAndMoveFile renames a file and moves it to directory.
func RenameAndMoveFile(
	oldDir, oldFileName, newDir, newFileName string) (string, error) {
	oldFilePath := filepath.Join(oldDir, oldFileName)
	if _, err := os.Stat(oldFilePath); err != nil {
		return "", fmt.Errorf("The old file is not existing %s", oldFilePath)
	}

	pos := strings.LastIndex(oldFileName, ".")
	var extension string
	if pos >= 0 {
		extension = oldFileName[pos:]
	}

	var newFilePath string
	for i := 0; i < math.MaxInt16; i++ {
		var suffix string
		if i != 0 {
			suffix = "-" + strconv.Itoa(i)
		}

		newFilePath = filepath.Join(newDir, newFileName+suffix+extension)
		fmt.Println(newFilePath)
		if _, err := os.Stat(newFilePath); err != nil {
			break
		}
	}

	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		return "", err
	}
	return newFilePath, nil
}

// RenameAndMoveFilesInDir all renames files in the directory and moves it.
func RenameAndMoveFilesInDir(dir, newDir, newFileName string) error {
	for _, entry := range GetFileInfoInDir(dir) {
		_, err := RenameAndMoveFile(dir, entry.Name(), newDir, newFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFileInfoInDir gets information for files in a directory.
func GetFileInfoInDir(dir string) []os.FileInfo {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) //- 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
