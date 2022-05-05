package main

import (
	"log"
	"os"

	"github.com/implocell/cleanc/localdeps"
	"github.com/implocell/cleanc/writer"
)

func main() {
	res, err := localdeps.FindAll("./testdata/javascript")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("results.json")

	if err != nil {
		log.Fatal(err)
	}

	err = writer.JSONFile(res, f)
	if err != nil {
		log.Fatal(err)
	}
}
