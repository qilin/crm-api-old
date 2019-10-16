package eventbus_test

import (
	"testing"
	"time"

	"github.com/qilin/crm-api/internal/eventbus"
	"github.com/stretchr/testify/assert"
)

func TestJsonMarshaller(t *testing.T) {
	a := assert.New(t)

	m := eventbus.NewJSONMarshaller()

	var inv2 eventbus.Invite
	inv := eventbus.Invite{
		TenantId:   1,
		UserId:     2,
		Expiration: time.Time{},
		Email:      "john@example.com",
		FirstName:  "John",
		LastName:   "Doe",
		Accepted:   false,
	}

	b, err := m.Marshall(inv)
	a.NoError(err)

	err = m.UnMarshall(b, &inv2)
	a.NoError(err)

	a.EqualValues(inv, inv2)
}

//func TestJsonWrapper(t *testing.T) {
//	a := assert.New(t)
//
//	attempt := 1
//
//	// 1. build Payload := Invite{ .... }
//	// 2. wrap(Payload{}): Event{ ..., Invite{}, ...}
//	// 3. Encode: marshal(Event{...})
//	// 4. Transmit
//	//
//	// 1. Receive
//	// 2. Decode Event := unmarshal([]byte{...})
//	// 3. unwrap(Event{}, &Invite)
//
//	w := eventbus.NewJsonWrapper(eventbus.NewJSONMarshaller())
//
//	invite := eventbus.Invite{
//		TenantId:   10,
//		UserId:     200,
//		Expiration: time.Now(),
//		Email:      "john@example.com",
//		FirstName:  "John",
//		LastName:   "Doe",
//		Accepted:   false,
//	}
//
//	evt, err := w.Wrap(invite, 1)
//	a.NoError(err)
//	a.NotEmpty(b)
//
//	var invite2
//	err = w.UnWrap(evt, invite2)
//	a.NoError(err)
//
//	a.Equal(evt.Name, invite.Name())
//	a.Equal(evt.Attempt, attempt)
//	a.Equal(evt.Version, eventbus.EventVersion)
//}

/*

rcv := []byte{}
event := unmarshal(rcv)
invite := unmarshal(event.payload)

...invite processing...

switch {
  case ack:
	...
  case resend:
    event.attempt := event.attempt + 1
}

*/
