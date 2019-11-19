package rambler

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
)

// if keys is empty than all params will be used for signature except 'sig'
func VerifySignature(params url.Values, secret string, keys ...string) bool {
	sig := Sign(params, secret, keys...)
	return sig == params.Get("sig")
}

// if keys is empty than all params will be used for signature except 'sig'
func Sign(params url.Values, secret string, keys ...string) string {
	if len(keys) == 0 {
		keys = make([]string, 0, len(params))
		for k := range params {
			if k != "sig" {
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	h := md5.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte("="))
		h.Write([]byte(params.Get(k)))
		h.Write([]byte("&"))
	}
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum(nil))
}

func SignUrl(uri string, secret string, keys ...string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	params := u.Query()
	sig := Sign(params, secret, keys...)
	params.Add("sig", sig)
	u.RawQuery = params.Encode()
	return u.String(), nil
}
