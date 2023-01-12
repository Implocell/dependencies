package writer_test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/implocell/dependencies/writer"
)

func TestJSONWriter(t *testing.T) {
	t.Run("it writes json data", func(t *testing.T) {
		data := struct {
			FileName string `json:"filename"`
			Hello    string `json:"hello"`
		}{
			FileName: "special",
			Hello:    "world",
		}

		file, err := ioutil.TempFile(".", "*")
		defer os.Remove(file.Name())

		if err != nil {
			log.Fatal(err)
		}

		err = writer.JSONFile(data, file)

		if err != nil {
			t.Fatal(err)
		}

		f, err := os.ReadFile(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		res := strings.Contains(string(f), `{"filename":"special","hello":"world"}`)

		if !res {
			t.Fatal("couldn't find correct data")
		}

	})
}
