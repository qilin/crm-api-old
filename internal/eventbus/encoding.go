package eventbus

import "encoding/json"

type Wrapper interface {
	Wrap(payload Payloader, attempt int) (Event, error)
}

type Marshaller interface {
	Marshall(v interface{}) ([]byte, error)
	UnMarshall(data []byte, v interface{}) error
}

// Wrapper

type jsonWrapper struct {
	marshaller Marshaller
}

func NewJsonWrapper(marshaller Marshaller) Wrapper {
	return &jsonWrapper{
		marshaller: marshaller,
	}
}

func (w jsonWrapper) Wrap(payload Payloader, attempt int) (Event, error) {
	msg, err := w.marshaller.Marshall(payload.Payload())
	if err != nil {
		return Event{}, err
	}
	return Event{
		Attempt: attempt,
		Name:    payload.Name(),
		Payload: msg,
		Version: EventVersion,
	}, nil
}

// Marshaller

type jsonMarshaller struct{}

func NewJSONMarshaller() Marshaller {
	return &jsonMarshaller{}
}

func (j *jsonMarshaller) Marshall(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonMarshaller) UnMarshall(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
