package stats

import (
	"github.com/labstack/gommon/log"
	"path/filepath"
	"testing"
)

func TestFiles(t *testing.T) {
	path, _ := filepath.Abs(".")
	s, err := New(path, OutputJson)
	if err != nil {
		t.Fail()
	}
	if s.String() != testData {
		log.Error("Unexpected data")
		t.Fail()
	}
}

var testData = `{"total_files":3,"total_lines":204,"avg_lines":68,"max_lines_file":{"name":"stats.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats.go","lines":145},"files":[{"name":"output.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/output.go","lines":38},{"name":"stats.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats.go","lines":145},{"name":"stats_test.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats_test.go","lines":21}]}`
