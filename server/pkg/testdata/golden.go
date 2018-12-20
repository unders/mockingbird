package testdata

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// ReadGolden reads the file named by filename and returns the content.
func ReadGolden(t *testing.T, filename string) []byte {
	t.Helper()

	golden := filepath.Join("testdata/golden", filename+".golden")
	buf, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatalf("\nWant: nil\n Got: %s", err)
	}
	return buf
}

// WriteGolden writes data to a file named by filename.
func WriteGolden(t *testing.T, filename string, data []byte) {
	t.Helper()

	golden := filepath.Join("testdata/golden", filename+".golden")
	err := ioutil.WriteFile(golden, data, 0644)
	if err != nil {
		t.Fatalf("\nWant: nil\n Got: %s", err)
	}
}
