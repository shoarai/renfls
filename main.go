// Copyright © 2017 shoarai

// The DirNameToFileName re-name the files in a directory to the directory.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"./monitored"}
	}

	for _, root := range roots {
		walkDirectory(root)
	}
}

func walkDirectory(dir string) {
	for _, entry := range getEntry(dir) {
		if entry.IsDir() {
			var folderName = entry.Name()
			subdir := filepath.Join(dir, entry.Name())
			dirNameToFileName(subdir, folderName)
		}
	}
}

func dirNameToFileName(dir string, name string) {
	// TODO: Add suffix in case of same extension.
	// TODO: Ignore "." file.
	// TODO: Make log file.
	// TODO: Remove the old directory.
	for _, entry := range getEntry(dir) {
		if !entry.IsDir() {
			oldFileName := entry.Name()
			pos := strings.LastIndex(oldFileName, ".")
			extension := oldFileName[pos:]

			newFileName := name + extension

			oldFilePath := dir + "/" + oldFileName
			newFilePath := dir + "/../" + newFileName
			os.Rename(oldFilePath, newFilePath)

			fmt.Println(dir + "/")
			fmt.Println("  ", oldFileName, "->", newFileName)
		}
	}
}

func getEntry(dir string) []os.FileInfo {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
