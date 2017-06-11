// Copyright © 2017 shoarai

package renfls_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoarai/renfls"
)

func TestToDirNames(t *testing.T) {
	createTestData()

	for _, test := range []struct {
		root      string
		wantFiles []string
	}{
		{"root", []string{
			"image.jpg",
			"image.jpg",
			"dir1.txt",
			"dir1.jpg",
			"dir1.mp4",
			"dir1.mp3",
			"dir1.no",
			"dir1",
			"dir2.txt",
			"dir2-2.txt",
			"ディレクトリ3.txt",
			"ディレクトリ3-2.txt",
			"ディレクトリ3-3.txt",
		}},
	} {
		e := renfls.ToSubDirsName(test.root)
		if e != nil {
			t.Errorf("ToDirNames(%v) error: %s\n", test.root, e)
		}

		for _, want := range test.wantFiles {
			wantPath := filepath.Join(test.root, want)
			if isNotFileExist(wantPath) {
				t.Errorf("%q didn't be created.\n", wantPath)
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
			"image.jpg",
			"dir1.txt",
			"dir2.txt",
			"dir2-2.txt",
			"ディレクトリ3.txt",
			"ディレクトリ3-2.txt",
			"ディレクトリ3-3.txt",
		}, []string{
			"dir1.jpg",
			"dir1.mp4",
			"dir1.mp3",
			"dir1.no",
			"dir1",
		}},
	} {
		e := renfls.ToSubDirsNamePattern(test.root, test.pattern)
		if e != nil {
			t.Errorf("ToDirNamesPattern(%v) error: %s\n", test.root, e)
		}

		for _, want := range test.wantFiles {
			wantPath := filepath.Join(test.root, want)
			if isNotFileExist(wantPath) {
				t.Errorf("%q didn't be created.\n", want)
			}
		}
		for _, want := range test.wantNotExists {
			wantPath := filepath.Join(test.root, want)
			if isExist(wantPath) {
				t.Errorf("%q is created.\n", want)
			}
		}

		clearTestDir()
	}
}

func TestToDirNamesIgnoreExt(t *testing.T) {
	createTestData()

	for _, test := range []struct {
		root          string
		exts          []string
		wantFiles     []string
		wantNotExists []string
	}{
		{"root", []string{"jpg", "mp4"}, []string{
			"dir1.txt",
			"dir2.txt",
			"dir1.mp3",
			"dir1.no",
			"dir1",
			"dir2-2.txt",
			"ディレクトリ3.txt",
			"ディレクトリ3-2.txt",
			"ディレクトリ3-3.txt",
		}, []string{
			// "image.jpg",
			"dir1.jpg",
			"dir1.mp4",
		}},
	} {
		e := renfls.ToSubDirsNameIgnoreExt(test.root, test.exts)
		if e != nil {
			t.Errorf("ToSubDirsNameIgnoreExt(%v) error: %s\n", test.root, e)
		}

		for _, want := range test.wantFiles {
			wantPath := filepath.Join(test.root, want)
			if isNotFileExist(wantPath) {
				t.Errorf("%q didn't be created.\n", want)
			}
		}
		for _, want := range test.wantNotExists {
			wantPath := filepath.Join(test.root, want)
			if isExist(wantPath) {
				t.Errorf("%q is created.\n", want)
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
				"movie.mp4",
				"ミュージック　.mp3",
				".no",
				"file",
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
