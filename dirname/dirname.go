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
	if !isFileExisting(oldPath) {
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
		if !isFileExisting(newPath) {
			break
		}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

// RenameAndMoveFileAll renames files in root and moves these to a directory.
func RenameAndMoveFileAll(root, newDir, newFileName string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		RenameAndMoveFile(path, newDir, newFileName)
		return nil
	})
}

func isFileExisting(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
