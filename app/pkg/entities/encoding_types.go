// Package entities defines custom data types to be used across the whole application.
package entities

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"reflect"
	"strconv"

	"golang.org/x/crypto/sha3"

	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

const (
	wordBits  = 32 << (uint64(^big.Word(0)) >> 63) // Number of bits in a big.Word
	wordBytes = wordBits / 8                       // Number of bytes in a big.Word
)

var (
	tt255   = PowBigInt(2, 255)
	tt256   = PowBigInt(2, 256)
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1))
)

const MaxUInt64 = 1<<64 - 1

// Int256 as a 256-bit integer
type Int256 struct {
	big.Int
}

// NewInt256FromInt creates a new 256-bit integer from a 64 bit input integer
func NewInt256FromInt(x int64) *Int256 {
	b := big.NewInt(x)
	h := Int256{*b}
	return &h
}

// NewInt256FromString parses the input string s as a 256-bit integer in decimal or hexadecimal syntax. Leading zeros are allowed. The empty string parses as zero.
func NewInt256FromString(s string) (*Int256, error) {
	if len(s) == 0 {
		return ZeroInt256(), nil
	}
	var bigInt *big.Int
	var ok bool
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		bigInt, ok = new(big.Int).SetString(s[2:], 16)
	} else {
		bigInt, ok = new(big.Int).SetString(s, 10)
	}
	if !ok {
		return nil, signererrors.InvalidArgument().WithMessage("the provided string [%s] is not a valid number", s)
	}
	if bigInt.BitLen() > 256 {
		return nil, signererrors.InvalidArgument().WithMessage("the provided number [%s] length is more than 256-bits", s)
	}
	return NewInt256(bigInt), nil
}

// NewInt256 returns a 256-bit integer from the input big integer
func NewInt256(i *big.Int) *Int256 {
	if i == nil {
		return nil
	}
	h := Int256{*i}
	return &h
}

// ZeroInt256 returns a zero valued 256-bit integer
func ZeroInt256() *Int256 {
	b := new(big.Int)
	h := Int256{*b}
	return &h
}

// IsInt256 parses the input string s as a 256-bit integer. It returns true if valid and false otherwise.
func IsInt256(s string) bool {
	_, err := NewInt256FromString(s)
	return err == nil
}

// BigInt converts b to a big.BigInt.
func (b *Int256) BigInt() *big.Int {
	return &b.Int
}

// U256 encodes to the 256-bit two's complement number.
func (b *Int256) U256() *big.Int {
	i := b.BigInt()
	return i.And(i, tt256m1)
}

// U256Bytes converts into a padded 256-bit EVM number.
func (b *Int256) U256Bytes() []byte {
	return NewInt256(b.U256()).PaddedBytes(32)
}

// S256 interprets b as a two's complement number. S256(0) = 0, S256(1) = 1, S256(2**255) = -2**255, S256(2**256-1) = -1.
func (b *Int256) S256() *big.Int {
	i := b.BigInt()
	if i.Cmp(tt255) < 0 {
		return i
	}
	return new(big.Int).Sub(i, tt256)
}

// FirstBitSet returns the index of the first 1 bit in v, counting from LSB.
func (b *Int256) FirstBitSet() int {
	i := b.BigInt()
	for k := 0; k < i.BitLen(); k++ {
		if i.Bit(k) > 0 {
			return k
		}
	}
	return i.BitLen()
}

// PaddedBytes encodes a big integer as a big-endian byte slice. The length of the slice is at least n bytes.
func (b *Int256) PaddedBytes(n int) []byte {
	if b.BigInt().BitLen()/8 >= n {
		return b.BigInt().Bytes()
	}
	ret := make([]byte, n)
	b.ReadBits(ret)
	return ret
}

// BigEndianByteAt returns the byte at position n, in big-endian encoding, so that n==0 returns the least significant byte.
func (b *Int256) BigEndianByteAt(n int) byte {
	words := b.BigInt().Bits()
	// Check word-bucket the byte will reside in
	l := n / wordBytes
	if l >= len(words) {
		return byte(0)
	}
	word := words[l]
	// Offset of the byte
	shift := 8 * uint(n%wordBytes)
	return byte(word >> shift)
}

// ByteAt returns the byte at position n, with the supplied pad length in little-endian encoding.
func (b *Int256) ByteAt(padLength, n int) byte {
	if n >= padLength {
		return byte(0)
	}
	return b.BigEndianByteAt(padLength - 1 - n)
}

// ReadBits encodes the absolute value of bigint as big-endian bytes.
func (b *Int256) ReadBits(buf []byte) {
	l := len(buf)
	for _, d := range b.BigInt().Bits() {
		for j := 0; j < wordBytes && l > 0; j++ {
			l--
			buf[l] = byte(d)
			d >>= 8
		}
	}
}

// MarshalJSON implements the json.Marshaler.
func (b Int256) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

// HexInt256 as an hex encoded 256-bit integer.
type HexInt256 struct {
	Int256
}

const invalidNibble = ^uint64(0)

var (
	typeHexBytes  = reflect.TypeOf(HexBytes(nil))
	typeHexInt256 = reflect.TypeOf((*HexInt256)(nil))
	typeHexUInt64 = reflect.TypeOf(HexUInt64{UInt64(0)})
)

var (
	errorHexEncodingEmptyString   = signererrors.InvalidArgument().WithMessage("hex: empty string")
	errorHexEncodingSyntax        = signererrors.InvalidArgument().WithMessage("hex: invalid hex string")
	errorHexEncodingMissingPrefix = signererrors.InvalidArgument().WithMessage("hex: string without 0x prefix")
	errorHexEncodingOddLength     = signererrors.InvalidArgument().WithMessage("hex: string of odd length")
	errorHexEncodingEmptyNumber   = signererrors.InvalidArgument().WithMessage("hex: string 0x")
	errorHexEncodingLeadingZero   = signererrors.InvalidArgument().WithMessage("hex: number with leading zero digits")
	errorHexEncodingUint64Range   = signererrors.InvalidArgument().WithMessage("hex: number > 64 bits")
	errorHexEncodingInt256Range   = signererrors.InvalidArgument().WithMessage("hex: number > 256-bits")
)

// Number of nibbles (groups of 4 bytes) in a big word.
var bigWordNibbles int

func init() {
	// Compute the number of nibbles required for big.Word on this architecture.
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1: // 64-bit architectures.
		bigWordNibbles = 16
	case 2: // 32-bit architectures.
		bigWordNibbles = 8
	default:
		panic("Invalid word size")
	}
}

// NewHexInt256FromString decodes a hex string and returns a new 256-bit integer. Numbers larger than 256-bits are not allowed.
func NewHexInt256FromString(input string) (*HexInt256, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return nil, err
	}
	if len(raw) > 64 {
		return nil, errorHexEncodingInt256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == invalidNibble {
				return nil, errorHexEncodingSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	decoded := new(big.Int).SetBits(words)
	h := HexInt256{Int256{*decoded}}
	return &h, nil
}

// NewHexInt256 returns a hex 256-bit integer from the input big integer
func NewHexInt256(i *big.Int) *HexInt256 {
	if i == nil {
		return nil
	}
	h := HexInt256{Int256{*i}}
	return &h
}

// MarshalJSON implements the json.Marshaler.
func (b *HexInt256) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Encode())
}

// MarshalText implements encoding.TextMarshaler
func (b HexInt256) MarshalText() ([]byte, error) {
	return []byte(b.Encode()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *HexInt256) UnmarshalJSON(input []byte) error {
	if !checkString(input) {
		return errorHexEncodingNonString(typeHexInt256)
	}
	return errorHexEncodingInvalidType(b.UnmarshalText(input[1:len(input)-1]), typeHexInt256)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *HexInt256) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 64 {
		return errorHexEncodingInt256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == invalidNibble {
				return errorHexEncodingSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	var decodedValue big.Int
	decodedValue.SetBits(words)
	*b = HexInt256{Int256{decodedValue}}
	return nil
}

// Encode encodes bigint as a hex string with 0x prefix. The sign of the integer is ignored.
func (b *HexInt256) Encode() string {
	i := b.BigInt()
	bitLength := i.BitLen()
	if bitLength == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", i)
}

// String implements fmt.Stringer.
func (b *HexInt256) String() string {
	return b.Encode()
}

// IsHexInt256 checks if the input string is an hex encoded 256-bit integer
func IsHexInt256(input string) bool {
	_, err := NewHexInt256FromString(input)
	return err == nil
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", errorHexEncodingEmptyString
	}
	if !checkPrefixFromString(input) {
		return "", errorHexEncodingMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return "", errorHexEncodingEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return "", errorHexEncodingLeadingZero
	}
	return input, nil
}

func checkString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func checkPrefixForBytes(input []byte) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func checkPrefixFromString(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func checkText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if checkPrefixForBytes(input) {
		input = input[2:]
	} else if wantPrefix {
		return nil, errorHexEncodingMissingPrefix
	}
	if len(input)%2 != 0 {
		return nil, errorHexEncodingOddLength
	}
	return input, nil
}

func checkNumberText(input []byte) (raw []byte, err error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if !checkPrefixForBytes(input) {
		return nil, errorHexEncodingMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return nil, errorHexEncodingEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return nil, errorHexEncodingLeadingZero
	}
	return input, nil
}

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return invalidNibble
	}
}

func errorHexEncoding(err error) error {
	var convertErr *strconv.NumError
	if errors.As(err, &convertErr) {
		if errors.Is(convertErr.Err, strconv.ErrRange) {
			return errorHexEncodingUint64Range
		}
		if errors.Is(convertErr.Err, strconv.ErrSyntax) {
			return errorHexEncodingSyntax
		}
	}
	var hexErr hex.InvalidByteError
	if errors.As(err, &hexErr) {
		return errorHexEncodingSyntax
	}
	if errors.Is(err, hex.ErrLength) {
		return errorHexEncodingOddLength
	}
	return err
}

func errorHexEncodingInvalidType(err error, typ reflect.Type) error {
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
	}
	return err
}

func errorHexEncodingNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}

// UInt64 as an unsigned 64 bit integer
type UInt64 uint64

// NewUInt64FromString parses s as an integer in decimal or hexadecimal syntax. Leading zeros are accepted. The empty string parses as zero.
func NewUInt64FromString(s string) (UInt64, error) {
	if s == "" {
		return UInt64(0), nil
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		parsedHex, parseHexErr := strconv.ParseUint(s[2:], 16, 64)
		if parseHexErr != nil {
			return 0, parseHexErr
		}
		return UInt64(parsedHex), nil
	}
	parsedDec, parsedDecErr := strconv.ParseUint(s, 10, 64)
	if parsedDecErr != nil {
		return 0, parsedDecErr
	}
	return UInt64(parsedDec), nil
}

// NewUInt64 wraps the given input as a UInt64.
func NewUInt64(input uint64) UInt64 {
	return UInt64(input)
}

func (ui UInt64) Sub(y uint64) (uint64, bool) {
	x := uint64(ui)
	diff, borrowOut := bits.Sub64(x, y, 0)
	return diff, borrowOut != 0
}

func (ui UInt64) Add(y uint64) (uint64, bool) {
	x := uint64(ui)
	sum, carryOut := bits.Add64(x, y, 0)
	return sum, carryOut != 0
}

func (ui UInt64) Mul(y uint64) (uint64, bool) {
	x := uint64(ui)
	hi, lo := bits.Mul64(x, y)
	return lo, hi != 0
}

// Uint64 converts b to a uint64.
func (ui UInt64) Uint64() uint64 {
	return uint64(ui)
}

// String implements fmt.Stringer
func (ui UInt64) String() string {
	return fmt.Sprintf("%d", ui.Uint64())
}

// HexUInt64 as an hex encoded 64 bit unsigned integer
type HexUInt64 struct {
	UInt64
}

// NewHexUInt64 wraps the given input as a HexUInt64.
func NewHexUInt64(input uint64) HexUInt64 {
	return HexUInt64{UInt64(input)}
}

// NewHexUInt64FromString decodes a hex string with 0x prefix and returns it as a HexUInt64.
func NewHexUInt64FromString(input string) (HexUInt64, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return NewHexUInt64(0), err
	}
	decoded, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		err = errorHexEncoding(err)
	}
	return NewHexUInt64(decoded), err
}

// MarshalJSON implements the json.Marshaler.
func (b *HexUInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.encode())
}

// MarshalText implements encoding.TextMarshaler.
func (b HexUInt64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 10)
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, b.Uint64(), 16)
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *HexUInt64) UnmarshalJSON(input []byte) error {
	if !checkString(input) {
		return errorHexEncodingNonString(typeHexUInt64)
	}
	return errorHexEncodingInvalidType(b.UnmarshalText(input[1:len(input)-1]), typeHexUInt64)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *HexUInt64) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 16 {
		return errorHexEncodingUint64Range
	}
	var dec uint64
	for _, byt := range raw {
		nib := decodeNibble(byt)
		if nib == invalidNibble {
			return errorHexEncodingSyntax
		}
		dec *= 16
		dec += nib
	}
	*b = NewHexUInt64(dec)
	return nil
}

// encode encodes a 64 bit unsigned integer as a hex string with 0x prefix.
func (b HexUInt64) encode() string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, b.Uint64(), 16))
}

// String returns the hex encoding of b.
func (b HexUInt64) String() string {
	return b.encode()
}

// HexBytes marshals/unmarshals with 0x prefix. The empty slice marshals as 0x.
type HexBytes []byte

// NewHexBytesFromString decodes a hex string with 0x prefix.
func NewHexBytesFromString(input string) (HexBytes, error) {
	if len(input) == 0 {
		return nil, errorHexEncodingEmptyString
	}
	if !checkPrefixFromString(input) {
		return nil, errorHexEncodingMissingPrefix
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = errorHexEncoding(err)
	}
	return b, err
}

// NewHexBytes wraps the given input as a HexBytes.
func NewHexBytes(input []byte) *HexBytes {
	if input == nil {
		return nil
	}
	result := HexBytes(input)
	return &result
}

// MarshalText implements encoding.TextMarshaler.
func (b HexBytes) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2+2)
	copy(result, `0x`)
	hex.Encode(result[2:], b)
	return result, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *HexBytes) UnmarshalText(input []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		err = errorHexEncoding(err)
	} else {
		*b = dec
	}
	return err
}

// MarshalJSON implements the json.Marshaler.
func (b HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Encode())
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *HexBytes) UnmarshalJSON(input []byte) error {
	if !checkString(input) {
		return errorHexEncodingNonString(typeHexBytes)
	}
	return errorHexEncodingInvalidType(b.UnmarshalText(input[1:len(input)-1]), typeHexBytes)
}

// Encode encodes b as a hex string with 0x prefix.
func (b HexBytes) Encode() string {
	i := []byte(b)
	enc := make([]byte, len(i)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], i)
	return string(enc)
}

// String implements Stringer.
func (b HexBytes) String() string {
	return b.Encode()
}

// Bytes returns the underlying slice of bytes of b.
func (b *HexBytes) Bytes() []byte {
	return *b
}

// HexBytes32 holds the internal representation for HexBytes32
type HexBytes32 [32]byte

// NewHexBytes32FromString decodes a hex string with 0x prefix.
func NewHexBytes32FromString(input string) (HexBytes32, error) {
	result := HexBytes32{}
	if len(input) == 0 {
		return result, errorHexEncodingEmptyString
	}
	if !checkPrefixFromString(input) {
		return result, errorHexEncodingMissingPrefix
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = errorHexEncoding(err)
		return result, err
	}
	result.FromBytes(b)
	return result, err
}

// FromBytes transform a slice of bytes into a HexBytes32.
func (b *HexBytes32) FromBytes(src []byte) {
	dst := HexBytes32{}
	copy(dst[:], src)
	*b = dst
}

// FromString transform a string into a HexBytes32.
func (b *HexBytes32) FromString(src string) {
	b.FromBytes([]byte(src))
}

// String implements fmt.Stringer.
func (b HexBytes32) String() string {
	n := bytes.IndexByte(b[:], 0)
	if n != -1 {
		return string(b[:n])
	}
	return string(b[:])
}

// Encode encodes b as a hex string with 0x prefix.
func (b HexBytes32) Encode() string {
	return hex.EncodeToString(b[:])
}

// MarshalJSON implements the json.Marshaler.
func (b HexBytes32) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Encode())
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *HexBytes32) UnmarshalJSON(input []byte) error {
	if !checkString(input) {
		return errorHexEncodingNonString(typeHexBytes)
	}
	return errorHexEncodingInvalidType(b.UnmarshalText(input[1:len(input)-1]), typeHexBytes)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *HexBytes32) UnmarshalText(input []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	dec := make([]byte, 32)
	if _, err = hex.Decode(dec, raw); err != nil {
		err = errorHexEncoding(err)
	} else {
		b.FromBytes(dec)
	}
	return err
}

// Bytes exports HexBytes32 as a slice of bytes.
func (b HexBytes32) Bytes() []byte {
	return b[:]
}

// FromHex returns the bytes represented by the hexadecimal string s which may be prefixed with "0x".
func FromHex(s string) ([]byte, error) {
	if HasHexPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return HexToBytes(s)
}

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// HasHexPrefix validates str begins with '0x' or '0X'.
func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// IsHexCharacter returns bool of c being a valid hexadecimal.
func IsHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// IsHex validates whether each byte is valid hexadecimal string.
func IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !IsHexCharacter(c) {
			return false
		}
	}
	return true
}

// HexToBytes returns the bytes represented by the hexadecimal string str.
func HexToBytes(str string) ([]byte, error) {
	if has0xPrefix(str) {
		return hex.DecodeString(str[2:])
	}
	return hex.DecodeString(str)
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// RightPadBytes zero-pads slice to the right up to length l.
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded, slice)
	return padded
}

// LeftPadBytes zero-pads slice to the left up to length l.
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)
	return padded
}

// TrimRightZeroes returns a subslice of s without trailing zeroes
func TrimRightZeroes(s []byte) []byte {
	idx := len(s)
	for ; idx > 0; idx-- {
		if s[idx-1] != 0 {
			break
		}
	}
	return s[:idx]
}

// PowBigInt returns a ** b as a 256-bit integer.
func PowBigInt(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// MaxBigInt returns the larger of x or y.
func MaxBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// MinBigInt returns the smaller of x or y.
func MinBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}

// HashKeccak256 returns the keccak256 hash of the input data
func HashKeccak256(data []byte) (*HexBytes, error) {
	d := sha3.NewLegacyKeccak256()
	for i := range data {
		_, err := d.Write(data[i : i+1])
		if err != nil {
			return nil, err
		}
	}
	return NewHexBytes(d.Sum(nil)), nil
}
