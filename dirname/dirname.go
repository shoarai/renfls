// Copyright © 2017 shoarai

package dirname

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// RenameAndMoveFile renames a file and moves it to directory.
func RenameAndMoveFile(
	oldDir, oldFileName, newDir, newFileName string) (string, error) {
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

		newFilePath = newDir + "/" + newFileName + suffix + extension
		fmt.Println(newFilePath)
		if _, err := os.Stat(newFilePath); err != nil {
			break
		}
	}

	err := os.Rename(oldDir+"/"+oldFileName, newFilePath)
	if err != nil {
		return "", err
	}
	return newFilePath, nil
}

func joinDir(dirs ...string) string {
	var fullDir string
	separator := "/"
	for _, dir := range dirs {
		fullDir += dir
		if !strings.HasSuffix(dir, separator) {
			dir += "/"
		}
	}
	return fullDir
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
