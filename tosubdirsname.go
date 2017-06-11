// Copyright © 2017 shoarai

// Package renfls provides interfaces to rename files in directory.
package renfls

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	tempDirName   = "fail"
	ignoreDirName = "ignore"
)

// ToSubDirsName renames all files in root
// by the directories name in root and moves these to a directory.
func ToSubDirsName(root string) error {
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

func renameToDirName(root, newDir string) error {
	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir.Name())
		if e := ToDirName(path, newDir); e != nil {
			return e
		}
	}
	return nil
}

// ToSubDirsNamePattern renames all files matching pattern in root
// by the directories name in root and moves these to a directory.
func ToSubDirsNamePattern(root, pattern string) error {
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

func renameToDirNamePattern(root, newDir, pattern string) error {
	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir.Name())
		if e := ToDirNamePattern(path, newDir, pattern); e != nil {
			return e
		}
	}
	return nil
}

// ToSubDirsNameIgnoreExt renames all files not matching extensions in root
// by the directories name in root and moves these to a directory.
func ToSubDirsNameIgnoreExt(root string, exts []string) error {
	tempDir, e := moveDirs(root, tempDirName)
	if e != nil {
		return e
	}
	if e := renameToDirNameIgnoreExt(tempDir, root, exts); e != nil {
		return e
	}
	if e := os.RemoveAll(tempDir); e != nil {
		return e
	}
	return nil
}

func renameToDirNameIgnoreExt(root, newDir string, exts []string) error {
	dirs, e := ioutil.ReadDir(root)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir.Name())
		if e := ToDirNameIgnoreExt(path, newDir, exts); e != nil {
			return e
		}
	}
	return nil
}

func moveDirs(root, newDir string) (string, error) {
	if isNotExist(root) {
		return "", errorNotExist("ToDirNames", root)
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
		if !dir.IsDir() {
			continue
		}
		path := filepath.Join(root, dir.Name())
		dirInTempDir := filepath.Join(tempDir + "/" + dir.Name())
		if e := os.Rename(path, dirInTempDir); e != nil {
			continue
		}
	}

	return tempDir, nil
}
