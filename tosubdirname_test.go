// Copyright © 2017 shoarai

package renfls_test

import (
	"path/filepath"
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
			[]string{"dir/text.txt", "dir/image.jpg"},
			"root", ".", renfls.Condition{},
			[]string{"dir.txt", "dir.jpg"},
			[]string{},
		},
		{
			[]string{"dir/text.txt", "dir/image.jpg"},
			"root", ".", renfls.Condition{Exts: []string{"txt"}},
			[]string{"dir.txt"},
			[]string{"ignore/dir/image.jpg"},
		},
		{
			[]string{"dir/text.txt", "dir/image.jpg"},
			"root", ".", renfls.Condition{Exts: []string{"txt"}, Ignore: true},
			[]string{"dir.jpg"},
			[]string{"ignore/dir/text.txt"},
		},
	} {
		createAlls(test.root, test.mockFiles)

		err := renfls.WalkToRootSubDirName(test.root, test.dest, test.condition)
		if err != nil {
			t.Errorf("WalkRename(%v) error: %s\n", test, err)
		}

		for _, path := range test.wantRenamedFilePaths {
			wantNewPath := filepath.Join(test.dest, path)
			if !isFileExist(wantNewPath) {
				t.Errorf("The new path %q didn't be created.\n", path)
			}
		}

		for _, path := range test.wantIgnoredFilePaths {
			ignorePath := filepath.Join(test.root, path)
			if !isFileExist(ignorePath) {
				t.Errorf("The path %q is not in ignore directory.\n", path)
			}
		}

		clearTestDir()
	}
}
