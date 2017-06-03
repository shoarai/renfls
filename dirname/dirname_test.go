// Copyright © 2017 shoarai

package dirname_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/shoarai/toDirName/dirname"
)

func TestMain(m *testing.M) {
	copyTestDir()
	code := m.Run()
	removeTestDir()
	os.Exit(code)
}

func TestRenameAndMoveFile(t *testing.T) {
	dir := getTestDataDir()

	tests := []struct {
		oldDir, oldFileName, newDir, newFileName, extension string
	}{
		{dir + "/dir1", "text.txt", dir, "newName", ".txt"},
	}

	for _, test := range tests {
		filePath, err := dirname.RenameAndMoveFile(
			test.oldDir, test.oldFileName, test.newDir, test.newFileName,
		)
		if err != nil {
			t.Errorf("RenameAndMoveFile(%v) error: %s\n", test, err)
			continue
		}

		if !isExisting(filePath) {
			t.Errorf("New file create error: %s", err)
		}
		if isExisting(test.oldDir + "/" + test.oldFileName) {
			t.Errorf("New file create error: %s", err)
		}
	}
}

func getTestDataDir() string {
	dir, _ := os.Getwd()
	return dir + "/.testdata"
}

func copyTestDir() {
	dir, _ := os.Getwd()
	exec.Command("cp", "-r", dir+"/testdata", dir+"/.testdata").Run()
}

func removeTestDir() {
	dir := getTestDataDir()
	os.RemoveAll(dir)
}

func isExisting(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
