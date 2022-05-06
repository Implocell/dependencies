package localdeps

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Exports struct {
	Export string   `json:"export"`
	UsedBy []string `json:"usedBy"`
}

type Import struct {
	Filename string   `json:"filename"`
	Imported []string `json:"imported"`
}

type FileDeps struct {
	Filename string    `json:"filename"`
	Imports  []Import  `json:"imports"`
	Exports  []Exports `json:"exports"`
}

func Find(fileName string, path string) (FileDeps, error) {
	filePath := filepath.Join(path, fileName)

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

func findExports(content []byte) []Exports {
	re := regexp.MustCompile(`export default (\w*)|export \w* (\w+)`)

	strs := re.FindAllStringSubmatch(string(content), -1)

	if len(strs) == 0 {
		return nil
	}

	if len(strs[0]) <= 1 {
		return nil
	}

	var exports []Exports

	for _, sub := range strs {

		if len(sub[1]) > 0 {
			exports = append(exports, Exports{Export: sub[1]})
		}

		if len(sub[2]) > 0 {
			exports = append(exports, Exports{Export: sub[2]})
		}

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

	fileDeps, err := findAllByDirectory(root)

	if err != nil {
		fmt.Println(err)
	}

	if len(fileDeps) == 0 {
		return nil, fmt.Errorf("failed to find any files")
	}
	findImportAttachExport(fileDeps)
	return fileDeps, nil
}

func findAllByDirectory(root string) ([]FileDeps, error) {
	var fileDeps []FileDeps

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if !isValidFileType(info.Name()) {
					return nil
				}

				dirPath := filepath.Dir(path)

				fileDep, err := Find(info.Name(), dirPath)
				if err != nil {
					return err
				}

				fileDeps = append(fileDeps, fileDep)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return fileDeps, nil
}

func isValidFileType(filename string) bool {
	extensions := [2]string{".js", ".jsx"}
	var isValid bool
	for _, ext := range extensions {
		if ok := strings.HasSuffix(filename, ext); ok {
			isValid = ok
		}
	}
	return isValid

}

func findImportAttachExport(fileDeps []FileDeps) {
	for _, f := range fileDeps {
		for _, fi := range f.Imports {
			for y, ff := range fileDeps {
				if ff.Filename == fi.Filename {
					for i, ex := range ff.Exports {
						for _, imports := range fi.Imported {
							if strings.Contains(ex.Export, imports) {
								fileDeps[y].Exports[i].UsedBy = append(ff.Exports[i].UsedBy, f.Filename)
							}
						}
					}
				}
			}
		}
	}
}
