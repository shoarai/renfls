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

const separator = ","

// Flag parameter
var dest string
var ext string
var reg string
var ignore bool

func main() {
	flag.StringVar(&dest, "dest", "", "Destination to which renamed files are moved")
	flag.StringVar(&ext, "ext", "",
		fmt.Sprintf("Extension list separated by %q", separator))
	flag.StringVar(&reg, "reg", "", "Regex")
	flag.BoolVar(&ignore, "ignore", false,
		"Flag whether files matching pattern are renamed or ignored.")
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		fmt.Println("Input root directory name as command argument")
		return
	}
	if dest == "" {
		dest = root
	}
	var exts []string
	if ext != "" {
		exts = strings.Split(ext, separator)
	}

	// DEBUG: Copy test files
	// createTestDir()

	fmt.Println(root)
	fmt.Println(dest)

	condition := renfls.Condition{Exts: exts, Reg: reg, Ignore: ignore}
	if e := renfls.WalkToRootSubDirName(root, dest, condition); e != nil {
		fmt.Println(e)
	}
}

func createTestDir() {
	dir := "rootForMain"
	os.RemoveAll(dir)
	exec.Command("cp", "-r", "testdata", dir).Run()
}
