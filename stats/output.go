package stats

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
)

const (
	OutputPretty = "pretty"
	OutputJson = "json"
)

type Output struct {
	format string
	stats  *Stats
}

func (s *Stats) NewOutput(format string) *Output {
	return &Output{format: format, stats: s}
}

func (s *Stats) String() string {
	o := s.NewOutput(OutputPretty)
	return o.String()
}

func (o *Output) Json() string {
	b, _ := json.Marshal(o.stats.Data)
	return string(b)
}

func (o *Output) String() string {
	if (o.format == OutputJson) {
		return o.Json()
	}
	return o.PrettyString()
}

func (o *Output) PrettyString() string {
	summary := `Stats for %s:
	Files: %s
	Lines: %s
	Max lines in file: %s (%s)
	Average lines in file: %s
	`

	yellow := color.New(color.FgYellow).SprintFunc()
	return fmt.Sprintf(summary,
		yellow(o.stats.Path),
		yellow(o.stats.Data.TotalFiles),
		yellow(o.stats.Data.TotalLines),
		yellow(o.stats.Data.MaxLinesFile.Lines),
		yellow(o.stats.Data.MaxLinesFile.Path),
		yellow(o.stats.Data.AvgLinesInFile),
	)
}
