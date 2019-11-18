package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	nonceLength = 16
)

type Payload struct {
	Host       string
	Uri        string
	Method     string
	Headers    string
	Ctime      string
	Nonce      string
	BodySha256 string `json:"body_sha256"`
}

type RPC struct {
	RPC    string
	Method string
	Params map[string]string
}

type IDClient struct {
	KID string
	Key *rsa.PrivateKey
}

func (id IDClient) post(body []byte) []byte {
	url := "https://id.rambler.ru/jsonrpc"

	h := sha256.New()
	h.Write([]byte(body))

	nonce := make([]byte, nonceLength)
	rand.Read(nonce)

	token := jwt.New(jwt.SigningMethodRS256)
	token.Header["kid"] = id.KID
	token.Claims = jwt.MapClaims{
		"host":        "id.rambler.ru",
		"uri":         "/jsonrpc",
		"method":      "POST",
		"headers":     "",
		"ctime":       time.Now().Unix(),
		"nonce":       nonce,
		"body_sha256": string(h.Sum(nil)),
	}
	JWT, err := token.SignedString(id.Key)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-RID-Signature", JWT)

	client := &http.Client{
		Timeout: time.Second * time.Duration(30),
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	return b
}

func (id IDClient) RamblerIdGetProfileInfo(rsid, ip, ua string) string {
	rpc := RPC{
		RPC:    "2.0",
		Method: "Rambler::id::get_profile_info",
		Params: map[string]string{
			"rsid":       rsid,
			"__clientIp": ip,
			"__clientUa": ua,
		},
	}
	body, _ := json.Marshal(rpc)
	profileInfo := id.post(body)

	// todo: parse
	return string(profileInfo)
}
