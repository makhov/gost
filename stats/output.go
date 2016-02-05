package stats

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
)

func (s *Stats) String() string {
	if s.output == OutputJson {
		return s.Json()
	}
	return s.PrettyString()
}

func (s *Stats) Json() string {
	b, _ := json.Marshal(s.Data)
	return string(b)
}

func (s *Stats) PrettyString() string {
	summary := `Stats for %s:
	Files: %s
	Lines: %s
	Max lines in file: %s (%s)
	Average lines in file: %s
	`

	yellow := color.New(color.FgYellow).SprintFunc()
	return fmt.Sprintf(summary,
		yellow(s.Path),
		yellow(s.Data.TotalFiles),
		yellow(s.Data.TotalLines),
		yellow(s.Data.MaxLinesFile.Lines),
		yellow(s.Data.MaxLinesFile.Path),
		yellow(s.Data.AvgLinesInFile),
	)
}
