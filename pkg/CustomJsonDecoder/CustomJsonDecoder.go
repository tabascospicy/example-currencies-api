package CustomJsonDecoder

import (
	"encoding/json"
	"io"
)

type DecodeError struct {
	Err     error
	Context string
}

// custom Error for decoding data
func (d *DecodeError) Error() string {
	return "Error decoding data: " + d.Err.Error() + " with context " + d.Context
}

// function to wrap the decoding of a json response process
func DecodeJson(response io.Reader, data interface{}, name string) error {

	val, err := io.ReadAll(response)

	if err != nil {
		return &DecodeError{Err: err, Context: name}
	}
	
	errJson := json.Unmarshal(val, &data)

	if errJson != nil {
		return &DecodeError{Err: err, Context: name}
	}

	return nil
}
