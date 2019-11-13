package utils

import "net/url"

func AddURLParams(uri string, params map[string]string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	q := u.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
