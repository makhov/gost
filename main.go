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

	s, err := stats.New(dir)
	if err != nil {
		log.Fatal(err)
	}

	o := s.NewOutput(*outputFlag)
	fmt.Println(o)
}
