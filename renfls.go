// Copyright © 2017 shoarai

// Package renfls provides interfaces to rename files in directory.
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
	fileSuffix  = "-%d"
	tempDirName = "fails"
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

	if e := os.Rename(oldPath, newPath); e != nil {
		return "", e
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

	return filepath.Walk(root, walkRenameFunc(newDir, newFileName,
		func(info os.FileInfo) bool {
			return !info.IsDir()
		}))
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

	reg, e := regexp.Compile(pattern)
	if e != nil {
		return e
	}

	return filepath.Walk(root, walkRenameFunc(newDir, newFileName,
		func(info os.FileInfo) bool {
			return !info.IsDir() && reg.MatchString(info.Name())
		}))
}

func walkRenameFunc(newDir, newFileName string,
	condition func(info os.FileInfo) bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !condition(info) {
			return nil
		}
		if _, err := Rename(path, newDir, newFileName); err != nil {
			return err
		}
		return nil
	}
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
	tempDir, e := moveDirs(root, tempDirName)
	if e != nil {
		return e
	}
	if e := renameToDirName(tempDir, root); e != nil {
		return e
	}
	if e := os.RemoveAll(tempDir); e != nil {
		return e
	}
	return nil
}

// ToDirNamesPattern renames all files matching pattern in root
// by the directories name in root and moves these to a directory.
func ToDirNamesPattern(root, pattern string) error {
	tempDir, e := moveDirs(root, tempDirName)
	if e != nil {
		return e
	}
	if e := renameToDirNamePattern(tempDir, root, pattern); e != nil {
		return e
	}
	if e := os.RemoveAll(tempDir); e != nil {
		return e
	}
	return nil
}

func renameToDirName(root, newDir string) error {
	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir.Name())
		if e := ToRootDirName(path, newDir); e != nil {
			return e
		}
	}
	return nil
}

func renameToDirNamePattern(root, newDir, pattern string) error {
	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir.Name())
		if e := ToRootDirNamePattern(path, newDir, pattern); e != nil {
			return e
		}
	}
	return nil
}

func moveDirs(root, newDir string) (string, error) {
	if isNotExist(root) {
		return "", fmt.Errorf("ToDirNames %s: no such file or directory", root)
	}

	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return "", e
	}

	tempDir := filepath.Join(root, newDir)
	if e := os.Mkdir(tempDir, os.ModePerm); e != nil {
		return "", fmt.Errorf("ToDirNames: Temporary directory can't be created. %s", e)
	}

	for _, dir := range dirs {
		// if !dir.IsDir() {
		// 	continue
		// }
		path := filepath.Join(root, dir.Name())
		dirInTempDir := filepath.Join(tempDir + "/" + dir.Name())
		if e := os.Rename(path, dirInTempDir); e != nil {
			continue
		}
	}

	return tempDir, nil
}

func isNotExist(path string) bool {
	_, e := os.Stat(path)
	return e != nil
}
