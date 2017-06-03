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
	dir := getTestDataDir()

	tests := []struct {
		oldDir, oldFileName, newDir, newFileName, wantFileName string
	}{
		{dir + "/dir1", "text.txt", dir, "new", "new.txt"},
		{dir + "/dir1", "music.mp3", dir, "new", "new.mp3"},
		{dir + "/dir1", ".no", dir, "new", "new.no"},
		{dir + "/dir1", "file", dir, "new", "new"},
		{dir + "/dir2", "a.txt", dir, "newText", "newText.txt"},
		{dir + "/dir2", "b.txt", dir, "newText", "newText-1.txt"},
		{dir + "/dir2", "c.txt", dir, "newText", "newText-2.txt"},
		{dir + "/dir3", "dir3-1", dir, "newDir", "newDir"},
		{dir + "/dir3", "dir3-2", dir, "newDir", "newDir-1"},
	}

	for _, test := range tests {
		filePath, err := dirname.RenameAndMoveFile(
			test.oldDir, test.oldFileName, test.newDir, test.newFileName,
		)

		if err != nil {
			t.Errorf("RenameAndMoveFile(%v) error: %s\n", test, err)
			continue
		}

		wantFileName := test.newDir + "/" + test.wantFileName
		if filePath != wantFileName {
			t.Errorf("RenameAndMoveFile() = %s, want %s", filePath, wantFileName)
		}

		if isFileExisting(test.oldDir + "/" + test.oldFileName) {
			t.Errorf("The old file didn't be removed.")
		}
		if !isFileExisting(wantFileName) {
			t.Errorf("The new file didn't be created.")
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
