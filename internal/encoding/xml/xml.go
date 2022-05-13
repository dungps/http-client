package xml

import (
	"encoding/xml"
	"github.com/dungps/http-client/internal/encoding"
)

const Name = "xml"

func init() {
	encoding.RegisterEncoder(xmlEncoder{})
}

type xmlEncoder struct{}

func (xmlEncoder) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (xmlEncoder) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (xmlEncoder) Name() string {
	return Name
}
