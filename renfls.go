// Copyright © 2017 shoarai

// Package renfls provides interfaces to rename files in directory.
package renfls

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const fileSuffix = "-%d"

// Rename renames a file or a directory and moves it to a directory.
func Rename(oldPath, newDir, newName string) (string, error) {
	return rename(oldPath, newDir, newName)
}

func rename(oldPath, newDir, newName string) (string, error) {
	if isNotExist(oldPath) {
		return "", errorNotExist("Rename", oldPath)
	}
	if isNotExist(newDir) {
		return "", errorNotExist("Rename", newDir)
	}

	_, oldFile := filepath.Split(oldPath)
	ext := filepath.Ext(oldFile)
	newPath := addSuffixIfExist(newDir, newName, ext)

	if e := os.Rename(oldPath, newPath); e != nil {
		return "", e
	}
	return newPath, nil
}

func addSuffixIfExist(dir, file, ext string) string {
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
	return renameAll(root, newDir, newFileName, nil)
}

func renameAll(root, newDir, newFileName string, condition Condition) error {
	if isNotExist(root) {
		return errorNotExist("RenameAll", root)
	}
	if isNotExist(newDir) {
		return errorNotExist("RenameAll", newDir)
	}
	return filepath.Walk(root, walkRenameFunc(newDir, newFileName, condition))
}

// Condition is condition whether rename or not.
type Condition func(info os.FileInfo) bool

// RenamePattern renames all files matching pattern in root
// and moves these to a directory.
func RenamePattern(root, newDir, newFileName, pattern string) error {
	reg, e := regexp.Compile(pattern)
	if e != nil {
		return e
	}
	return renameAll(root, newDir, newFileName, func(info os.FileInfo) bool {
		return reg.MatchString(info.Name())
	})
}

// RenameExt renames all files matching extensions in root
// and moves these to a directory.
func RenameExt(root, newDir, newFileName string, exts []string) error {
	return renameAll(root, newDir, newFileName, func(info os.FileInfo) bool {
		return hasExt(info.Name(), exts)
	})
}

// RenameIgnoreExt renames all files not matching extension ins root
// and moves these to a directory.
func RenameIgnoreExt(root, newDir, newFileName string, exts []string) error {
	return renameAll(root, newDir, newFileName, func(info os.FileInfo) bool {
		return !hasExt(info.Name(), exts)
	})
}

func hasExt(file string, exts []string) bool {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	return contains(exts, ext)
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

// RenameCondition renames all files matching condition in root
// and moves these to a directory.
func RenameCondition(root, newDir, newFileName string, condition Condition) error {
	return renameAll(root, newDir, newFileName, condition)
}

func errorNotExist(funcName, path string) error {
	return fmt.Errorf("%s %s: no such file or directory", funcName, path)
}

func walkRenameFunc(newDir, newFileName string,
	condition Condition) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if condition != nil && !condition(info) {
			return nil
		}
		if _, err := Rename(path, newDir, newFileName); err != nil {
			return err
		}
		return nil
	}
}

func isNotExist(path string) bool {
	_, e := os.Stat(path)
	return e != nil
}
