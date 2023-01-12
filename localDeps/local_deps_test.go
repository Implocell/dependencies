package localdeps_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/implocell/dependencies/localdeps"
)

const JAVSCRIPT_DIR = "../testdata/javascript"

func TestLocalDeps(t *testing.T) {
	t.Run("it fetches testdata", func(t *testing.T) {
		_, err := os.ReadDir(JAVSCRIPT_DIR)
		if err != nil {
			t.Fatal("failed to read testdata folder")
		}
	})
	t.Run("it finds file 2", func(t *testing.T) {
		d, err := os.ReadDir(JAVSCRIPT_DIR)
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

		d, _ := os.ReadDir(JAVSCRIPT_DIR)
		path, _ := filepath.Abs(JAVSCRIPT_DIR)

		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "2.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file.Name(), path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: filepath.Join(path, file.Name()),
			Imports: []localdeps.Import{
				{
					Filename: filepath.Join(path, "1.js"),
					Imported: []string{"hello"},
				},
				{
					Filename: filepath.Join(path, "3.js"),
					Imported: []string{"Nothing", "Something"},
				},
			},
			Exports: nil,
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})

	t.Run("it finds exports in file 1", func(t *testing.T) {

		d, _ := os.ReadDir(JAVSCRIPT_DIR)
		path, _ := filepath.Abs(JAVSCRIPT_DIR)
		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "1.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file.Name(), path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: filepath.Join(path, file.Name()),
			Imports:  []localdeps.Import{},
			Exports:  []localdeps.Exports{{Export: "hello"}},
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})

	t.Run("it finds exports in file 3", func(t *testing.T) {

		d, _ := os.ReadDir(JAVSCRIPT_DIR)
		path, _ := filepath.Abs(JAVSCRIPT_DIR)
		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "3.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file.Name(), path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: filepath.Join(path, file.Name()),
			Imports:  []localdeps.Import{},
			Exports: []localdeps.Exports{
				{
					Export: "Nothing",
				},
				{
					Export: "Something",
				},
				{
					Export: "someone",
				},
			},
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})

	t.Run("it finds exports and imports in file 4", func(t *testing.T) {

		d, _ := os.ReadDir(JAVSCRIPT_DIR)
		path, _ := filepath.Abs(JAVSCRIPT_DIR)
		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "4.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file.Name(), path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: filepath.Join(path, file.Name()),
			Imports: []localdeps.Import{
				{
					Filename: filepath.Join(path, "1.js"),
					Imported: []string{"hello"},
				},
				{
					Filename: filepath.Join(path, "3.js"),
					Imported: []string{"Nothing", "Something"},
				}},
			Exports: []localdeps.Exports{
				{
					Export: "what",
				},
				{
					Export: "krumspring",
				},
				{
					Export: "nojo",
				},
			},
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})
	t.Run("it finds exports and imports in components 1", func(t *testing.T) {

		d, _ := os.ReadDir(JAVSCRIPT_DIR + "/components")
		path, _ := filepath.Abs(JAVSCRIPT_DIR + "/components")
		var file fs.DirEntry
		for _, f := range d {
			if f.Name() == "comp1.js" {
				file = f
			}
		}

		res, err := localdeps.Find(file.Name(), path)

		if err != nil {
			t.Fatalf("error from find function: %s\n", err)
		}

		expected := localdeps.FileDeps{
			Filename: filepath.Join(path, file.Name()),
			Imports: []localdeps.Import{
				{
					Filename: filepath.Join(path, "../3.js"),
					Imported: []string{"Nothing"},
				}},
			Exports: nil,
		}

		if equal := reflect.DeepEqual(res, expected); !equal {
			t.Fatal("structs are not equal")
		}

	})
}
