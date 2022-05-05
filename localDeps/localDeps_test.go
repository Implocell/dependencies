package localdeps_test

import (
	"io/fs"
	"os"
	"reflect"
	"testing"

	localdeps "github.com/implocell/cleanc/localDeps"
)

func TestLocalDeps(t *testing.T) {
	t.Run("it fetches testdata", func(t *testing.T) {
		_, err := os.ReadDir("../testdata")
		if err != nil {
			t.Fatal("failed to read testdata folder")
		}
	})
	t.Run("it finds file 2", func(t *testing.T) {
		d, err := os.ReadDir("../testdata")
		if err != nil {
			t.Fatal("failed to read testdata folder")
		}
		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "2.js" {
				file = f
			}
		}

		if file == nil {
			t.Fatal("failed to find file")
		}
	})
	t.Run("it finds import in file 2", func(t *testing.T) {
		d, _ := os.ReadDir("../testdata")

		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "2.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file, "../testdata")

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: "2.js",
			Imports: []localdeps.Import{
				{
					Filename:  "./1.js",
					Functions: []string{"hello"},
				},
			},
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})
}
