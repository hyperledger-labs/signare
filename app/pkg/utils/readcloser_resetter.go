package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ReadAndResetCloser(reader *io.ReadCloser, pointer any) error {
	body, err := io.ReadAll(*reader)
	if err != nil {
		return fmt.Errorf("error reading body [%w]", err)
	}

	var buf bytes.Buffer
	_, err = buf.Write(body)
	if err != nil {
		return fmt.Errorf("error writing body [%w]", err)
	}

	*reader = io.NopCloser(&buf)
	err = json.Unmarshal(body, pointer)
	if err != nil {
		return fmt.Errorf("error unmarshaling body [%w]", err)
	}

	return nil
}
