package encoding

import (
	"strings"
)

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Name() string
}

var registeredEncoder = make(map[string]Encoder)

func RegisterEncoder(encoder Encoder) {
	if encoder == nil {
		panic("cannot register a nil Encoder")
	}
	if encoder.Name() == "" {
		panic("cannot register Encoder with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(encoder.Name())
	registeredEncoder[contentSubtype] = encoder
}

func GetEncoder(encoderType string) Encoder {
	return registeredEncoder[encoderType]
}
