// Copyright © 2017 shoarai

package dirname

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

const fileSuffix = "-%d"

// Rename renames a file or a directory and moves it to a directory.
func Rename(oldPath, newDir, newName string) (string, error) {
	if !isExist(oldPath) {
		return "", fmt.Errorf("Rename %s: no such file or directory", oldPath)
	}
	if !isExist(newDir) {
		return "", fmt.Errorf("Rename %s: no such file or directory", newDir)
	}

	_, oldFile := filepath.Split(oldPath)
	ext := filepath.Ext(oldFile)
	newPath := addSuffixIfSamePath(newDir, newName, ext)

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

func addSuffixIfSamePath(dir, file, ext string) string {
	path := filepath.Join(dir, file)
	if p := path + ext; !isExist(p) {
		return p
	}

	for i := 2; i < math.MaxInt16; i++ {
		suff := fmt.Sprintf(fileSuffix, i)
		if p := path + suff + ext; !isExist(p) {
			return p
		}
	}
	return ""
}

// RenameAll renames all files in root and moves these to a directory.
func RenameAll(root, newDir, newFileName string) error {
	if !isExist(newDir) {
		return fmt.Errorf("RenameAll %s: no such file or directory", newDir)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		Rename(path, newDir, newFileName)
		return nil
	})
}

// RenamePattern renames all files matching pattern in root
// and moves these to a directory.
func RenamePattern(root, newDir, newFileName string, pattern string) error {
	if !isExist(newDir) {
		return fmt.Errorf("RenamePattern %s: no such file or directory", newDir)
	}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !reg.MatchString(info.Name()) {
			return nil
		}
		Rename(path, newDir, newFileName)
		return nil
	})
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
