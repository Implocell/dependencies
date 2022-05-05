package localdeps

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Import struct {
	Filename string   `json:"filename"`
	Imported []string `json:"imported"`
}

type FileDeps struct {
	Filename string   `json:"filename"`
	Imports  []Import `json:"imports"`
	Exports  []string `json:"exports"`
}

func Find(file fs.DirEntry, path string) (FileDeps, error) {
	filePath := filepath.Join(path, file.Name())

	content, err := os.ReadFile(filePath)
	if err != nil {
		return FileDeps{}, err
	}

	imports := findImports(content, path)

	exports := findExports(content)

	fileDeps := FileDeps{
		Filename: filePath,
		Imports:  imports,
		Exports:  exports,
	}

	return fileDeps, nil
}

func findImports(content []byte, path string) []Import {
	re := regexp.MustCompile(`import {?(.*?)}? from ['|"](.*)['|"]`)

	strs := re.FindAllStringSubmatch(string(content), -1)

	if len(strs) == 0 {
		return []Import{}
	}

	if len(strs[0]) <= 1 {
		return []Import{}
	}

	var imports []Import

	for _, sub := range strs {

		imported := strings.Split(strings.ReplaceAll(sub[1], " ", ""), ",")
		importFrom, err := findImportFile(sub[len(sub)-1], path)
		if err != nil {
			return []Import{}
		}

		imports = append(imports, Import{
			Filename: importFrom,
			Imported: imported,
		})

	}
	return imports
}

func findExports(content []byte) []string {
	re := regexp.MustCompile(`export default (\w*)|export \w* (\w+)`)

	strs := re.FindAllStringSubmatch(string(content), -1)

	if len(strs) == 0 {
		return nil
	}

	if len(strs[0]) <= 1 {
		return nil
	}

	var exports []string

	for _, sub := range strs {

		var exported []string

		if len(sub[2]) > 0 {
			exports = append(exports, sub[2])
		}

		if len(sub[1]) > 0 {
			exports = append(exports, sub[1])
		}

		exports = append(exports, exported...)

	}
	return exports
}

//!OH BOY NEEDS REWRITE
func findImportFile(f, path string) (string, error) {
	extensions := [2]string{".js", ".jsx"}

	fullPath := filepath.Join(path, f)

	fileName := filepath.Base(fullPath)

	folderPath := filepath.Dir(fullPath)

	d, err := os.ReadDir(folderPath)

	if err != nil {
		return "", err
	}
	for _, f := range d {
		if ok := strings.Contains(f.Name(), fileName); ok {
			for _, ext := range extensions {
				_, err := os.Stat(fullPath + ext)
				if err != nil {
					continue
				}
				return fullPath + ext, nil
			}
		}
	}
	return "", nil
}

func FindAll(root string) ([]FileDeps, error) {
	var fileDeps []FileDeps
	d, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	for _, f := range d {
		fileDep, err := Find(f, path)
		if err != nil {
			continue
		}
		fileDeps = append(fileDeps, fileDep)
	}

	if len(fileDeps) == 0 {
		return nil, fmt.Errorf("failed to find any files")
	}

	return fileDeps, nil
}
