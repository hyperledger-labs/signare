package entities_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

func TestNewInt256FromString(t *testing.T) {
	_, err := entities.NewInt256FromString("junk")
	if err == nil {
		t.Error("NewInt256FromString should've failed")
	}
}

func TestNewInt256FromInt(t *testing.T) {
	in := int64(50000)
	out := big.NewInt(in)
	value := entities.NewInt256FromInt(in)
	if value.BigInt().Cmp(out) != 0 {
		t.Error("NewInt256FromInt failed to construct the correct value")
	}
}

func TestIsInt256(t *testing.T) {
	value := entities.IsInt256("junk")
	if value {
		t.Error("IsInt256 should've returned false for string")
	}
	value = entities.IsInt256("0X12345678")
	if !value {
		t.Error("IsInt256 should've returned true for valid integer")
	}
	value = entities.IsInt256("1234")
	if !value {
		t.Error("IsInt256 should've returned true for valid integer")
	}
}

func TestZeroInt256(t *testing.T) {
	result := entities.ZeroInt256()
	if result == nil {
		t.Error("ZeroInt256 failed")
	}
}

func TestNewHexInt256FromStringFailure(t *testing.T) {
	_, err := entities.NewHexInt256FromString("junk")
	if err == nil {
		t.Error("NewHexInt256FromString should've failed for non numbers")
	}
	_, err = entities.NewHexInt256FromString("0x115792089237316195423570985008687907853269984665640564039457584007913129639936")
	if err == nil {
		t.Error("NewHexInt256FromString should've failed for too large a number")
	}
}

func TestNewHexInt256FromStringSuccess(t *testing.T) {
	_, err := entities.NewHexInt256FromString("0x12345678")
	if err != nil {
		t.Error("NewHexInt256FromString should've succeeded for valid hex integer")
	}
}

func TestIsHexInt256True(t *testing.T) {
	ok := entities.IsHexInt256("0x427364885")
	if !ok {
		t.Error("NewHexInt256FromString should've returned true for a valid hex integer")
	}
}

func TestIsHexInt256False(t *testing.T) {
	ok := entities.IsHexInt256("zzz")
	if ok {
		t.Error("IsHexInt256 should've returned false for invalid hex integer")
	}
}

func TestHexInt256UnmarshalText(t *testing.T) {
	tests := []struct {
		input string
		num   *big.Int
		ok    bool
	}{
		{"", big.NewInt(0), true},
		{"0", nil, false},
		{"0x0", big.NewInt(0), true},
		{"12345678", nil, false},
		{"0x12345678", big.NewInt(0x12345678), true},
		{"0X12345678", big.NewInt(0x12345678), true},
		{"0123456789", nil, false},
		{"00", nil, false},
		{"0x00", nil, false},
		{"0x012345678abc", nil, false},
		{"abcdef", nil, false},
		{"0xgg", nil, false},
		{"0x115792089237316195423570985008687907853269984665640564039457584007913129639936", nil, false},
	}
	for _, test := range tests {
		var num entities.HexInt256
		err := num.UnmarshalText([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("HexInt256.UnmarshalText(%q) -> (err == nil) == %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if test.num != nil && num.BigInt().Cmp(test.num) != 0 {
			t.Errorf("HexInt256.UnmarshalText(%q) -> %d, want %d", test.input, num.BigInt(), test.num)
		}
	}
}

func TestHexInt256MarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"0x0", "0x0", true},
		{"0xFFFFFFFFFFFF", "0xffffffffffff", true},
		{"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", true},
	}
	for _, test := range tests {
		num, _ := entities.NewHexInt256FromString(test.input)
		m, err := num.MarshalText()
		if (err == nil) != test.ok {
			t.Errorf("HexInt256(%q).MarshalText() -> (err == nil) == %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if test.expected != "" && !bytes.Equal([]byte(test.expected), m) {
			t.Errorf("HexInt256(%q).MarshalText() -> %s, want %s", test.input, m, test.expected)
		}
	}
}

func TestHexInt256UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		num   *big.Int
		ok    bool
	}{
		{"\"\"", big.NewInt(0), true},
		{"\"0\"", nil, false},
		{"\"0x0\"", big.NewInt(0), true},
		{"\"12345678\"", nil, false},
		{"\"0x12345678\"", big.NewInt(0x12345678), true},
		{"\"0X12345678\"", big.NewInt(0x12345678), true},
		{"\"0x115792089237316195423570985008687907853269984665640564039457584007913129639936\"", nil, false},
		{"0x123", nil, false},
	}
	for _, test := range tests {
		var num entities.HexInt256
		err := num.UnmarshalJSON([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("HexInt256(%q).UnmarshalJSON -> (err == nil) == %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if test.num != nil && num.BigInt().Cmp(test.num) != 0 {
			t.Errorf("HexInt256(%q).UnmarshalJSON -> %d, want %d", test.input, num.BigInt(), test.num)
		}
	}
}

func TestHexInt256MarshalJSON(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
		ok       bool
	}{
		{int64(119690537), "\"0x7225529\"", true},
		{int64(0), "\"0x0\"", true},
	}
	for _, test := range tests {
		num := entities.NewHexInt256(big.NewInt(test.input))
		m, err := json.Marshal(num)
		if (err == nil) != test.ok {
			t.Errorf("HexInt256(%d).MarshalJSON -> (err == nil) == %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if test.expected != "" && !bytes.Equal(m, []byte(test.expected)) {
			t.Errorf("HexInt256(%q).MarshalJSON -> %s, want %s", test.input, m, test.expected)
		}
	}
}

func TestMaxBigInt(t *testing.T) {
	a := big.NewInt(8)
	b := big.NewInt(5)
	max1 := entities.MaxBigInt(a, b)
	if max1 != a {
		t.Errorf("Expected %d got %d", a, max1)
	}
	max2 := entities.MaxBigInt(b, a)
	if max2 != a {
		t.Errorf("Expected %d got %d", a, max2)
	}
}

func TestMinBigInt(t *testing.T) {
	a := big.NewInt(5)
	b := big.NewInt(3)
	min1 := entities.MinBigInt(a, b)
	if min1 != b {
		t.Errorf("Expected %d got %d", b, min1)
	}
	min2 := entities.MinBigInt(b, a)
	if min2 != b {
		t.Errorf("Expected %d got %d", b, min2)
	}
}

func TestInt256FirstBitSet(t *testing.T) {
	tests := []struct {
		num *entities.Int256
		ix  int
	}{
		{entities.NewInt256(big.NewInt(0)), 0},
		{entities.NewInt256(big.NewInt(1)), 0},
		{entities.NewInt256(big.NewInt(2)), 1},
		{entities.NewInt256(big.NewInt(0x100)), 8},
	}
	for _, test := range tests {
		if ix := test.num.FirstBitSet(); ix != test.ix {
			t.Errorf("Int256(%q).FirstBitSet() = %d, want %d", test.num.BigInt(), ix, test.ix)
		}
	}
}

func TestInt256PaddedBytes(t *testing.T) {
	tests := []struct {
		num    *big.Int
		n      int
		result []byte
	}{
		{num: big.NewInt(0), n: 4, result: []byte{0, 0, 0, 0}},
		{num: big.NewInt(1), n: 4, result: []byte{0, 0, 0, 1}},
		{num: big.NewInt(512), n: 4, result: []byte{0, 0, 2, 0}},
		{num: entities.PowBigInt(2, 32), n: 4, result: []byte{1, 0, 0, 0, 0}},
	}
	for _, test := range tests {
		if result := entities.NewInt256(test.num).PaddedBytes(test.n); !bytes.Equal(result, test.result) {
			t.Errorf("Int256(%q).PaddedBytes(%d) = %v, want %v", test.num, test.n, result, test.result)
		}
	}
}

func TestInt256ReadBits(t *testing.T) {
	check := func(input string) {
		want, _ := hex.DecodeString(input)
		integer, _ := new(big.Int).SetString(input, 16)
		buf := make([]byte, len(want))
		entities.NewInt256(integer).ReadBits(buf)
		if !bytes.Equal(buf, want) {
			t.Errorf("Int256(%q).ReadBits(..) = %x, want: %x", integer, buf, want)
		}
	}
	check("000000000000000000000000000000000000000000000000000000FEFCF3F8F0")
	check("0000000000012345000000000000000000000000000000000000FEFCF3F8F0")
	check("18F8F8F1000111000110011100222004330052300000000000000000FEFCF3F8F0")
}

func TestInt256U256(t *testing.T) {
	tests := []struct{ x, y *big.Int }{
		{x: big.NewInt(0), y: big.NewInt(0)},
		{x: big.NewInt(1), y: big.NewInt(1)},
		{x: entities.PowBigInt(2, 255), y: entities.PowBigInt(2, 255)},
		{x: entities.PowBigInt(2, 256), y: big.NewInt(0)},
		{x: new(big.Int).Add(entities.PowBigInt(2, 256), big.NewInt(1)), y: big.NewInt(1)},
		{x: big.NewInt(-1), y: new(big.Int).Sub(entities.PowBigInt(2, 256), big.NewInt(1))},
		{x: big.NewInt(-2), y: new(big.Int).Sub(entities.PowBigInt(2, 256), big.NewInt(2))},
		{x: entities.PowBigInt(2, -255), y: big.NewInt(1)},
	}
	for _, test := range tests {
		if y := entities.NewInt256(new(big.Int).Set(test.x)).U256(); y.Cmp(test.y) != 0 {
			t.Errorf("Int256(%q).U256() = %x, want %x", test.x, y, test.y)
		}
	}
}

func TestInt256U256Bytes(t *testing.T) {
	ubytes := make([]byte, 32)
	ubytes[31] = 1
	unsigned := entities.NewInt256(big.NewInt(1)).U256Bytes()
	if !bytes.Equal(unsigned, ubytes) {
		t.Errorf("Int256(1).U256Bytes() = %x, want %x", unsigned, ubytes)
	}
}

func TestInt256BigEndianByteAt(t *testing.T) {
	tests := []struct {
		x   string
		y   int
		exp byte
	}{
		{"00", 0, 0x00},
		{"01", 1, 0x00},
		{"00", 1, 0x00},
		{"01", 0, 0x01},
		{"0000000000000000000000000000000000000000000000000000000000102030", 0, 0x30},
		{"0000000000000000000000000000000000000000000000000000000000102030", 1, 0x20},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 31, 0xAB},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 32, 0x00},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 500, 0x00},
	}
	for _, test := range tests {
		hexBytes, _ := entities.HexToBytes(test.x)
		v := new(big.Int).SetBytes(hexBytes)
		actual := entities.NewInt256(v).BigEndianByteAt(test.y)
		if actual != test.exp {
			t.Fatalf("Int256(%s).BigEndianByteAt(%v) = %v, want %v", test.x, test.y, actual, test.exp)
		}
	}
}

func TestInt256LittleEndianByteAt(t *testing.T) {
	tests := []struct {
		x   string
		y   int
		exp byte
	}{
		{"00", 0, 0x00},
		{"01", 1, 0x00},
		{"00", 1, 0x00},
		{"01", 0, 0x00},
		{"0000000000000000000000000000000000000000000000000000000000102030", 0, 0x00},
		{"0000000000000000000000000000000000000000000000000000000000102030", 1, 0x00},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 31, 0x00},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 32, 0x00},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 0, 0xAB},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 1, 0xCD},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", 0, 0x00},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", 1, 0xCD},
		{"0000000000000000000000000000000000000000000000000000000000102030", 31, 0x30},
		{"0000000000000000000000000000000000000000000000000000000000102030", 30, 0x20},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 32, 0x0},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 31, 0xFF},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0xFFFF, 0x0},
	}
	for _, test := range tests {
		hexBytes, _ := entities.HexToBytes(test.x)
		v := entities.NewInt256(new(big.Int).SetBytes(hexBytes))
		actual := v.ByteAt(32, test.y)
		if actual != test.exp {
			t.Fatalf("Int256(%s).LittleEndianByteAt(%v) = %v, want %v", test.x, test.y, test.exp, actual)
		}

	}
}

func TestInt256S256(t *testing.T) {
	tests := []struct{ x, y *big.Int }{
		{x: big.NewInt(0), y: big.NewInt(0)},
		{x: big.NewInt(1), y: big.NewInt(1)},
		{x: big.NewInt(2), y: big.NewInt(2)},
		{
			x: new(big.Int).Sub(entities.PowBigInt(2, 255), big.NewInt(1)),
			y: new(big.Int).Sub(entities.PowBigInt(2, 255), big.NewInt(1)),
		},
		{
			x: entities.PowBigInt(2, 255),
			y: new(big.Int).Neg(entities.PowBigInt(2, 255)),
		},
		{
			x: new(big.Int).Sub(entities.PowBigInt(2, 256), big.NewInt(1)),
			y: big.NewInt(-1),
		},
		{
			x: new(big.Int).Sub(entities.PowBigInt(2, 256), big.NewInt(2)),
			y: big.NewInt(-2),
		},
	}
	for _, test := range tests {
		if y := entities.NewInt256(test.x).S256(); y.Cmp(test.y) != 0 {
			t.Errorf("Int256(%x).S256() = %x, want %x", test.x, y, test.y)
		}
	}
}

type operation byte

const (
	sub operation = iota
	add
	mul
)

func TestNewUInt64FromStringValid(t *testing.T) {
	if v, _ := entities.NewUInt64FromString("1"); v.Uint64() != 1 {
		t.Errorf(`NewUInt64FromString("1") = %d, want 1`, v)
	}
	if v, _ := entities.NewUInt64FromString("0x1"); v.Uint64() != 1 {
		t.Errorf(`NewUInt64FromString("0x1") = %d, want 1`, v)
	}
	if v, _ := entities.NewUInt64FromString(""); v.Uint64() != 0 {
		t.Errorf(`NewUInt64FromString("") = %d, want 1`, v)
	}
}

func TestNewUInt64Valid(t *testing.T) {
	if v := entities.NewUInt64(uint64(3)); v.Uint64() != 3 {
		t.Errorf(`NewUInt64("3") = %d, want 3`, v)
	}
}

func TestUInt64String(t *testing.T) {
	if str := entities.NewUInt64(uint64(3)).String(); str != "3" {
		t.Errorf(`entities.NewUInt64(uint64(3)).String() = %s, want 3`, str)
	}
}

func TestUInt64Overflow(t *testing.T) {
	for i, test := range []struct {
		x        uint64
		y        uint64
		overflow bool
		op       operation
	}{
		{entities.MaxUInt64, 1, true, add},
		{entities.MaxUInt64 - 1, 1, false, add},
		{0, 1, true, sub},
		{0, 0, false, sub},
		{0, 0, false, mul},
		{10, 10, false, mul},
		{entities.MaxUInt64, 2, true, mul},
		{entities.MaxUInt64, 1, false, mul},
	} {
		var overflows bool
		switch test.op {
		case sub:
			_, overflows = entities.UInt64(test.x).Sub(test.y)
		case add:
			_, overflows = entities.UInt64(test.x).Add(test.y)
		case mul:
			_, overflows = entities.UInt64(test.x).Mul(test.y)
		}
		if test.overflow != overflows {
			t.Errorf("Uint64 overflow test %d failed. Got %v, want %v", i, overflows, test.overflow)
		}
	}
}

func TestHexUInt64UnmarshalText(t *testing.T) {
	tests := []struct {
		input string
		num   uint64
		ok    bool
	}{
		{"", 0, true},
		{"0", 0, false},
		{"0x0", 0, true},
		{"12345678", 12345678, false},
		{"0x12345678", 0x12345678, true},
		{"0X12345678", 0x12345678, true},
		{"0123456789", 123456789, false},
		{"0x00", 0, false},
		{"0x012345678abc", 0x12345678abc, false},
		{"abcdef", 0, false},
		{"0xgg", 0, false},
		{"18446744073709551617", 0, false},
	}
	for _, test := range tests {
		var num entities.HexUInt64
		err := num.UnmarshalText([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("HexUInt64.UnmarshalText(%q) -> (err == nil) = %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		if err == nil && num.Uint64() != test.num {
			t.Errorf("HexUInt64.UnmarshalText(%q) -> %d, want %d", test.input, num, test.num)
		}
	}
}

func TestHexBytesUnmarshalText(t *testing.T) {
	tests := []struct {
		input string
		val   string
		ok    bool
	}{
		{"", "0x", true},
		{"0", "0", false},
		{"0x0", "0x0", false},
		{"12345678", "12345678", false},
		{"0x12345678", "0x12345678", true},
		{"0X12345678", "0x12345678", true},
		{"0123456789", "123456789", false},
		{"0x00", "0x00", true},
		{"0x012345678abc", "0x012345678abc", true},
		{"abcdef", "", false},
		{"0xgg", "", false},
		{"18446744073709551617", "", false},
	}
	for _, test := range tests {
		var val entities.HexBytes
		err := val.UnmarshalText([]byte(test.input))
		if (err == nil) != test.ok {
			t.Errorf("HexBytes.UnmarshalText(%q) -> (err == nil) = %t, want %t", test.input, err == nil, test.ok)
			continue
		}
		valStr := val.String()
		if err == nil && valStr != test.val {
			t.Errorf("HexBytes.UnmarshalText(%q) -> %s, want %s", test.input, val, test.val)
		}
	}
}

func TestNewHexUInt64FromStringValid(t *testing.T) {
	if v, _ := entities.NewHexUInt64FromString("0xAF2C"); v.Uint64() != 44844 {
		t.Errorf(`NewHexUInt64FromString("0xAF2C") = %d, want 44844`, v)
	}
}

func TestNewHexUInt64FromStringInvalid(t *testing.T) {
	_, err := entities.NewHexUInt64FromString("ggg")
	if err == nil {
		t.Error("NewHexUInt64FromString should've failed")
	}
}

func TestNewHexBytesFromStringValid(t *testing.T) {
	v, _ := entities.NewHexBytesFromString("0x1234")
	if v.String() != "0x1234" {
		t.Errorf(`NewHexBytesFromString("0x1234") = %s, want 0x1234`, v)
	}
}

func TestNewHexBytesFromStringInvalidEmpty(t *testing.T) {
	_, err := entities.NewHexBytesFromString("")
	if err == nil {
		t.Error("NewHexBytesFromString should've failed")
	}
}

func TestNewHexBytesFromStringInvalidHexPrefix(t *testing.T) {
	_, err := entities.NewHexBytesFromString("123")
	if err == nil {
		t.Error("NewHexBytesFromString should've failed")
	}
}

func TestNewHexBytes32FromStringInvalidEmpty(t *testing.T) {
	_, err := entities.NewHexBytes32FromString("")
	if err == nil {
		t.Error("NewHexBytes32FromString should've failed")
	}
}

func TestNewHexBytes32FromStringInvalidHexPrefix(t *testing.T) {
	_, err := entities.NewHexBytes32FromString("123")
	if err == nil {
		t.Error("NewHexBytes32FromString should've failed")
	}
}

func TestCopyBytes(t *testing.T) {
	input := []byte{1, 2, 3, 4}

	v := entities.CopyBytes(input)
	if !bytes.Equal(v, []byte{1, 2, 3, 4}) {
		t.Fatal("CopyBytes not equal after copy")
	}
	v[0] = 99
	if bytes.Equal(v, input) {
		t.Fatal("CopyBytes result is not a copy")
	}
}

func TestLeftPadBytes(t *testing.T) {
	val := []byte{1, 2, 3, 4}
	padded := []byte{0, 0, 0, 0, 1, 2, 3, 4}

	if r := entities.LeftPadBytes(val, 8); !bytes.Equal(r, padded) {
		t.Fatalf("LeftPadBytes(%v, 8) == %v", val, r)
	}
	if r := entities.LeftPadBytes(val, 2); !bytes.Equal(r, val) {
		t.Fatalf("LeftPadBytes(%v, 2) == %v", val, r)
	}
}

func TestRightPadBytes(t *testing.T) {
	val := []byte{1, 2, 3, 4}
	padded := []byte{1, 2, 3, 4, 0, 0, 0, 0}

	if r := entities.RightPadBytes(val, 8); !bytes.Equal(r, padded) {
		t.Fatalf("RightPadBytes(%v, 8) == %v", val, r)
	}
	if r := entities.RightPadBytes(val, 2); !bytes.Equal(r, val) {
		t.Fatalf("RightPadBytes(%v, 2) == %v", val, r)
	}
}

func TestFromHex(t *testing.T) {
	input := "0x01"
	expected := []byte{1}
	result, _ := entities.FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected %x got %x", expected, result)
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"", true},
		{"0", false},
		{"00", true},
		{"a9e67e", true},
		{"A9E67E", true},
		{"0xa9e67e", false},
		{"a9e67e001", false},
		{"0xHELLO_KITTY_@#$^&*", false},
	}
	for _, test := range tests {
		if ok := entities.IsHex(test.input); ok != test.ok {
			t.Errorf("IsHex(%q) = %v, want %v", test.input, ok, test.ok)
		}
	}
}

func TestFromHexOddLength(t *testing.T) {
	input := "0x1"
	expected := []byte{1}
	result, _ := entities.FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("FromHex(%q) = %x, want %x", input, result, expected)
	}
}

func TestNoPrefixShortHexOddLength(t *testing.T) {
	input := "1"
	expected := []byte{1}
	result, _ := entities.FromHex(input)
	if !bytes.Equal(expected, result) {
		t.Errorf("FromHex(%q) = %x, want %x", input, result, expected)
	}
}

func TestTrimRightZeroes(t *testing.T) {
	tests := []struct {
		arr []byte
		exp []byte
	}{
		{fromHex("0x00ffff00ff0000"), fromHex("0x00ffff00ff")},
		{fromHex("0x00000000000000"), []byte{}},
		{fromHex("0xff"), fromHex("0xff")},
		{[]byte{}, []byte{}},
		{fromHex("0x00ffffffffffff"), fromHex("0x00ffffffffffff")},
	}
	for i, test := range tests {
		got := entities.TrimRightZeroes(test.arr)
		if !bytes.Equal(got, test.exp) {
			t.Errorf("TrimRightZeroes test %d failed, got %x, want %x", i, got, test.exp)
		}
	}
}

func fromHex(s string) []byte {
	result, _ := entities.FromHex(s)
	return result
}
