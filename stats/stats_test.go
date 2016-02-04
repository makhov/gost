package stats

import (
	"path/filepath"
	"testing"
	"github.com/labstack/gommon/log"
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

var testData = `{"total_files":2,"total_lines":190,"avg_lines":95,"max_lines":169,"files":[{"name":"stats.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats.go","lines":169},{"name":"stats_test.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats_test.go","lines":21}]}`
