package shadowid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"sync/atomic"
)

// defaultSalt is the default salt used to conceal the autoincr and random ID (UUIDv4).
// This acts as a private key.
var defaultSalt atomic.Uint64

// SetSalt sets the default salt.
func SetSalt(s uint64) { defaultSalt.Store(s) }

// ShadowID is a 24-byte ID that conceals the autoincr and random ID (UUIDv4).
// 24 bytes = 8 bytes of autoincr + 16 bytes of UUIDv4
type ShadowID [24]byte

// New creates a new ShadowID from the autoincr and generates a random ID (UUIDv4).
func New(autoincr int64) ShadowID {
	return NewShadowID(autoincr, uuid.New())
}

// NewShadowID creates a new ShadowID from the autoincr and random ID (UUIDv4).
func NewShadowID(autoincr int64, randomid uuid.UUID) ShadowID {
	var id ShadowID

	// NOTE 1: Take 8 bytes from the random ID as the random salt. We can take any 8 bytes from the random ID,
	//         but for this case, we take the last 8 bytes.
	// NOTE 2: We take 8 bytes since the defaultSalt and autoincr are 8 bytes.
	randomSalt := binary.LittleEndian.Uint64(randomid[8:16])

	// NOTE 3: Generate the salted ID by XOR-ing the autoincr, random salt, and default salt.
	// NOTE 4: We XOR because we want to ensure that this ID is reversible.
	salted := uint64(autoincr) ^ randomSalt ^ defaultSalt.Load()

	// NOTE 5: Put the 8 bytes of the UUID's LSB into the first 8 bytes of the ShadowID.
	copy(id[:8], randomid[:8])

	// NOTE 6: Put the 4 bytes of the salted ID's LSB into the 8th-12th byte of the ShadowID.
	// NOTE 7: We use BigEndian because we want to ensure that the byte order matches the hex encoded salted ID.
	//         For example, if the salted ID in hex is c82e0b54_0495ed56, the id[8:12] will be 0495ed56, and the
	//         id[12:20] will be c82e0b54.
	binary.BigEndian.PutUint32(id[8:12], uint32(salted))

	// Same as before, the only difference is the offset.
	copy(id[12:20], randomid[8:16])

	// NOTE 8: We need to shift half of the bits to the right to move the MSB to the LSB.
	//         Converting by uint32 will only take 4 bytes from the LSB.
	binary.BigEndian.PutUint32(id[20:], uint32(salted>>32))

	return id
}

// String returns the string representation of the ShadowID.
func (id ShadowID) String() string {
	text, _ := id.MarshalText()
	return string(text)
}

// MarshalText implements the encoding.TextMarshaler interface, this also covers json.Marshal.
func (id ShadowID) MarshalText() ([]byte, error) {
	enc := make([]byte, hex.EncodedLen(len(id)))
	hex.Encode(enc, id[:])
	return enc, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface, this also covers json.Unmarshal.
func (id *ShadowID) UnmarshalText(text []byte) error {
	p, err := Parse(text)
	if err != nil {
		return fmt.Errorf("shadowid: unmarshal text: %w", err)
	}
	*id = p
	return nil
}

// RandomID returns the random ID (UUIDv4) from the ShadowID.
func (id ShadowID) RandomID() uuid.UUID {
	var uid uuid.UUID
	// Based on the ShadowID anatomy, the random ID consists of the first 8 bytes and the 12th-20th bytes.
	// So, let's copy the first 8 bytes and the 12th-20th bytes to the UUID.
	copy(uid[:8], id[0:8])
	copy(uid[8:], id[12:20])
	return uid
}

// Autoincr returns the autoincr from the ShadowID.
func (id ShadowID) Autoincr() int64 {
	var autoincr [8]byte
	// This is a bit tricky: in NewShadowID, we placed the salted ID's LSB into the 8th-12th and 20th-24th bytes.
	// Since we used BigEndian for both, we need to reverse the order of the salted ID's LSB.
	copy(autoincr[:4], id[20:])
	copy(autoincr[4:], id[8:12])

	// Converts the salted ID to uint64
	salted := binary.BigEndian.Uint64(autoincr[:])

	// We take the random salt from the UUID.
	randomSalt := binary.LittleEndian.Uint64(id[12:20])

	// Apply the same XOR operation as in NewShadowID.
	return int64(salted ^ randomSalt ^ defaultSalt.Load())
}

// Parse parses a ShadowID from a hex string.
func Parse(src []byte) (ShadowID, error) {
	var id ShadowID
	if len(src) != hex.EncodedLen(len(id)) {
		return id, fmt.Errorf("shadowid: unmarshal src: invalid length")
	}

	_, err := hex.Decode(id[:], src)
	if err != nil {
		return id, fmt.Errorf("shadowid: unmarshal src: %w", err)
	}

	return id, nil
}
