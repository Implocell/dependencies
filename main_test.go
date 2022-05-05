package main_test

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/implocell/cleanc/localdeps"
	"github.com/implocell/cleanc/writer"
)

const JAVSCRIPT_DIR = "./testdata/javascript"

func TestMain(t *testing.T) {
	t.Run("it finds import in file 2 and writes to file", func(t *testing.T) {
		d, _ := os.ReadDir(JAVSCRIPT_DIR)
		path, _ := filepath.Abs(JAVSCRIPT_DIR)

		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "2.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file, path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		tempFile, err := ioutil.TempFile(".", "*")
		defer os.Remove(tempFile.Name())

		if err != nil {
			t.Fatal(err)
		}

		err = writer.JSONFile(res, tempFile)

		if err != nil {
			t.Fatal(err)
		}

		f, err := os.ReadFile(tempFile.Name())

		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf(
			`{"filename":"%s","imports":[{"filename":"%s","imported":["hello"]},{"filename":"%s","imported":["Nothing","Something"]}],"exports":null}`,
			filepath.Join(
				path, file.Name()),
			filepath.Join(path, "1.js"),
			filepath.Join(path, "3.js"),
		)

		if string(f) != expected {
			t.Fatalf("failed to validate string, expected %s, but got %s", expected, string(f))
		}

	})
}
