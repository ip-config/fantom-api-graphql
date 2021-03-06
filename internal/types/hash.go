package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"reflect"
)

// HashLength is the expected length of the hash type.
const HashLength = 32

// Hash represents the 32 byte hash of arbitrary data.
// It's inspired by go-ethereum Hash type, just doesn't follow some of the naming conventions.
type Hash [HashLength]byte

// hashT holds the reflected type of Hash.
var hashT = reflect.TypeOf(Hash{})

// Hex converts a hash to a hex string.
func (h Hash) Hex() string {
	return hexutil.Encode(h[:])
}

// String implements the stringer interface and is used also by the logger to attach Hash information.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, the byte slice to be formatted as is, without going through the stringer.
func (h Hash) Format(s fmt.State, c rune) {
	if _, err := fmt.Fprintf(s, "%"+string(c), h[:]); err != nil {
		return
	}
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax from JSON input.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of Hash.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// Scan implements Scanner interface for reading Hash from database/sql.
func (h *Hash) Scan(src interface{}) error {
	// get raw bytes
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can not scan %T into Hash", src)
	}

	// do we get the right size
	if len(srcB) != HashLength {
		return fmt.Errorf("can not scan []byte of len %d into Hash, expected %d", len(srcB), HashLength)
	}

	copy(h[:], srcB)
	return nil
}

// Value implements Valuer interface for storing Hash in a database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// UnmarshalGraphQL unmarshal the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Hash", input)
	}
	return err
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type, in this case it implements Hash.
func (Hash) ImplementsGraphQLType(name string) bool {
	return name == "Hash"
}
