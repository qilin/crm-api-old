package eventbus

type Bus struct {
	in         string
	out        string
	marshaller Marshaller
}

func NewBus(in, out string, marshaller Marshaller) *Bus {
	return &Bus{
		in:         in,
		out:        out,
		marshaller: marshaller,
	}
}
