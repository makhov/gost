package stats

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Stats struct {
	Path   string
	Data   statData
}

type statData struct {
	TotalFiles     int        `json:"total_files"`
	TotalLines     int        `json:"total_lines"`
	AvgLinesInFile int        `json:"avg_lines"`
	MaxLinesFile   FileInfo   `json:"max_lines_file"`
	Files          []FileInfo `json:"files"`
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Lines int    `json:"lines"`
}

func New(path string) (*Stats, error) {
	if !CheckFileExist(path) {
		return nil, errors.New("Directory not found: " + path)
	}
	s := &Stats{Path: path}

	done := make(chan struct{})
	defer close(done)

	fi, errc := sumFiles(done, path)
	for f := range fi {
		s.Data.Files = append(s.Data.Files, f)
		s.Data.TotalLines += f.Lines
		if f.Lines > s.Data.MaxLinesFile.Lines {
			s.Data.MaxLinesFile = f
		}
	}
	s.Data.TotalFiles = len(s.Data.Files)
	if s.Data.TotalFiles > 0 {
		s.Data.AvgLinesInFile = s.Data.TotalLines / s.Data.TotalFiles
	}
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
