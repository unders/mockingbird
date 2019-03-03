package rand

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/unders/mockingbird/server/pkg/testdata"
)

// URLEncoding is the alternate base64 encoding
// defined in RFC 4648:
// https://tools.ietf.org/html/rfc4648
// The "URL and Filename safe" Base 64 Alphabet
// 0 - 61: [A-Za-z0-9]
// 62: + becomes: -
// 63: / becomes: _
// padding: =
// This encoding may be referred to as "base64url".
var base64URLFormat = regexp.MustCompile(
	`^[\_\-a-zA-Z\d]+={0,3}$`,
)

func TestBytes_mustBeUnique(t *testing.T) {
	b := make([]byte, 20*100)

	for i := 0; i < 100; i++ {
		subslice := make([]byte, 20)
		err := Bytes(b)
		testdata.AssertNil(t, err)

		if bytes.Contains(b, subslice) {
			t.Error("random subslice should be unique")
		}

		n := i * 20
		copy(b[n:], subslice)
	}
}

func TestBase64URL(t *testing.T) {
	for i := 0; i < 100; i++ {
		s, err := Base64URL(32)

		testdata.AssertNil(t, err)

		if len(s) != 44 {
			t.Error("Expected string of len 44, got: ", len(s))
		}

		if !base64URLFormat.MatchString(s) {
			t.Error("Expected base64 encoded string, got: ", s)
		}
	}
}

func TestKey32_shouldHaveCorrectLength(t *testing.T) {
	b, err := Key32()
	testdata.AssertNil(t, err)

	l := len(b)

	if 32 != l {
		t.Error("Expected key length 32; got: ", l, string(b))
	}
}
