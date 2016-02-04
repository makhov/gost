package main

import (
	"flag"
	"fmt"
	"github.com/makhov/gost/stats"
	"log"
	"path/filepath"
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
