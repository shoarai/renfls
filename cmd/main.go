// Copyright © 2017 shoarai

// The ToDirName re-name the files in a directory to the directory.
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/shoarai/renfls"
)

const toDir = "toDirNameDirectories"

func main() {
	createTestDir()
	if err := renfls.ToDirNames(toDir); err != nil {
		fmt.Println(err)
	}
}

func createTestDir() {
	os.RemoveAll(toDir)
	exec.Command("cp", "-r", "testdata", toDir).Run()
}
