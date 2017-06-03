// Copyright © 2017 shoarai

// The ToDirName re-name the files in a directory to the directory.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shoarai/toDirName/dirname"
)

const toDir = "toDirNameDirectories"

func main() {
	// DEBUG:
	createTestDir()

	if entry, err := os.Stat(toDir); err != nil || !entry.IsDir() {
		fmt.Printf("Make a directory and place folders in it: %q\n", toDir)
		return
	}

	workingDir := "." + toDir
	if _, err := os.Stat(workingDir); err == nil {
		fmt.Printf("The directory is already existing: %q\n", workingDir)
		return
	}
	os.Rename(toDir, workingDir)
	os.Mkdir(toDir, 0777)

	renameToDirName(workingDir, toDir)

	os.RemoveAll(workingDir)
}

func createTestDir() {
	os.RemoveAll(toDir)
	exec.Command("cp", "-r", "dirname/testdata", toDir).Run()
}

func renameToDirName(dir, newDir string) error {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
		// if !entry.IsDir() {
		// 	continue
		// }
		dirName := entry.Name()
		path := filepath.Join(dir, dirName)
		dirname.RenameAndMoveFileAll(path, newDir, dirName)
	}
	return nil
}
