package main

import (
	"github.com/makhov/gost/stats"
	"path/filepath"
	"fmt"
	"log"
	"flag"
)

func main() {
	path := flag.String("path", ".", "Path")
	outputFlag := flag.String("output", "pretty", "Output format")
	flag.Parse()

	dir, _ := filepath.Abs(*path)

	o := stats.OutputPretty
	if *outputFlag == "json" {
		o = stats.OutputJson
	}

	s, err := stats.New(dir, o)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}


