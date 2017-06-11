// Copyright © 2017 shoarai

// The renfls renames files in a directory.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shoarai/renfls"
)

const toDir = "toSubDirsName"

// Flag
var ext string
var ignore bool

func main() {
	// DEBUG:
	// createTestDir()

	flag.BoolVar(&ignore, "ignore", false, "bool flag")
	flag.StringVar(&ext, "ext", "", "extensions splited by \",\"")
	flag.Parse()

	exts := strings.Split(ext, ",")
	var err error
	if !ignore {
		// err = renfls.ToSubDirsNameExt(toDir, exts)
	} else {
		err = renfls.ToSubDirsNameIgnoreExt(toDir, exts)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func createTestDir() {
	os.RemoveAll(toDir)
	exec.Command("cp", "-r", "testdata", toDir).Run()
}
