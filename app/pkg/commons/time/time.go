// Package time defines time based on the domain logic.
// It ensures that time within the application is aligned with the requirements and constraints dictated by the domain.
package time

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Timestamp defines an instant in time with milliseconds precision.
type Timestamp struct {
	time int64
}

// ToInt64 formats Timestamp to int64
func (t Timestamp) ToInt64() int64 {
	return t.time
}

// TimestampFromInt64 creates a new Timestamp from the provided int64
func TimestampFromInt64(src int64) Timestamp {
	return Timestamp{time: src}
}

// Now creates a new Timestamp with the current instant of time in milliseconds since January 1, 1970, UTC
func Now() Timestamp {
	return Timestamp{time: time.Now().UnixMilli()}
}

// UnixNano returns the timestamp as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC.
func (t Timestamp) UnixNano() int64 {
	return t.time * 1000000
}

// String returns a string representation of Timestamp
func (t Timestamp) String() string {
	return fmt.Sprintf("%d", t.time)
}

// Add creates a new Timestamp with the result of the sum for the timestamp and the amount provided as parameter.
func (t Timestamp) Add(amount int64) Timestamp {
	return Timestamp{time: t.time + amount}
}

// Sub creates a new Timestamp with the result of the rest for the timestamp and the amount provided as parameter.
func (t Timestamp) Sub(amount int64) Timestamp {
	return Timestamp{time: t.time - amount}
}

// MarshalJSON returns Timestamp JSON-marshalled output or error if it fails
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON unmarshalls JSON data into Timestamp or returns an error if it fails
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var timeString string
	err := json.Unmarshal(data, &timeString)
	if err != nil {
		// Cannot use Timestamp as destination because it forces an infinite loop
		// https://biscuit.ninja/posts/go-avoid-an-infitine-loop-with-custom-json-unmarshallers/
		var tmpJSON map[string]interface{}
		err = json.Unmarshal(data, &tmpJSON)
		if err != nil {
			return err
		}
		if _, ok := tmpJSON["time"]; ok {
			timeString = tmpJSON["time"].(string)
		}
	}

	var timeInt int64
	if len(timeString) > 0 {
		timeInt, err = strconv.ParseInt(timeString, 10, 64)
		if err != nil {
			return err
		}
	}

	t.time = timeInt
	return nil
}
