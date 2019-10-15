package eventbus

import "encoding/json"

type Wrapper interface {
	Wrap(typ string, payload []byte) Event
	UnWrap(event []byte) Event
}

type Marshaller interface {
	Marshall(interface{}) ([]byte, error)
	UnMarshaller(source []byte, destination interface{}) error
}

type jsonMarshaller struct{}

func NewJSONMarshaller() Marshaller {
	return &jsonMarshaller{}
}

func (j *jsonMarshaller) Marshall(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonMarshaller) UnMarshaller(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
