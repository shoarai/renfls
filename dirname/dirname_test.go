// Copyright © 2017 shoarai

package dirname_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoarai/toDirName/dirname"
)

const (
	tmpTestData = ".testdata"
	fileMode    = 0777
)

func TestMain(m *testing.M) {
	createTestDir()
	code := m.Run()
	removeTestDir()
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

		filePath, err := dirname.Rename(
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

		err := dirname.RenameAll(test.root, test.newDir, test.newFileName)

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
			[]string{"dir4-1/text.txt", "dir4-2/text.txt", "dir4-2/画像.jpg"},
			[]string{"new4.txt", "new4-2.txt"},
			[]string{"new4.jpg"}},
	} {
		createAlls(test.root, test.files)

		err := dirname.RenamePattern(
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

func createAll(path string) {
	dir, _ := filepath.Split(path)
	os.MkdirAll(dir, fileMode)
	os.Create(path)
}

func createAlls(root string, path []string) {
	for _, p := range path {
		createAll(filepath.Join(root, p))
	}
}

func createTestDir() {
	os.Mkdir(tmpTestData, fileMode)
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

func isExist(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
