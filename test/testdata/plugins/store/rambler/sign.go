package rambler

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
)

func VerifySignature(params url.Values, secret string, keys ...string) bool {
	sig := Sign(params, secret, keys...)
	return sig == params.Get("sig")
}

func Sign(params url.Values, secret string, keys ...string) string {
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
