package stats

import (
	"testing"
	"path/filepath"
)

func TestFiles(t *testing.T)  {
	path, _ := filepath.Abs(".")
	s, err := New(path, OutputJson)
	if err != nil {
		t.Fail()
	}
	if s.String() != testData {
		t.Fail()
	}
}

var testData = `{"total_files":2,"total_lines":155,"avg_lines":77,"max_lines":137,"files":[{"name":"stats_test.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats_test.go","lines":18},{"name":"stats.go","path":"/Users/amakhov/go/src/github.com/makhov/gost/stats/stats.go","lines":137}]}`