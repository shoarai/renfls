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
	testDataDir := getTestDataDir()

	for _, test := range []struct {
		oldPath, newDir, newFileName, wantFileName string
	}{
		//Rename file
		{"/dir1/text.txt", ".", "new1", "new1.txt"},
		{"/dir1/image.jpg", ".", "new1", "new1.jpg"},
		{"/dir1/ミュージック　.mp3", ".", "　新　", "　新　.mp3"},
		{"/dir1/.no", ".", "new1", "new1.no"},
		{"/dir1/file", ".", "new1", "new1"},
		{"/dir2/a.txt", ".", "new text2", "new text2.txt"},
		{"/dir2/b.txt", ".", "new text2", "new text2-1.txt"},
		{"/dir2/c.txt", ".", "new text2", "new text2-2.txt"},
		// Rename folder
		{"/dir3/dir3-1", ".", "newDir 3", "newDir 3"},
		{"/dir3/dir3-2", ".", "newDir 3", "newDir 3-1"},
	} {
		test.oldPath = filepath.Join(testDataDir, test.oldPath)
		test.newDir = filepath.Join(testDataDir, test.newDir)

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

func TestRenameAndMoveFileAll(t *testing.T) {
	testDataDir := getTestDataDir()

	for _, test := range []struct {
		root, newDir, newFileName string
		wantFileNames             []string
	}{
		{"dir4", ".", "new4",
			[]string{"new4.txt", "new4-1.txt"}},
	} {
		test.root = filepath.Join(testDataDir, test.root)
		test.newDir = filepath.Join(testDataDir, test.newDir)

		err := dirname.RenameAndMoveFileAll(test.root, test.newDir, test.newFileName)

		if err != nil {
			t.Errorf("RenameAndMoveFileAll(%v) error: %s\n", test, err)
			continue
		}

		for _, wantFileName := range test.wantFileNames {
			wantNewPath := filepath.Join(test.newDir, wantFileName)
			if !isFileExisting(wantNewPath) {
				t.Errorf("The new file didn't be created.")
			}
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
