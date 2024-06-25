package rlp_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger-labs/signare/app/pkg/commons/rlp"
)

var loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipisicing elit"
var encodedLoremIPsum = []byte{0xb8, 0x38, 0x4c, 0x6f, 0x72, 0x65, 0x6d, 0x20, 0x69, 0x70, 0x73, 0x75, 0x6d, 0x20, 0x64, 0x6f, 0x6c, 0x6f, 0x72, 0x20, 0x73, 0x69, 0x74, 0x20, 0x61, 0x6d, 0x65, 0x74, 0x2c, 0x20, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x63, 0x74, 0x65, 0x74, 0x75, 0x72, 0x20, 0x61, 0x64, 0x69, 0x70, 0x69, 0x73, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x20, 0x65, 0x6c, 0x69, 0x74}

var zeroInteger = 0

var emptyString = ""
var encodedEmptyString = []byte{0x80}

var testInteger = 1024

var testInteger2 = 15

var testString = "dog"
var encodedTestString = []byte{0x83, 'd', 'o', 'g'}

var testList = []string{"cat", "dog"}
var encodedTestList = []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}
var decodedTestList = []interface{}{"cat", "dog"}

var testBytes = []byte{0x01, 0x02, 0x03, 0x04}
var encodedTestBytes = []byte{0x84, 0x01, 0x02, 0x03, 0x04}
var decodedTestBytes = "\x01\x02\x03\x04"

var encodedEmptyList = []byte{0xc0}

var nestedList interface{} = []interface{}{
	"hello world",
	"1",
	"",
	[]interface{}{"1", "2", "3"},
	[]interface{}{"dog", "cat", "frog"},
}
var encodedNestedList = "0xe08b68656c6c6f20776f726c643180c3313233cd83646f67836361748466726f67"

var blockNumberAndBlockHash = []string{
	"0xb6cf",
	"0x2d40765edb5cda24af3f26ec207c78f6a754b68b56657a5a653a5cb0af7206a8",
}

var encodedBlockHeaderData = "f9027ea016fc3e618906c161b951738ea6392375d409d6640880e6c25ce7380d40ff97dfa01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794ca31306798b41bc81c43094a1e0462890ce7a673a0e51931b4ee95ac77368978c168ee0a173124a84ed70d96f2a352beefd111a7fda0217bcffbfd58413206859c981d8d09c6f2e45f03c93ed5996a1fbb26a5a8a974a048e29eeab42547c238445f65e4e45e10ba0fabe1c72f63357b05c8a49520adbeb90100000000000000080000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000181b28405f5e10083017c59846512bc6db884f882a00000000000000000000000000000000000000000000000000000000000000000d594ca31306798b41bc81c43094a1e0462890ce7a673808400000000f843b8412dda242ab37c9c7507e3c77c9398a93544fd75ef63392045da12863cd29b267b42a6599dfd9985358e6f89090341e6e7d89b8fb85ed0c80bb74137f332f75abf00a063746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365880000000000000000"

var encodedBlockHeaderData2 = "f9040301830650d3b9010000000000000000000004000000000000000800400002000000400000000000000000000000000000000000000000000200000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000001000000000000000000000000f902f8f8b99457f1defacaafaa664bbb0f4f3347b786d874a25ce1a0f0e25c63981e9a617375c8244c8ac144e4e520acc2cb08e2d4c3781b1a02c067b8800000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004038383632303030316331653239326566626463636433316638396436633832326464353263656138643234336635663434373263623434383763373962316639f9023a94c13141df25ac03eb6e249ab4319f61d2be3d4254e1a07a752c4d100a96be23f60f072fed1f33da2890fbe45eb15af52750faa2772a92b90200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000037bcb3cac66f4d859a4ef77dcd97eec146bbc425000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001788903901f000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000083066376138323034000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b4854474247423030555344000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b48545553555330305553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002ac3a5e96e4a09b6613273aca6649f453f94257120000000000000000"

var encodedEthProofInHexString = "0xf9040301830650d3b9010000000000000000000004000000000000000800400002000000400000000000000000000000000000000000000000000200000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000001000000000000000000000000f902f8f8b99457f1defacaafaa664bbb0f4f3347b786d874a25ce1a0f0e25c63981e9a617375c8244c8ac144e4e520acc2cb08e2d4c3781b1a02c067b8800000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004038383632303030316331653239326566626463636433316638396436633832326464353263656138643234336635663434373263623434383763373962316639f9023a94c13141df25ac03eb6e249ab4319f61d2be3d4254e1a07a752c4d100a96be23f60f072fed1f33da2890fbe45eb15af52750faa2772a92b90200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000037bcb3cac66f4d859a4ef77dcd97eec146bbc425000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001788903901f000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000083066376138323034000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b4854474247423030555344000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b48545553555330305553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002ac3a5e96e4a09b6613273aca6649f453f94257120000000000000000"

var encodedData = []byte{0xc8, 0x83, 'f', 'o', 'o', 0x83, 'b', 'a', 'r', 0x83, 'c', 'a', 't'}

var encodedBlockHeaderExtraData = "f90144a00000000000000000000000000000000000000000000000000000000000000000f854948823902f9a09b9a590f6b74c3b5b3d870e744fed94d16e8d7fe64fae7691be7e29454902a3394d2ff294e9837db2ef07d50a522edbea0a6009cd433f2b8594f9849f03096a6a59c9383dca4109675b83dadb84c080f8c9b841d9f48527a0652b8e4837f451f131a810ac653233d904d275cc574bb5dba4eefc0c7eb208d93c5d8cfa9f253aef4543aff424c08facd9fd5d097bec45ed9f3f6301b841cdd3fe1fe3013e30b28da08c66b0974b4852266538aaad6747df4ced76fc51ce317b63aae8845cb4bc907edce4397c8939bf26c0ca47c83f28aff588d82004b901b841219c44d7270cba0aa70b9152b6e0ec0f1b9eef4ee4b339bda833109c5c9ee5727facae1787cb585a8be4b31d37a7b6029fa45497edb767cb4e6dbf9a6b61861d01"

var status = uint(1)
var cumulativeGasUsed = uint(413907)
var logsBloom = "0x00000000000000000004000000000000000800400002000000400000000000000000000000000000000000000000000200000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000001000000000000000000000000"
var logs = []interface{}{
	[]interface{}{
		hexStringToBytes("0x57f1defacaafaa664bbb0f4f3347b786d874a25c"),
		[]interface{}{hexStringToBytes("0xf0e25c63981e9a617375c8244c8ac144e4e520acc2cb08e2d4c3781b1a02c067")},
		hexStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004038383632303030316331653239326566626463636433316638396436633832326464353263656138643234336635663434373263623434383763373962316639"),
	},
	[]interface{}{
		hexStringToBytes("0xc13141df25ac03eb6e249ab4319f61d2be3d4254"),
		[]interface{}{hexStringToBytes("0x7a752c4d100a96be23f60f072fed1f33da2890fbe45eb15af52750faa2772a92")},
		hexStringToBytes("0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000037bcb3cac66f4d859a4ef77dcd97eec146bbc425000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001788903901f000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000083066376138323034000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b4854474247423030555344000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b48545553555330305553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002ac3a5e96e4a09b6613273aca6649f453f94257120000000000000000"),
	},
}

var encodedBlockNumberAndBlockHash = "0xf84b86307862366366b842307832643430373635656462356364613234616633663236656332303763373866366137353462363862353636353761356136353361356362306166373230366138"

var receipt interface{} = []interface{}{
	status,
	cumulativeGasUsed,
	hexStringToBytes(logsBloom),
	logs,
}

type foo struct {
	Foo string
	Bar string
}

type nestedListType struct {
	Value  string
	Value2 string
	Value3 string
	List   []string
	List2  []string
}

type loremIPsumStruct struct {
	Text string
}

type blockHeader struct {
	Value   string
	Value2  string
	Value3  string
	Value4  string
	Value5  string
	Value6  string
	Value7  string
	Value8  string
	Value9  string
	Value10 string
	Value11 string
	Value12 string
	Value13 string
	Value14 string
	Value15 string
}

type ethProofReceipt struct {
	Val1 string
	Val2 string
	Val3 string
	Val4 EthProofNested
}

type EthProofNested struct {
	Val1 EthProofNested2
	Val2 EthProofNested2
}

type WrongEthProofNested struct {
	Val1 EthProofNested2
	Val2 EthProofNested4
}

type EthProofNested4 struct {
}

type EthProofNested2 struct {
	Val1 string
	Val2 EthProofNested3
	Val3 string
}

type EthProofNested3 struct {
	Val1 string
}

type BlockHeaderExtraData struct {
	Value1 string
	Value2 []string
	Value3 []string
	Value4 string
	Value5 []string
}

type DeserializeTestStruct struct {
	Bytes  []byte
	String string
	Uint   uint
}

func Test_RLP_Encode_CorrectExecution(t *testing.T) {
	output, err := rlp.Encode(testString)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, encodedTestString, output)

	output, err = rlp.Encode(emptyString)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, encodedEmptyString, output)

	output, err = rlp.Encode(testList)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, output, encodedTestList)

	output, err = rlp.Encode([]string{})
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, output, encodedEmptyList)

	output, err = rlp.Encode(nestedList)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, bytesToHexString(output), encodedNestedList)

	output, err = rlp.Encode(loremIpsum)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, output, encodedLoremIPsum)

	output, err = rlp.Encode(blockNumberAndBlockHash)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, bytesToHexString(output), encodedBlockNumberAndBlockHash)

	output, err = rlp.Encode(testBytes)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, encodedTestBytes, output)
}

func Test_RLP_Encode_Error_UnsupportedType(t *testing.T) {
	output, err := rlp.Encode(errors.New("dummy error"))
	require.Nil(t, output)
	require.Error(t, err)
	require.Equal(t, err, rlp.ErrEncodingUnhandledInputType)

	output, err = rlp.Encode(-12)
	require.Nil(t, output)
	require.Error(t, err)
	require.Equal(t, err, rlp.ErrEncodingUnhandledInputType)

	output, err = rlp.Encode(zeroInteger)
	require.Error(t, err)
	require.Nil(t, output)
	require.Equal(t, err, rlp.ErrEncodingUnhandledInputType)

	output, err = rlp.Encode(testInteger)
	require.Error(t, err)
	require.Nil(t, output)
	require.Equal(t, err, rlp.ErrEncodingUnhandledInputType)

	output, err = rlp.Encode(testInteger2)
	require.Error(t, err)
	require.Nil(t, output)
	require.Equal(t, err, rlp.ErrEncodingUnhandledInputType)
}

func Test_RLP_Decode_CorrectExecution(t *testing.T) {
	output, err := rlp.Decode(encodedTestList)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, decodedTestList, output)

	output, err = rlp.Decode(encodedTestBytes)
	require.Nil(t, err)
	require.NotNil(t, output)

	output, err = rlp.Decode(encodedLoremIPsum)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, loremIpsum, output)

	output, err = rlp.Decode(encodedEmptyString)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, "", output)

	output, err = rlp.Decode(encodedTestString)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, testString, output)

	output, err = rlp.Decode(encodedTestBytes)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, decodedTestBytes, output)

	output, err = rlp.Decode(encodedTestBytes)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, decodedTestBytes, output)

	encodeNestedList, err := rlp.Encode(nestedList)
	require.Nil(t, err)
	require.NotNil(t, encodedNestedList)
	output, err = rlp.Decode(encodeNestedList)
	require.Nil(t, err)
	require.NotNil(t, output)
	require.Equal(t, nestedList, output)

	receiptProofInBytes, decodeErr := hex.DecodeString(encodedBlockHeaderData)
	require.Nil(t, decodeErr)
	output, err = rlp.Decode(receiptProofInBytes)
	require.Nil(t, err)
	require.NotNil(t, output)

	receiptProofInBytes, decodeErr = hex.DecodeString(encodedBlockHeaderData2)
	require.Nil(t, decodeErr)
	output, err = rlp.Decode(receiptProofInBytes)
	require.Nil(t, err)
	require.NotNil(t, output)
}

func Test_RLP_Error_NilValue(t *testing.T) {
	output, err := rlp.Decode([]byte{})
	require.Error(t, err)
	require.Nil(t, output)
	require.Equal(t, err, rlp.ErrDecodingNullInput)
}

func Test_RLP_DecodeAndDeserialize_CorrectExecution(t *testing.T) {
	expectedStruct := foo{
		Foo: "foo",
		Bar: "bar",
	}
	var decodedStruct foo
	err := rlp.DecodeAndDeserialize(encodedData, &decodedStruct)
	require.Nil(t, err)
	require.Equal(t, decodedStruct, expectedStruct)

	var nestedListStruct nestedListType
	encodedNestedList, err := rlp.Encode(nestedList)
	require.Nil(t, err)
	require.NotNil(t, encodedNestedList)
	err = rlp.DecodeAndDeserialize(encodedNestedList, &nestedListStruct)
	require.Nil(t, err)

	var lorem loremIPsumStruct
	err = rlp.DecodeAndDeserialize(encodedLoremIPsum, &lorem)
	require.Nil(t, err)

	var loremString string
	err = rlp.DecodeAndDeserialize(encodedLoremIPsum, &loremString)
	require.Nil(t, err)

	var blockHeaderStruct blockHeader
	encodedBlockHeaderDataInBytes, decodeErr := hex.DecodeString(encodedBlockHeaderData)
	require.Nil(t, decodeErr)
	err = rlp.DecodeAndDeserialize(encodedBlockHeaderDataInBytes, &blockHeaderStruct)
	require.Nil(t, err)

	var ethProofReceiptOne ethProofReceipt
	encodedEthProofReceipt, decodeErr := rlp.Encode(receipt)
	require.Nil(t, decodeErr)
	err = rlp.DecodeAndDeserialize(encodedEthProofReceipt, &ethProofReceiptOne)
	require.Nil(t, err)

	var blockHeaderExtraDataStruct BlockHeaderExtraData
	encodedBlockHeaderExtraDataInBytes, decodeErr := hex.DecodeString(encodedBlockHeaderExtraData)
	require.Nil(t, decodeErr)
	err = rlp.DecodeAndDeserialize(encodedBlockHeaderExtraDataInBytes, &blockHeaderExtraDataStruct)
	require.Nil(t, err)

	var deserializeTestStruct interface{} = []interface{}{
		[]byte{65, 66, 67},
		"hello",
		uint(42),
	}
	expectedDeserialization := DeserializeTestStruct{
		Bytes:  []byte{65, 66, 67},
		String: "hello",
		Uint:   uint(42),
	}
	var deserializedStruct DeserializeTestStruct
	encodedDeserializeTestStruct, decodeErr := rlp.Encode(deserializeTestStruct)
	require.Nil(t, decodeErr)
	err = rlp.DecodeAndDeserialize(encodedDeserializeTestStruct, &deserializedStruct)
	require.Nil(t, err)
	require.Equal(t, deserializedStruct, expectedDeserialization)
}

func Test_RLP_DecodeAndDeserialize_Error(t *testing.T) {
	var wrongEthProof WrongEthProofNested
	encodedEthProofReceipt, decodeErr := rlp.Encode(receipt)
	require.Nil(t, decodeErr)
	err := rlp.DecodeAndDeserialize(encodedEthProofReceipt, &wrongEthProof)
	require.Error(t, err)
	require.Equal(t, err, rlp.ErrDeserializationMismatchedLength)
}

func Test_RLP_EncodeAndDecode_Tx_Receipt_CorrectExecution(t *testing.T) {
	encodedEthProof, err := rlp.Encode(receipt)
	require.Nil(t, err)
	require.NotNil(t, encodedEthProof)
	require.Equal(t, encodedEthProofInHexString, bytesToHexString(encodedEthProof))

	decodedEthProof, err := rlp.Decode(encodedEthProof)
	require.Nil(t, err)
	require.NotNil(t, decodedEthProof)
}

func hexStringToBytes(input string) []byte {
	if len(input) == 0 {
		panic("empty string")
	}
	if !has0xPrefix(input) {
		panic("missing '0x' prefix")
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		panic(err.Error())
	}
	return b
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func bytesToHexString(bytes []byte) string {
	hexValue := "0x"
	for _, b := range bytes {
		hexValue += fmt.Sprintf("%02x", b)
	}

	return hexValue
}
