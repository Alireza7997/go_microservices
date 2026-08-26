package utils

import (
	"encoding/json"
	"errors"
	"io"
)

func ParseBody(body io.ReadCloser, output interface{}) {
	bytes, err1 := io.ReadAll(body)
	err2 := json.Unmarshal(bytes, output)
	if err1 != nil {
		panic(errors.New("BodyNotProvidedProperly"))
	} else if err2 != nil {
		panic(errors.New("BodyNotProvidedProperly"))
	}
}
