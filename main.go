// Copyright © 2017 shoarai

// The ToDirName re-name the files in a directory to the directory.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shoarai/toDirName/dirname"
)

func main() {
	walkDirectory("./todirname")
}

func walkDirectory(dir string) {
	for _, entry := range dirname.GetFileInfoInDir(dir) {
		if !entry.IsDir() {
			continue
		}
		var folderName = entry.Name()
		subDir := filepath.Join(dir, folderName)
		dirNameToFileName(subDir, dir, folderName)
	}
}

func dirNameToFileName(dir, toDir, name string) {
	// TODO: Ignore "." file.
	// TODO: Make log file.

	err := dirname.RenameAndMoveFilesInDir(dir, toDir, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RenameAndMoveFilesInDir(): %v\n", err)
		return
	}

	err = os.RemoveAll(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RemoveAll(): %v\n", err)
		return
	}
}
