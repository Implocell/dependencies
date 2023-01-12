package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/implocell/dependencies/localdeps"
	"github.com/implocell/dependencies/writer"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("target -path <PATH> is required")
	fmt.Println("target -path <PATH> -all - collects all data")
	fmt.Println("target -path <PATH> -unused - finds all empty exports")
	fmt.Println("to save as a file add -save behind arguments, saves it as result.json at current location")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) collectAll(target string, args []string) {
	res, err := localdeps.FindAll(target)
	if err != nil {
		log.Fatal(err)
	}

	for _, arg := range args {
		if strings.Contains(arg, "save") {
			cli.saveToFile(res)
		}
	}

}

func (cli *CommandLine) collectUnusedExports(target string, args []string) {
	res, err := localdeps.FindAll(target)
	if err != nil {
		log.Fatal(err)
	}

	uExports := localdeps.EmptyExports(res)

	for _, arg := range args {
		if strings.Contains(arg, "save") {
			cli.saveToFile(uExports)
		}
	}
}

func (cli *CommandLine) saveToFile(res interface{}) {
	f, err := os.Create("results.json")

	if err != nil {
		log.Fatal(err)
	}

	err = writer.JSONFile(res, f)
	if err != nil {
		log.Fatal(err)
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addTarget := flag.NewFlagSet("target", flag.ExitOnError)
	addTargetDir := addTarget.String("path", "", "Target directory")
	addTargetAll := addTarget.Bool("all", false, "scans all directories")
	addTargetUnused := addTarget.Bool("unused", false, "scans all unused exports")

	switch os.Args[1] {
	case "target":
		err := addTarget.Parse(os.Args[2:])
		if err != nil {
			os.Exit(1)
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addTarget.Parsed() {
		if *addTargetDir == "" {
			fmt.Println("empty target")
			addTarget.Usage()
			runtime.Goexit()
		}
		if *addTargetAll {
			cli.collectAll(*addTargetDir, os.Args)
		}
		if *addTargetUnused {
			cli.collectUnusedExports(*addTargetDir, os.Args)
		}
	}

	for _, arg := range os.Args {
		if strings.Contains(arg, "web") {
			fmt.Println("Serving report at localhost:3000/index.html")
			server(":3000")
		}
	}

}

func handleData(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./results.json")
}

//go:embed frontend/index.html
//go:embed frontend/sorter.js
var content embed.FS

func server(port string) {

	fs := http.FileServer(http.FS(content))
	http.Handle("/", fs)
	http.HandleFunc("/data", handleData)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		runtime.Goexit()
	}

}

func main() {
	defer os.Exit(0)

	cli := CommandLine{}

	cli.run()
}
