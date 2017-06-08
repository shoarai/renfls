// Copyright © 2017 shoarai

package renfls

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

const (
	fileSuffix = "-%d"
	tempDir    = "fails"
)

// Rename renames a file or a directory and moves it to a directory.
func Rename(oldPath, newDir, newName string) (string, error) {
	if isNotExist(oldPath) {
		return "", fmt.Errorf("Rename %s: no such file or directory", oldPath)
	}
	if isNotExist(newDir) {
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
	if p := path + ext; isNotExist(p) {
		return p
	}

	for i := 2; i < math.MaxInt16; i++ {
		suff := fmt.Sprintf(fileSuffix, i)
		if p := path + suff + ext; isNotExist(p) {
			return p
		}
	}
	return ""
}

// RenameAll renames all files in root and moves these to a directory.
func RenameAll(root, newDir, newFileName string) error {
	if isNotExist(root) {
		return fmt.Errorf("RenameAll %s: no such file or directory", root)
	}
	if isNotExist(newDir) {
		return fmt.Errorf("RenameAll %s: no such file or directory", newDir)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if _, err := Rename(path, newDir, newFileName); err != nil {
			return err
		}
		return nil
	})
}

// RenamePattern renames all files matching pattern in root
// and moves these to a directory.
func RenamePattern(root, newDir, newFileName, pattern string) error {
	if isNotExist(root) {
		return fmt.Errorf("RenamePattern %s: no such file or directory", root)
	}
	if isNotExist(newDir) {
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
		if _, err := Rename(path, newDir, newFileName); err != nil {
			return err
		}
		return nil
	})
}

// ToRootDirName renames all files in root
// by the root directory name and moves these to a directory.
func ToRootDirName(root, newDir string) error {
	_, file := filepath.Split(root)
	return RenameAll(root, newDir, file)
}

// ToRootDirNamePattern renames all files matching pattern in root
// by the root directory name and moves these to a directory.
func ToRootDirNamePattern(root, newDir, pattern string) error {
	_, file := filepath.Split(root)
	return RenamePattern(root, newDir, file, pattern)
}

// ToDirNames renames all files in root
// by the directories name in root and moves these to a directory.
func ToDirNames(root string) error {
	if isNotExist(root) {
		return fmt.Errorf("ToDirNames %s: no such file or directory", root)
	}

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	tempDirInRoot := filepath.Join(root, tempDir)
	if err := os.Mkdir(tempDirInRoot, os.ModePerm); err != nil {
		return fmt.Errorf("ToDirNames: Temporary directory can't create. %s", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		path := filepath.Join(root, dir.Name())
		dirInTempDIr := filepath.Join(tempDirInRoot + "/" + dir.Name())
		if err := os.Rename(path, dirInTempDIr); err != nil {
			continue
		}
	}

	dirs, err = ioutil.ReadDir(tempDirInRoot)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		path := filepath.Join(tempDirInRoot, dir.Name())
		if err := ToRootDirName(path, root); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(tempDirInRoot); err != nil {
		return err
	}
	return nil
}

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}
