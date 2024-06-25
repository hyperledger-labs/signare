package time_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/time"

	"github.com/stretchr/testify/require"
)

func TestTimestamp_Add(t *testing.T) {
	initialTimeStamp := int64(1000)

	timestampFromInt64 := time.TimestampFromInt64(initialTimeStamp)

	addedTimeStamp := int64(800)
	time2 := timestampFromInt64.Add(addedTimeStamp)

	require.Equal(t, initialTimeStamp+addedTimeStamp, time2.ToInt64())
}

func TestTimestamp_Sub(t *testing.T) {
	initialTimeStamp := int64(2000)

	timestampFromInt64 := time.TimestampFromInt64(initialTimeStamp)

	subtractedTimeStamp := int64(800)
	time2 := timestampFromInt64.Sub(subtractedTimeStamp)

	require.Equal(t, initialTimeStamp-subtractedTimeStamp, time2.ToInt64())
}

func TestTimestamp_MarhsalAndUnmarshalJSON_Success(t *testing.T) {
	expectedTimestamp := time.TimestampFromInt64(1200)

	bytesTimestamp, err := expectedTimestamp.MarshalJSON()
	require.NoError(t, err)

	newTimestamp := time.Timestamp{}
	err = newTimestamp.UnmarshalJSON(bytesTimestamp)
	require.NoError(t, err)

	require.Equal(t, expectedTimestamp, newTimestamp)
}

func TestTimestamp_MarhsalAndUnmarshalJSON_Failure(t *testing.T) {
	invalidString := []byte(`[1, 2, 3]`)
	invalidInt64 := []byte(`"this is a string"`)

	newTimestamp := time.Timestamp{}

	err1 := newTimestamp.UnmarshalJSON(invalidString)
	require.Error(t, err1)

	err2 := newTimestamp.UnmarshalJSON(invalidInt64)
	require.Error(t, err2)

	require.NotEqual(t, err1, err2)
}

func TestTimestamp_UnixNano_Success(t *testing.T) {
	newTimestamp := time.Now()

	require.Equal(t, newTimestamp.UnixNano(), newTimestamp.ToInt64()*1000000)
}
