package stats

import (
	"io"
	"os"
	"bytes"
	"sync"
	"path/filepath"
	"regexp"
	"errors"
	"encoding/json"
	"log"
	"fmt"
	"github.com/fatih/color"
)

type outputType int

const (
	OutputPretty outputType = iota
	OutputJson
)

type Stats struct {
	wg     sync.WaitGroup
	output outputType
	Path   string
	Data   statData
}

type statData struct {
	TotalFiles     int `json:"total_files"`
	TotalLines     int `json:"total_lines"`
	AvgLinesInFile int `json:"avg_lines"`
	MaxLinesInFile int `json:"max_lines"`
	Files          []FileInfo `json:"files"`
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Lines int `json:"lines"`
}

func New(path string, o outputType) (*Stats, error) {
	if !CheckFileExist(path) {
		return nil, errors.New("Directory not found: " + path)
	}
	s := &Stats{output: o, Path: path}
	filepath.Walk(path, s.visit)
	s.wg.Wait()
	s.Data.AvgLinesInFile = s.Data.TotalLines / s.Data.TotalFiles

	return s, nil
}

func (s *Stats) String() string {
	if (s.output == OutputJson) {
		b, _ := json.Marshal(s.Data)
		return string(b)
	}

	summary := `Stats for %s:
	Files: %s
	Lines: %s
	Max lines in file: %s
	Average lines in file: %s
	`

	yellow := color.New(color.FgYellow).SprintFunc()
	return fmt.Sprintf(summary,
		yellow(s.Path),
		yellow(s.Data.TotalFiles),
		yellow(s.Data.TotalLines),
		yellow(s.Data.MaxLinesInFile),
		yellow(s.Data.AvgLinesInFile),
	)
}

func (s *Stats) visit(path string, fi os.FileInfo, err error) error {
	re := regexp.MustCompile(`.+\/\w+\.go$`)
	if re.MatchString(path) {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if !fi.IsDir() {
				f := FileInfo{Name: filepath.Base(path), Path: path}
				f.Lines, err = f.getLinesCount()
				if err != nil {
					log.Fatal(err)
				}
				s.Data.Files = append(s.Data.Files, f)
				s.Data.TotalFiles += 1
				s.Data.TotalLines += f.Lines
				if f.Lines > s.Data.MaxLinesInFile {
					s.Data.MaxLinesInFile = f.Lines
				}

			}
		}()
	}

	return nil
}

func (fi *FileInfo) getLinesCount() (int, error) {
	f, _ := os.Open(fi.Path)
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
