package common

import "encoding/json"

type Wrapper interface {
	Wrap(payload Payloader, attempt int) (Event, error)
	UnWrap(data []byte) (Event, error)
}

type Marshaller interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
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
	msg, err := w.marshaller.Marshal(payload.Payload())
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

func (w jsonWrapper) UnWrap(data []byte) (Event, error) {
	var evt Event
	err := w.marshaller.Unmarshal(data, &evt)
	return evt, err
}

// Marshaller

type jsonMarshaller struct{}

func NewJSONMarshaller() Marshaller {
	return &jsonMarshaller{}
}

func (j *jsonMarshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonMarshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
