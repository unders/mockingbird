package uuid_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/unders/mockingbird/server/pkg/testdata"
	"github.com/unders/mockingbird/server/pkg/uuid"
)

var (
	UUIDFormat = regexp.MustCompile(
		`[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}`,
	)
	format      = "\nValue: '%s'\nError: Is not a valid UUID V4 format"
	sliceFormat = "\nExpected: slice == s\nslice: %s\ns:     %s"
)

func TestNewV4(t *testing.T) {
	buf, err := uuid.NewV4()
	testdata.AssertNil(t, err)

	if !UUIDFormat.Match(buf) {
		t.Fatalf(format, buf)
	}

	str := string(buf)
	if !UUIDFormat.MatchString(str) {
		t.Fatalf(format, str)
	}

	want := 36
	got := len(str)
	if want != got {
		t.Fatalf("\nWant: %d\n Got: %d\n", want, got)
	}
}

func TestNewV4_ConvertFromSliceToString(t *testing.T) {
	slice, err := uuid.NewV4()
	testdata.AssertNil(t, err)

	s := string(slice)

	if !bytes.Equal(slice, []byte(s)) {
		t.Fatalf(sliceFormat, slice, s)
	}

	if s != string(slice) {
		t.Fatalf(sliceFormat, slice, s)
	}
}

// go test -bench=. github.com/unders/mockingbird/server/pkg/uuid
//
func BenchmarkNewV4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = uuid.NewV4()
	}
}
