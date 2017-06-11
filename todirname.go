// Copyright © 2017 shoarai

// Package renfls provides interfaces to rename files in directory.
package renfls

import "path/filepath"

// ToDirName renames all files in root
// by the root directory name and moves these to a directory.
func ToDirName(root, newDir string) error {
	_, name := filepath.Split(root)
	return RenameAll(root, newDir, name)
}

// ToDirNamePattern renames all files matching pattern in root
// by the root directory name and moves these to a directory.
func ToDirNamePattern(root, newDir, pattern string) error {
	_, name := filepath.Split(root)
	return RenamePattern(root, newDir, name, pattern)
}

// ToDirNameExt renames all files matching extensions in root
// by the root directory name and moves these to a directory.
func ToDirNameExt(root, newDir string, exts []string) error {
	_, name := filepath.Split(root)
	return RenameExt(root, newDir, name, exts)
}

// ToDirNameIgnoreExt renames all files not matching extensions in root
// by the root directory name and moves these to a directory.
func ToDirNameIgnoreExt(root, newDir string, exts []string) error {
	_, name := filepath.Split(root)
	return RenameIgnoreExt(root, newDir, name, exts)
}

// ToDirNameCondition renames all files matching condition in root
// by the root directory name and moves these to a directory.
func ToDirNameCondition(root, newDir string, condition Condition) error {
	_, name := filepath.Split(root)
	return RenameCondition(root, newDir, name, condition)
}
