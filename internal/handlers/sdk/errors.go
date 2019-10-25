package sdk

const (
	// Custom backend communications errors [100]1xxx
	errBackendTimeout = 1001000

	// Common Bad Request [400]00xx
	errBadRequest = 400
	// Auth Bad Request [400]11xx
	errAuthRequestURLEmpty = 4001101
	// Order Bad Request [400]12xx
	errOrderRequestDataEmpty = 4001201

	// Internal server error
	errInternalServerError = 500
)

var statusText = map[int]string{
	// 400
	errBadRequest:            "request malformed",
	errAuthRequestURLEmpty:   "url can not be empty",
	errOrderRequestDataEmpty: "data can not be empty",
	// 500
	errInternalServerError: "internal server error",
}

func StatusText(code int) string {
	return statusText[code]
}
