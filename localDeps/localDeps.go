package localdeps

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Import struct {
	Filename  string   `json:"filename"`
	Functions []string `json:"functions"`
}

type FileDeps struct {
	Filename string   `json:"filename"`
	Imports  []Import `json:"imports"`
}

func Find(file fs.DirEntry, path string) (FileDeps, error) {
	res := filepath.Join(path, file.Name())
	content, err := os.ReadFile(res)
	if err != nil {
		return FileDeps{}, err
	}

	re, err := regexp.Compile(`import (\w+) from '(.*)'`)
	if err != nil {
		return FileDeps{}, err
	}

	strs := re.FindAllStringSubmatch(string(content), -1)
	if len(strs[0]) <= 1 {
		return FileDeps{}, fmt.Errorf("failed to find anything in file %s", res)
	}

	fileDeps := FileDeps{
		Filename: file.Name(),
		Imports: []Import{
			{
				Filename:  strs[0][len(strs[0])-1] + ".js",
				Functions: strs[0][1 : len(strs[0])-1],
			},
		},
	}

	return fileDeps, nil
}
