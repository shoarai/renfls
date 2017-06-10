// Copyright © 2017 shoarai

package renfls_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoarai/renfls"
)

const (
	tmpTestData = ".testdata"
)

func TestMain(m *testing.M) {
	createTestDir()
	code := m.Run()
	// removeTestDir()
	os.Exit(code)
}

func TestRename(t *testing.T) {
	for _, test := range []struct {
		oldPath, newDir, newName, wantFileName string
	}{
		// Rename file
		{"dir1/text.txt", ".", "new1", "new1.txt"},
		{"dir1/image.jpg", ".", "new1", "new1.jpg"},
		{"dir1/ミュージック　.mp3", ".", "　新　", "　新　.mp3"},
		{"dir1/.no", ".", "new1", "new1.no"},
		{"dir1/file", ".", "new1", "new1"},
		{"dir2/a.txt", ".", "new text2", "new text2.txt"},
		{"dir2/b.txt", ".", "new text2", "new text2-2.txt"},
		{"dir2/c.txt", ".", "new text2", "new text2-3.txt"},
		// Rename folder
		{"dir3/dir3-1", ".", "newDir 3", "newDir 3"},
		{"dir3/dir3-2", ".", "newDir 3", "newDir 3-2"},
	} {
		createAll(test.oldPath)

		filePath, err := renfls.Rename(
			test.oldPath, test.newDir, test.newName,
		)

		if err != nil {
			t.Errorf("Rename(%v) error: %s\n", test, err)
		}

		wantNewPath := filepath.Join(test.newDir, test.wantFileName)
		if !isExist(wantNewPath) {
			t.Errorf("The new file %q didn't be created.", wantNewPath)
		}
		if isExist(test.oldPath) {
			t.Errorf("The old file %q didn't be removed.", test.oldPath)
		}
		if filePath != wantNewPath {
			t.Errorf("Rename() = %s, want %s", filePath, wantNewPath)
		}

		// clearTestDir()
	}
}

func TestRenameAll(t *testing.T) {
	for _, test := range []struct {
		root, newDir, newFileName string
		files                     []string
		wantFileNames             []string
	}{
		{"dir4", ".", "new4",
			[]string{"dir4-1/text.txt", "dir4-2/text.txt", "dir4-2/画像.jpg"},
			[]string{"new4.txt", "new4-2.txt", "new4.jpg"}},
	} {
		createAlls(test.root, test.files)

		err := renfls.RenameAll(test.root, test.newDir, test.newFileName)

		if err != nil {
			t.Errorf("RenameAll(%v) error: %s\n", test, err)
		}

		for _, wantFileName := range test.wantFileNames {
			wantNewPath := filepath.Join(test.newDir, wantFileName)
			if !isExist(wantNewPath) {
				t.Errorf("The new path %q didn't be created.\n", wantNewPath)
			}
		}

		clearTestDir()
	}
}

func TestRenamePattern(t *testing.T) {
	for _, test := range []struct {
		root, newDir, newFileName, pattern string
		files                              []string
		wantFileNames                      []string
		wantRemovedFileNames               []string
	}{
		{"dir4", ".", "new4", `text*`,
			[]string{"dir4-1/text.txt", "dir4-2/text.txt", "dir4-2/image.jpg"},
			[]string{"new4.txt", "new4-2.txt"},
			[]string{"new4.jpg"}},
	} {
		createAlls(test.root, test.files)

		err := renfls.RenamePattern(
			test.root, test.newDir, test.newFileName, test.pattern)

		if err != nil {
			t.Errorf("RenameAll(%v) error: %s\n", test, err)
		}

		for _, wantFileName := range test.wantFileNames {
			wantNewPath := filepath.Join(test.newDir, wantFileName)
			if !isExist(wantNewPath) {
				t.Errorf("The new path %q didn't be created.\n", wantNewPath)
			}
		}

		for _, wantRemovedFileName := range test.wantRemovedFileNames {
			wantNewPath := filepath.Join(test.newDir, wantRemovedFileName)
			if isExist(wantNewPath) {
				t.Errorf("The path not matched %q is created.\n", wantNewPath)
			}
		}

		clearTestDir()
	}
}

func TestToDirNames(t *testing.T) {
	createTestData()

	for _, test := range []struct {
		root      string
		wantFiles []string
	}{
		{"root", []string{
			"image.jpg",
			"dir1.txt",
			"dir1.jpg",
			"dir2.txt",
			"dir2-2.txt",
			"ディレクトリ3.txt",
			"ディレクトリ3-2.txt",
			"ディレクトリ3-3.txt",
		}},
	} {
		e := renfls.ToDirNames(test.root)
		if e != nil {
			t.Errorf("ToDirNames(%v) error: %s\n", test.root, e)
		}

		for _, want := range test.wantFiles {
			wantPath := filepath.Join(test.root, want)
			if isNotFileExist(wantPath) {
				t.Errorf("The path %q didn't be created.\n", wantPath)
			}
		}

		clearTestDir()
	}
}

func TestToDirNamesPattern(t *testing.T) {
	createTestData()

	for _, test := range []struct {
		root, pattern string
		wantFiles     []string
		wantNotExists []string
	}{
		{"root", `.txt$`, []string{
			"dir1.txt",
			"dir2.txt",
			"dir2-2.txt",
			"ディレクトリ3.txt",
			"ディレクトリ3-2.txt",
			"ディレクトリ3-3.txt",
		}, []string{
			"image.jpg",
			"dir1.jpg",
		}},
	} {
		e := renfls.ToDirNamesPattern(test.root, test.pattern)
		if e != nil {
			t.Errorf("ToDirNamesPattern(%v) error: %s\n", test.root, e)
		}

		for _, want := range test.wantFiles {
			wantPath := filepath.Join(test.root, want)
			if isNotFileExist(wantPath) {
				t.Errorf("The path %q didn't be created.\n", want)
			}
		}
		for _, want := range test.wantNotExists {
			wantPath := filepath.Join(test.root, want)
			if !isNotFileExist(wantPath) {
				t.Errorf("The path %q is created.\n", want)
			}
		}

		clearTestDir()
	}
}

func getFiles(dir string) []string {
	dirs, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil
	}
	var strs []string
	for _, dir := range dirs {
		strs = append(strs, dir.Name())
	}
	return strs
}

func equalNoOrder(strs1, strs2 []string) (string, bool) {
	if len(strs1) != len(strs2) {
		return "", false
	}
	for _, s1 := range strs1 {
		if !contains(strs2, s1) {
			return s1, false
		}
	}
	for _, s2 := range strs2 {
		if !contains(strs1, s2) {
			return s2, false
		}
	}
	return "", true
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

type DirTree struct {
	dir    string
	files  []string
	subDir []*DirTree
}

func createTestData() []*DirTree {
	data := []*DirTree{
		{"root", []string{
			"image.jpg"}, []*DirTree{
			{"dir1", []string{
				"text.txt",
				"image.jpg",
			}, nil},
			{"dir2", []string{
				"text2.txt",
				"text2-1.txt",
			}, nil},
			{"ディレクトリ3", []string{
				"テキスト.txt",
			}, []*DirTree{
				{"dir3-1", []string{
					"text3-1.txt",
					"あいうえお　.txt",
				}, nil}},
			},
			{"dir4", nil, nil},
		}},
	}
	createDirTrees(data)
	return data
}

func createDirTrees(trees []*DirTree) {
	createDirTree("", trees)
}

func createDirTree(root string, trees []*DirTree) {
	if trees == nil {
		return
	}
	for _, tree := range trees {
		dir := filepath.Join(root, tree.dir)
		create(dir, tree.files)
		createDirTree(dir, tree.subDir)
	}
}

func create(dir string, files []string) {
	os.Mkdir(dir, os.ModePerm)
	for _, f := range files {
		os.Create(filepath.Join(dir, f))
	}
}

func createAll(path string) {
	dir, _ := filepath.Split(path)
	os.MkdirAll(dir, os.ModePerm)
	os.Create(path)
}

func createAlls(root string, path []string) {
	for _, p := range path {
		createAll(filepath.Join(root, p))
	}
}

func createTestDir() {
	os.Mkdir(tmpTestData, os.ModePerm)
	os.Chdir(tmpTestData)
}

func removeTestDir() {
	os.Chdir("..")
	os.RemoveAll(tmpTestData)
}

func clearTestDir() {
	removeTestDir()
	createTestDir()
}

func isFileExist(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func isNotFileExist(path string) bool {
	return !isFileExist(path)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isNotExist(path string) bool {
	return !isExist(path)
}
