package main

import (
	"flag"
	"log"
	"os"

	"github.com/edfun317/ereader/internal/format/epub"
	"github.com/edfun317/ereader/internal/viewer/cli"
)

func main() {

	run()

}

func run() {

	filePath := flag.String("file", "", "Path to EPUB file")
	flag.Parse()

	if *filePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	reader := epub.NewEPUBReader()
	viewer := cli.NewCLIViewer(reader)
	if err := viewer.Start(*filePath); err != nil {
		log.Fatal(err)
	}

}
