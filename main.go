// Copyright © 2017 shoarai

// The ToDirName re-name the files in a directory to the directory.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shoarai/toDirName/dirname"
)

const toDir = "toDirNameDirectories"

func main() {
	createTestDir()

	entry, err := os.Stat(toDir)
	if err != nil {
		os.Mkdir(toDir, 0777)
		fmt.Printf("Place directories in the %q folder.\n", toDir)
		return
	}
	if !entry.IsDir() {
		fmt.Printf("The %q isn't folder.\n", toDir)
		return
	}

	walkDirectory(toDir)
}

func createTestDir() {
	os.RemoveAll(toDir)
	exec.Command("cp", "-r", "dirname/testdata", toDir).Run()
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

	err := dirname.RenameAndMoveFilesInDirRecursive(dir, toDir, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RenameAndMoveFilesInDir(): %v\n", err)
		return
	}

	// err = os.RemoveAll(dir)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "RemoveAll(): %v\n", err)
	// 	return
	// }
}
