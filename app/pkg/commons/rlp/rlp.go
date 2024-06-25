// Package rlp provides functions for encoding and decoding data using the Recursive Length Prefix (RLP) encoding scheme.
package rlp

import (
	"errors"
	"math"
	"math/big"
	"reflect"
)

const (
	singleByteValue                        = 0x80
	singleBytePrefix                       = 0x7f
	longItemPrefix                         = 0xb7
	itemListOfLongerThan55BytesItemsPrefix = 0xbf
	itemListOfShortItemsPrefix             = 0xf7
	shortListPrefix                        = 0xc0 // A short list in rlp would be one between 0 and 55 bytes long

	stringType = "string"
	listType   = "list"
)

var (
	ErrEncodingUnhandledInputType = errors.New("unhandled input type for RLP encoding")
	ErrEncodingInputSizeTooLong   = errors.New("input size is too long for RLP encoding")

	ErrDecodingUnhandledInputType  = errors.New("unhandled input type for RLP decoding")
	ErrDecodingLengthFieldTooShort = errors.New("cannot read data, incomplete length field")
	ErrDecodingNullInput           = errors.New("input cannot be null")

	ErrEncodingNegativeBigInt = errors.New("RLP can not encode negative big.Int")

	ErrDeserializationMismatchedLength          = errors.New("mismatched number of attributes in the output struct and the decoded data")
	ErrDeserializationUnsupportedTypeConversion = errors.New("unsupported type conversion in deserialization")
	ErrDeserializationNonStructOnNestedList     = errors.New("cannot deserialize nested lists on non struct field")
)

// Encode function takes in an item as input, encodes it into its RLP representation and returns the encoded bytes.
// A nil pointer to a struct, slice or array is encoded as an empty RLP list unless the slice or array type is bytes.
// A nil pointer to any other value is encoded as the empty string.
// RLP is a standard encoding for Ethereum data structures
// src: https://ethereum.org/es/developers/docs/data-structures-and-encoding/rlp/
//
// Considerations:
// Encode does not support hexstring encoding as it will interpret it as a normal string, thus if
// you need to encode a hexstring always convert it to []byte
// Example: The hexstring "0x832728" will encode the utf8 unicode values of "0", "x", "8"..., not the actual bytes representation
func Encode(input interface{}) ([]byte, error) {
	switch item := input.(type) {
	case string:
		return encodeString(item)
	case *string:
		return encodeStringPointer(item)
	case []string:
		return encodeStringArray(item)
	case *[]string:
		return encodeStringArrayPointer(item)
	case []interface{}:
		return encodeList(item)
	case *[]interface{}:
		return encodeListPointer(item)
	case []byte:
		return encodeBytes(item)
	case *[]byte:
		return encodeBytesPointer(item)
	case uint:
		return encodeUint(item)
	case *big.Int:
		return encodeBigInt(item)
	default:
		return nil, ErrEncodingUnhandledInputType
	}
}

func encodeBytes(item []byte) ([]byte, error) {
	if len(item) == 1 && item[0] < 0x80 {
		return []byte{item[0]}, nil
	}
	encodedLength, err := encodeLength(len(item), 0x80)
	if err != nil {
		return nil, err
	}
	return append(encodedLength, item...), nil
}

func encodeBytesPointer(item *[]byte) ([]byte, error) {
	if item == nil {
		return []byte{shortListPrefix}, nil
	}
	return encodeBytes(*item)
}

func encodeUint(input interface{}) ([]byte, error) {
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Uint && value.Kind() != reflect.Uint64 {
		return nil, ErrEncodingUnhandledInputType
	}

	return encodeBytes(new(big.Int).SetUint64(value.Uint()).Bytes())
}

func encodeList(list []interface{}) ([]byte, error) {
	var encodedItems []byte

	for _, element := range list {
		switch item := element.(type) {
		case string:
			encodedStr, err := encodeString(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedStr...)
		case *string:
			encodedStrPointer, err := encodeStringPointer(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedStrPointer...)
		case []string:
			encodedStrArray, err := encodeStringArray(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedStrArray...)
		case *[]string:
			encodedStrArrayPointer, err := encodeStringArrayPointer(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedStrArrayPointer...)
		case []interface{}:
			encodedNestedList, err := encodeList(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedNestedList...)
		case *[]interface{}:
			encodedNestedListPointer, err := encodeListPointer(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedNestedListPointer...)
		case []byte:
			encodedBytes, err := encodeBytes(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedBytes...)
		case *[]byte:
			encodedBytesPointer, err := encodeBytesPointer(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedBytesPointer...)
		case uint:
			encodedUint, err := encodeUint(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedUint...)
		case *big.Int:
			encodedBigInt, err := encodeBigInt(item)
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, encodedBigInt...)
		default:
			return nil, ErrEncodingUnhandledInputType
		}
	}

	encodedLength, err := encodeLength(len(encodedItems), shortListPrefix)
	if err != nil {
		return nil, err
	}
	return append(encodedLength, encodedItems...), nil
}

func encodeListPointer(list *[]interface{}) ([]byte, error) {
	if list == nil {
		return []byte{shortListPrefix}, nil
	}
	return encodeList(*list)
}

func encodeStringArray(item []string) ([]byte, error) {
	var encodedBytes []byte
	for _, str := range item {
		encodedStr, err := Encode(str)
		if err != nil {
			return nil, err
		}
		encodedBytes = append(encodedBytes, encodedStr...)
	}
	encodedLength, err := encodeLength(len(encodedBytes), shortListPrefix)
	if err != nil {
		return nil, err
	}
	return append(encodedLength, encodedBytes...), nil
}

func encodeStringArrayPointer(item *[]string) ([]byte, error) {
	if item == nil {
		return []byte{shortListPrefix}, nil
	}
	return encodeStringArray(*item)
}

func encodeString(item string) ([]byte, error) {
	if len(item) == 1 && int(item[0]) < singleByteValue {
		return []byte{item[0]}, nil
	}

	encodedLength, err := encodeLength(len(item), singleByteValue)
	if err != nil {
		return nil, err
	}
	return append(encodedLength, item...), nil
}

func encodeStringPointer(item *string) ([]byte, error) {
	if item == nil {
		return encodeString("")
	}
	return encodeString(*item)
}

func encodeLength(length int, offset int) ([]byte, error) {
	if length < 56 {
		return []byte{byte(length + offset)}, nil
	}
	if length < int(math.Pow(256, 8)) {
		return nil, ErrEncodingInputSizeTooLong
	}
	bl := toBinary(length)
	return append([]byte{byte(len(bl) + offset + 55)}, bl...), nil
}

func encodeBigInt(value *big.Int) ([]byte, error) {
	if value == nil {
		return []byte{singleByteValue}, nil
	}
	if value.Sign() == -1 {
		return nil, ErrEncodingNegativeBigInt
	}

	if value.BitLen() <= 64 {
		return encodeUint(value.Uint64())
	}

	encodedLength, err := encodeLength(len(value.Bytes()), 0x80)
	if err != nil {
		return nil, err
	}
	return append(encodedLength, value.Bytes()...), nil
}

func toBinary(value int) []byte {
	if value == 0 {
		return []byte{}
	}

	bigIntLength := big.NewInt(int64(value))
	quotient, reminder := new(big.Int).DivMod(bigIntLength, big.NewInt(256), new(big.Int))
	return append(toBinary(int(quotient.Int64())), byte(reminder.Int64()))
}

// Decode function takes in an array of binary data. The RLP decoding process is as follows:
// according to the first byte (i.e. prefix) of input data and decoding the data type, the length of the actual data and offset;
// according to the type and offset of data, decode the data correspondingly;
// continue to decode the rest of the input;
// src: https://ethereum.org/es/developers/docs/data-structures-and-encoding/rlp/
func Decode(input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, ErrDecodingNullInput
	}
	offset, dataLength, dataType, err := decodeLength(input)
	if err != nil {
		return nil, err
	}
	switch dataType {
	case stringType:
		return string(input[offset : offset+dataLength]), nil
	case listType:
		return decodeList(input[offset : offset+dataLength])
	default:
		return nil, ErrDecodingUnhandledInputType
	}
}

// DecodeAndDeserialize decodes a byte slice and stores the result in the data type pointed to by the 'output' variable
// this variable must be a non-nil pointer.
//
// Considerations: the struct being used as output has to define its attributes in the order expected to be decoded
// nolint:gocyclo
func DecodeAndDeserialize(b []byte, output interface{}) error {
	if output == nil {
		return errors.New("output cannot be nil ")
	}

	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}

	outputElement := outputValue.Elem()

	// Deserialization in a string
	if outputElement.Kind() == reflect.String {
		decodedData, err := Decode(b)
		if err != nil {
			return err
		}

		if stringValue, isSourceString := decodedData.(string); isSourceString {
			outputElement.SetString(stringValue)
			return nil
		}

		return ErrDeserializationUnsupportedTypeConversion
	}

	// Deserialization in a []string
	if outputElement.Kind() == reflect.Slice && outputElement.Type().Elem().Kind() == reflect.String {
		decodedData, err := Decode(b)
		if err != nil {
			return err
		}

		decodedSlice, isInterface := decodedData.([]interface{})
		if !isInterface {
			return ErrDeserializationUnsupportedTypeConversion
		}

		stringSlice := make([]string, len(decodedSlice))
		for i, item := range decodedSlice {
			if stringValue, isString := item.(string); isString {
				stringSlice[i] = stringValue
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		}

		outputElement.Set(reflect.ValueOf(stringSlice))
		return nil
	}

	// Deserialization in a struct
	if outputElement.Kind() != reflect.Struct {
		return errors.New("output must be a pointer to a struct")
	}

	decodedData, err := Decode(b)
	if err != nil {
		return err
	}

	decodedSlice, isInterface := decodedData.([]interface{})
	if !isInterface {
		if outputElement.NumField() == 1 {
			field := outputElement.Field(0)
			targetType := field.Type()
			if targetType.Kind() == reflect.String {
				field.Set(reflect.ValueOf(decodedData))
				return nil
			}
		}
	}

	numFields := outputElement.NumField()
	if numFields != len(decodedSlice) {
		return ErrDeserializationMismatchedLength
	}

	for i := 0; i < numFields; i++ {
		field := outputElement.Field(i)
		targetType := field.Type()
		sourceValue := reflect.ValueOf(decodedSlice[i])
		if !field.CanSet() {
			continue
		}

		// Precondition: our encoding method always encodes as strings, that's why this code always assumes that the decoded value is a string
		//nolint:exhaustive
		switch targetType.Kind() {
		case reflect.String:
			if stringValue, isSourceString := sourceValue.Interface().(string); isSourceString {
				field.SetString(stringValue)
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		case reflect.Slice:
			if targetType.Elem().Kind() == reflect.String {
				if interfaceSlice, isInterfaceSlice := decodedSlice[i].([]interface{}); isInterfaceSlice {
					stringSlice := make([]string, len(interfaceSlice))
					for j, item := range interfaceSlice {
						if stringValue, isString := item.(string); isString {
							stringSlice[j] = stringValue
						} else {
							return ErrDeserializationUnsupportedTypeConversion
						}
					}
					field.Set(reflect.ValueOf(stringSlice))
				} else {
					return ErrDeserializationUnsupportedTypeConversion
				}
			} else if targetType.Elem().Kind() == reflect.Uint8 {
				if stringValue, isSourceString := sourceValue.Interface().(string); isSourceString {
					field.Set(reflect.ValueOf([]byte(stringValue)))
				} else {
					return ErrDeserializationUnsupportedTypeConversion
				}
			}
		case reflect.Uint:
			if uintInString, isSourceString := sourceValue.Interface().(string); isSourceString {
				// Since our Decode() function decodes in strings, this means that each decoded uint will be converted to its
				// unicode value. We'll need to get the unicode code point of it in order to convert it to uint again.
				unicodeCodePoint := rune(uintInString[0])
				uintValue := uint(unicodeCodePoint)
				field.Set(reflect.ValueOf(uintValue))
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		case reflect.Struct:
			if interfaceSlice, isInterfaceSlice := decodedSlice[i].([]interface{}); isInterfaceSlice {
				if errorDeserializingStruct := nestedDeserialize(interfaceSlice, field.Addr().Interface()); err != nil {
					return errorDeserializingStruct
				}
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		default:
			return ErrDeserializationUnsupportedTypeConversion
		}
	}
	return nil
}

// nestedDeserialize defines the recursive deserialization for interface slices inside a struct
func nestedDeserialize(interfaceSlice []interface{}, output interface{}) error {
	outputValue := reflect.ValueOf(output).Elem()
	if outputValue.Kind() != reflect.Struct {
		return ErrDeserializationNonStructOnNestedList
	}

	numFields := outputValue.NumField()
	if len(interfaceSlice) != numFields {
		return ErrDeserializationMismatchedLength
	}

	for i := 0; i < numFields; i++ {
		field := outputValue.Field(i)
		targetType := field.Type()
		sourceValue := reflect.ValueOf(interfaceSlice[i])
		if !field.CanSet() {
			continue
		}

		//nolint:exhaustive
		switch targetType.Kind() {
		case reflect.String:
			if stringValue, isSourceString := sourceValue.Interface().(string); isSourceString {
				field.SetString(stringValue)
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		case reflect.Slice:
			if targetType.Elem().Kind() == reflect.String {
				stringSlice := make([]string, len(interfaceSlice))
				for j, elem := range interfaceSlice {
					if stringValue, isString := elem.(string); isString {
						stringSlice[j] = stringValue
					} else {
						return ErrDeserializationUnsupportedTypeConversion
					}
				}
				field.Set(reflect.ValueOf(stringSlice))
			} else if targetType.Elem().Kind() == reflect.Uint8 {
				if stringValue, isSourceString := sourceValue.Interface().(string); isSourceString {
					field.Set(reflect.ValueOf([]byte(stringValue)))
				} else {
					return ErrDeserializationUnsupportedTypeConversion
				}
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		case reflect.Uint:
			if uintInString, isSourceString := sourceValue.Interface().(string); isSourceString {
				// Since our Decode() function decodes in strings, this means that each decoded uint will be converted to its
				// unicode value. We'll need to get the unicode code point of it in order to convert it to uint again.
				unicodeCodePoint := rune(uintInString[0])
				uintValue := uint(unicodeCodePoint)
				field.Set(reflect.ValueOf(uintValue))
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		case reflect.Struct:
			if nestedInterfaceSlice, isInterfaceSlice := interfaceSlice[i].([]interface{}); isInterfaceSlice {
				if err := nestedDeserialize(nestedInterfaceSlice, field.Addr().Interface()); err != nil {
					return err
				}
			} else {
				return ErrDeserializationUnsupportedTypeConversion
			}
		default:
			return ErrDeserializationUnsupportedTypeConversion
		}
	}
	return nil
}

// decodeLength takes a byte slice which was part of an RLP-encoded message
// and checks the prefix byte to determine the length and the type of the data
// it returns the position where the data begins, the length, the data type and an error if something fails
func decodeLength(input []byte) (offset int, dataLength int, dataType string, err error) {
	length := len(input)
	if length == 0 {
		return 0, 0, "", ErrDecodingNullInput
	}

	prefix := input[0]
	switch {
	case prefix <= singleBytePrefix:
		return 0, 1, stringType, nil
	case prefix <= longItemPrefix:
		return 1, int(prefix - singleByteValue), stringType, nil
	case prefix <= itemListOfLongerThan55BytesItemsPrefix:
		lengthOfLengthField := int(prefix - longItemPrefix)
		if length <= 1+lengthOfLengthField {
			return 0, 0, "", ErrDecodingLengthFieldTooShort
		}
		return 1 + lengthOfLengthField, toInt(input[1 : 1+lengthOfLengthField]), stringType, nil
	case prefix <= itemListOfShortItemsPrefix:
		return 1, int(prefix - shortListPrefix), listType, nil
	default:
		lengthOfLengthField := int(prefix - itemListOfShortItemsPrefix)
		if length <= 1+lengthOfLengthField {
			err = ErrDecodingLengthFieldTooShort
			return
		}
		return 1 + lengthOfLengthField, toInt(input[1 : 1+lengthOfLengthField]), listType, nil
	}
}

func toInt(b []byte) int {
	length := len(b)
	if length == 0 {
		return 0
	} else if length == 1 {
		return int(b[0])
	}
	return int(b[length-1]) + toInt(b[:length-1])*256
}

func decodeList(input []byte) ([]interface{}, error) {
	var output []interface{}
	i := 0
	for i < len(input) {
		offset, dataLength, dataType, err := decodeLength(input[i:])
		if err != nil {
			return nil, err
		}
		switch dataType {
		case stringType:
			output = append(output, string(input[i+offset:i+offset+dataLength]))
		case listType:
			nestedList, decodeListErr := decodeList(input[i+offset : i+offset+dataLength])
			if decodeListErr != nil {
				return nil, decodeListErr
			}
			//nolint:asasalint
			output = append(output, nestedList)
		default:
			return nil, ErrDecodingUnhandledInputType
		}
		i += offset + dataLength
	}
	return output, nil
}
