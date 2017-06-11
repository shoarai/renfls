// Copyright © 2017 shoarai

// renfls renames all files or files matching patterns in directories.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shoarai/renfls"
)

const (
	separator = ","
	toDir     = "toSubDirsName"
)

// Flag parameter
var ext string
var ignore bool

func main() {
	flag.StringVar(&ext, "ext", "",
		fmt.Sprintf("extensions splited by %q", separator))
	flag.BoolVar(&ignore, "ignore", false, "bool flag")
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		fmt.Println("Input root directory name as command argument")
		return
	}
	var exts []string
	if ext != "" {
		exts = strings.Split(ext, separator)
	}

	// DEBUG:
	// createTestDir()

	if e := toSubDirsName(root, exts, ignore); e != nil {
		fmt.Println(e)
	}
}

func createTestDir() {
	os.RemoveAll(toDir)
	exec.Command("cp", "-r", "testdata", toDir).Run()
}

func toSubDirsName(root string, exts []string, ignore bool) error {
	if len(ext) == 0 {
		return renfls.ToSubDirsName(root)
	}
	if ignore {
		return renfls.ToSubDirsNameIgnoreExt(root, exts)
	}
	return renfls.ToSubDirsNameExt(root, exts)
}
