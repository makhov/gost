package stats

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type outputType int

const (
	OutputPretty outputType = iota
	OutputJson
)

type Stats struct {
	//wg     sync.WaitGroup
	output outputType
	Path   string
	Data   statData
}

type statData struct {
	TotalFiles     int        `json:"total_files"`
	TotalLines     int        `json:"total_lines"`
	AvgLinesInFile int        `json:"avg_lines"`
	MaxLinesInFile int        `json:"max_lines"`
	Files          []FileInfo `json:"files"`
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Lines int    `json:"lines"`
}

func New(path string, o outputType) (*Stats, error) {
	if !CheckFileExist(path) {
		return nil, errors.New("Directory not found: " + path)
	}
	s := &Stats{output: o, Path: path}

	done := make(chan struct{})
	defer close(done)

	fi, errc := sumFiles(done, path)
	for f := range fi {
		s.Data.Files = append(s.Data.Files, f)
		s.Data.TotalLines += f.Lines
		if f.Lines > s.Data.MaxLinesInFile {
			s.Data.MaxLinesInFile = f.Lines
		}
	}
	s.Data.TotalFiles = len(s.Data.Files)
	s.Data.AvgLinesInFile = s.Data.TotalLines / s.Data.TotalFiles
	if err := <-errc; err != nil {
		return nil, err
	}

	return s, nil
}

func sumFiles(done <-chan struct{}, root string) (<-chan FileInfo, <-chan error) {
	c := make(chan FileInfo)
	errc := make(chan error, 1)
	re := regexp.MustCompile(`.+\/\w+\.go$`)
	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if re.MatchString(path) && !info.IsDir() {
				if err != nil {
					return err
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				wg.Add(1)
				go func() {
					f := FileInfo{Name: filepath.Base(path), Path: path}
					f.Lines, err = f.getLinesCount()
					select {
					case c <- f:
					case <-done:
					}
					wg.Done()
				}()

				select {
				case <-done:
					return errors.New("walk canceled")
				default:
					return nil
				}
			}
			return nil
		})

		go func() {
			wg.Wait()
			close(c)
		}()

		errc <- err
	}()
	return c, errc
}

func (s *Stats) String() string {
	if s.output == OutputJson {
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
