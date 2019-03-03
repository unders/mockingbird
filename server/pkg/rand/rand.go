package rand

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

// Bytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
// Usage:
//         b := make([]byte, 16)
//         if err := rand.Bytes(b); err != nil { return err }
//
func Bytes(b []byte) error {
	// io.ReadFull(rand.Reader, b) reads exactly len(b) bytes into b.
	// It returns the number of bytes copied and an error
	// if fewer bytes were read.
	_, err := io.ReadFull(rand.Reader, b)
	return errors.Wrap(err, "failed to generate random bytes")
}

// Base64URL returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
// Usage:
//         s, err := rand.Base64URL(32)
//         if err != nil {
//             return err
//         }
//
func Base64URL(n int) (string, error) {
	b := make([]byte, n)
	err := Bytes(b)
	return base64.URLEncoding.EncodeToString(b), err
}

// Key32 returns a URL-safe, base64 encoded
// securely generated random string that is 32 bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
// Usage:
//         buf, err := rand.Key32()
//         if err != nil {
//             return err
//         }
//
func Key32() ([]byte, error) {
	b := make([]byte, 32)
	err := Bytes(b)
	str := base64.URLEncoding.EncodeToString(b)
	buf := []byte(str)

	return buf[0:32], err
}
