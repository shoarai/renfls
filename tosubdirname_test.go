// Copyright © 2017 shoarai

package renfls_test

import (
	"testing"

	"github.com/shoarai/renfls"
)

func TestWalkToRootSubDirName(t *testing.T) {
	for _, test := range []struct {
		mockFiles            []string
		root, dest           string
		condition            renfls.Condition
		wantRenamedFilePaths []string
		wantIgnoredFilePaths []string
	}{
		{
			[]string{"root/dir/text.txt", "root/dir/image.jpg"},
			"root", ".", renfls.Condition{},
			[]string{"dir.txt", "dir.jpg"},
			[]string{},
		},
	} {
		createFiles(test.mockFiles)
		defer clearTestDir()

		err := renfls.WalkToRootSubDirName(test.root, test.dest, test.condition)

		if err != nil {
			t.Errorf("WalkRename(%v) error: %s\n", test, err)
		}

		for _, path := range test.wantRenamedFilePaths {
			if !isFileExist(path) {
				t.Errorf("The new path %q didn't be created.\n", path)
			}
		}

		for _, path := range test.wantIgnoredFilePaths {
			if isFileExist(path) {
				t.Errorf("The path not matched %q is created.\n", path)
			}
		}
	}
}
