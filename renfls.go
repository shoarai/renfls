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
func Rename(oldPath, dest, newName string) (string, error) {
	return rename(oldPath, dest, newName)
}

func rename(oldPath, dest, newName string) (string, error) {
	if isNotExist(oldPath) {
		return "", errorNotExist("Rename", oldPath)
	}
	if isNotExist(dest) {
		return "", errorNotExist("Rename", dest)
	}

	_, oldFile := filepath.Split(oldPath)
	ext := filepath.Ext(oldFile)
	newPath, e := addSuffixIfExist(dest, newName, ext)
	if e != nil {
		return "", e
	}

	if e := os.Rename(oldPath, newPath); e != nil {
		return "", e
	}
	return newPath, nil
}

func addSuffixIfExist(dir, file, ext string) (string, error) {
	path := filepath.Join(dir, file)
	if p := path + ext; isNotExist(p) {
		return p, nil
	}

	for i := 2; i < math.MaxInt16; i++ {
		suff := fmt.Sprintf(fileSuffix, i)
		if p := path + suff + ext; isNotExist(p) {
			return p, nil
		}
	}
	return "", fmt.Errorf("Add file suffix failed")
}

// WalkRenameAll renames all files in a root directory
// and moves them to a destination directory.
func WalkRenameAll(root, dest, newFileName string) error {
	return walkRename(root, dest, newFileName, nil)
}

// Condition is condition to rename files.
type Condition struct {
	Exts   []string
	Reg    string
	Ignore bool
}

// WalkRename renames files that match a condition in a root directory
// and moves them to a destination directory.
func WalkRename(root, dest, newFileName string, condition Condition) error {
	var reg *regexp.Regexp

	if condition.Reg != "" {
		var e error
		reg, e = regexp.Compile(condition.Reg)
		if e != nil {
			return e
		}
	}

	isMatch := func(info os.FileInfo) bool {
		if reg == nil && (condition.Exts == nil || len(condition.Exts) == 0) {
			return true
		}
		if reg != nil && reg.MatchString(info.Name()) {
			return true
		}
		if condition.Exts != nil && len(condition.Exts) > 0 && hasExt(info.Name(), condition.Exts) {
			return true
		}
		return false
	}

	return walkRename(root, dest, newFileName, func(info os.FileInfo) bool {
		if condition.Ignore {
			return !isMatch(info)
		} else {
			return isMatch(info)
		}
	})
}

// NeedRename returns whether the file needs to be rename.
type NeedRename func(info os.FileInfo) bool

func walkRename(root, dest, newFileName string, needRename NeedRename) error {
	if isNotExist(root) {
		return errorNotExist("RenameAll", root)
	}
	if isNotExist(dest) {
		return errorNotExist("RenameAll", dest)
	}
	return filepath.Walk(root, walkRenameFunc(dest, newFileName, needRename))
}

func walkRenameFunc(dest, newFileName string, needRename NeedRename) filepath.WalkFunc {
	if needRename == nil {
		needRename = func(nfo os.FileInfo) bool { return true }
	}

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !needRename(info) {
			return nil
		}
		if _, err := Rename(path, dest, newFileName); err != nil {
			return err
		}
		return nil
	}
}

// RenamePattern renames all files matching pattern in root
// and moves these to a directory.
func RenamePattern(root, dest, newFileName, pattern string) error {
	reg, e := regexp.Compile(pattern)
	if e != nil {
		return e
	}
	return walkRename(root, dest, newFileName, func(info os.FileInfo) bool {
		return reg.MatchString(info.Name())
	})
}

// RenameExt renames all files matching extensions in root
// and moves these to a directory.
func RenameExt(root, dest, newFileName string, exts []string) error {
	return walkRename(root, dest, newFileName, func(info os.FileInfo) bool {
		return hasExt(info.Name(), exts)
	})
}

// RenameIgnoreExt renames all files not matching extension ins root
// and moves these to a directory.
func RenameIgnoreExt(root, dest, newFileName string, exts []string) error {
	return walkRename(root, dest, newFileName, func(info os.FileInfo) bool {
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
func RenameCondition(root, dest, newFileName string, needRename NeedRename) error {
	return walkRename(root, dest, newFileName, needRename)
}

func errorNotExist(funcName, path string) error {
	return fmt.Errorf("%s %s: no such file or directory", funcName, path)
}

func isNotExist(path string) bool {
	_, e := os.Stat(path)
	return e != nil
}
