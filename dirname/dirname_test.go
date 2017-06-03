// Copyright © 2017 shoarai

package dirname_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/shoarai/toDirName/dirname"
)

const testData = "testdata"
const tmpTestData = "." + testData

func TestMain(m *testing.M) {
	createTestDir()
	code := m.Run()
	removeTestDir()
	os.Exit(code)
}

func TestRenameAndMoveFile(t *testing.T) {
	root := getTestDataDir()

	tests := []struct {
		oldPath, newDir, newFileName, wantFileName string
	}{
		{root + "/dir1/text.txt", root, "new", "new.txt"},
		{root + "/dir1/music.mp3", root, "new", "new.mp3"},
		{root + "/dir1/.no", root, "new", "new.no"},
		{root + "/dir1/file", root, "new", "new"},
		{root + "/dir2/a.txt", root, "newText", "newText.txt"},
		{root + "/dir2/b.txt", root, "newText", "newText-1.txt"},
		{root + "/dir2/c.txt", root, "newText", "newText-2.txt"},
		// Rename folder
		{root + "/dir3/dir3-1", root, "newDir", "newDir"},
		{root + "/dir3/dir3-2", root, "newDir", "newDir-1"},
	}

	for _, test := range tests {
		filePath, err := dirname.RenameAndMoveFile(
			test.oldPath, test.newDir, test.newFileName,
		)

		if err != nil {
			t.Errorf("RenameAndMoveFile(%v) error: %s\n", test, err)
			continue
		}

		wantNewPath := filepath.Join(test.newDir, test.wantFileName)
		if !isFileExisting(wantNewPath) {
			t.Errorf("The new file didn't be created.")
		}
		if isFileExisting(test.oldPath) {
			t.Errorf("The old file didn't be removed.")
		}
		if filePath != wantNewPath {
			t.Errorf("RenameAndMoveFile() = %s, want %s", filePath, wantNewPath)
		}
	}
}

func createTestDir() {
	dir, _ := os.Getwd()
	exec.Command("cp", "-r",
		filepath.Join(dir, testData), filepath.Join(dir, tmpTestData),
	).Run()
}

func removeTestDir() {
	dir := getTestDataDir()
	os.RemoveAll(dir)
}

func getTestDataDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, tmpTestData)
}

func isFileExisting(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
