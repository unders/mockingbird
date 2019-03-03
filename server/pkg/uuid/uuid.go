// Package uuid generates unique UUIDs of V4 format:
//
//      xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//

package uuid

import (
	"encoding/hex"

	"github.com/unders/mockingbird/server/pkg/rand"
)

const dash = byte('-')

// NewV4 returns a securely generated UUID of V4 format.
//
// Usage:
//
//         buf, err := uuid.NewV4()
//         if err != nil {
//             return err
//         }
//
//
// Returns:
//
//         * buf
//                 A UUID representation compliant with specification
//                 described in RFC 4122:
//                 xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
//
//         * error
//                 if the system's secure random number
//                 generator fails to function correctly, in which
//                 case the caller should not continue.
func NewV4() ([]byte, error) {
	src := [16]byte{}
	if err := rand.Bytes(src[:]); err != nil {
		return nil, err
	}

	dst := [36]byte{}
	encode(dst[:], src[:])

	return dst[:], nil
}

func encode(dst []byte, src []byte) {
	hex.Encode(dst[0:8], src[0:4])
	dst[8] = dash
	hex.Encode(dst[9:13], src[4:6])
	dst[13] = dash
	hex.Encode(dst[14:18], src[6:8])
	dst[18] = dash
	hex.Encode(dst[19:23], src[8:10])
	dst[23] = dash
	hex.Encode(dst[24:], src[10:])
}
