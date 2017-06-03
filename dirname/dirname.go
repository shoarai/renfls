// Copyright © 2017 shoarai

package dirname

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

// RenameAndMoveFile renames a file and moves it to directory.
func RenameAndMoveFile(oldPath, newDir, newFileName string) (string, error) {
	if _, err := os.Stat(oldPath); err != nil {
		return "", fmt.Errorf("oldPath %q is not existing", oldPath)
	}

	_, oldFileName := filepath.Split(oldPath)
	extension := filepath.Ext(oldFileName)

	var newPath string
	for i := 0; i < math.MaxInt16; i++ {
		var suffix string
		if i != 0 {
			suffix = "-" + strconv.Itoa(i)
		}
		newPath = filepath.Join(newDir, newFileName+suffix+extension)
		if _, err := os.Stat(newPath); err != nil {
			break
		}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

// RenameAndMoveFilesInDir all renames files in the directory and moves it.
func RenameAndMoveFilesInDir(dir, newDir, newFileName string) error {
	for _, entry := range GetFileInfoInDir(dir) {
		_, err := RenameAndMoveFile(
			filepath.Join(dir, entry.Name()), newDir, newFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// RenameAndMoveFilesInDirRecursive all renames files in the directory and moves it.
func RenameAndMoveFilesInDirRecursive(dir, newDir, newFileName string) error {
	for _, entry := range GetFileInfoInDir(dir) {
		if entry.IsDir() {
			subDir := filepath.Join(dir, entry.Name())
			RenameAndMoveFilesInDirRecursive(subDir, newDir, newFileName)
			continue
		}
		_, err := RenameAndMoveFile(
			filepath.Join(dir, entry.Name()), newDir, newFileName)
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
