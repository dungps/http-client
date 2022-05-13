package json

import (
	"encoding/json"
	"github.com/dungps/http-client/internal/encoding"
	"reflect"
)

const Name = "json"

func init() {
	encoding.RegisterEncoder(jsonEncoder{})
}

type jsonEncoder struct{}

func (jsonEncoder) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	default:
		return json.Marshal(v)
	}
}

func (jsonEncoder) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		return json.Unmarshal(data, m)
	}
}

func (jsonEncoder) Name() string {
	return Name
}
