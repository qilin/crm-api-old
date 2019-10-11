package errors

type ClientErr struct {
	err error
}

func (c *ClientErr) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

type PanicErr struct {
	err error
}

func (c *PanicErr) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

type AccessDeniedErr struct {
	err error
}

func (c *AccessDeniedErr) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

func WrapPanicErr(err error) error {
	return &PanicErr{err: err}
}

func WrapClientErr(err error) error {
	return &ClientErr{err: err}
}

func WrapAccessDeniedErr(err error) error {
	return &AccessDeniedErr{err: err}
}
