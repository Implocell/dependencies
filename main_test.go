package main_test

import (
	"fmt"
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

		res, err := localdeps.FindAll("./testdata/javascript")

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
			`[{"filename":"%s","imports":[],"exports":[{"export":"hello","usedBy":["%s","%s"]}]},{"filename":"%s","imports":[{"filename":"%s","imported":["hello"]},{"filename":"%s","imported":["Nothing","Something"]}],"exports":null},{"filename":"%s","imports":[],"exports":[{"export":"Nothing","usedBy":["%s","%s","%s"]},{"export":"Something","usedBy":["%s","%s"]},{"export":"someone","usedBy":null}]},{"filename":"%s","imports":[{"filename":"%s","imported":["hello"]},{"filename":"%s","imported":["Nothing","Something"]}],"exports":[{"export":"what","usedBy":null},{"export":"krumspring","usedBy":null},{"export":"nojo","usedBy":null}]},{"filename":"%s","imports":[{"filename":"%s","imported":["Nothing"]}],"exports":null}]`,
			filepath.Join(JAVSCRIPT_DIR, "1.js"),
			filepath.Join(JAVSCRIPT_DIR, "2.js"),
			filepath.Join(JAVSCRIPT_DIR, "4.js"),
			filepath.Join(JAVSCRIPT_DIR, "2.js"),
			filepath.Join(JAVSCRIPT_DIR, "1.js"),
			filepath.Join(JAVSCRIPT_DIR, "3.js"),
			filepath.Join(JAVSCRIPT_DIR, "3.js"),
			filepath.Join(JAVSCRIPT_DIR, "2.js"),
			filepath.Join(JAVSCRIPT_DIR, "4.js"),
			filepath.Join(JAVSCRIPT_DIR, "/components/comp1.js"),
			filepath.Join(JAVSCRIPT_DIR, "2.js"),
			filepath.Join(JAVSCRIPT_DIR, "4.js"),
			filepath.Join(JAVSCRIPT_DIR, "4.js"),
			filepath.Join(JAVSCRIPT_DIR, "1.js"),
			filepath.Join(JAVSCRIPT_DIR, "3.js"),
			filepath.Join(JAVSCRIPT_DIR, "/components/comp1.js"),
			filepath.Join(JAVSCRIPT_DIR, "3.js"),
		)

		if string(f) != expected {
			t.Fatalf("failed to validate string, expected %s, but got %s", expected, string(f))
		}

	})
}
