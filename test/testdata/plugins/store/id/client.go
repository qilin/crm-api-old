package id

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	nonceLength  = 16
	ramblerIDURL = "https://id.rambler.ru/jsonrpc"
)

type Client struct {
	KID string
	Key *rsa.PrivateKey
}

func (id Client) post(body []byte) []byte {
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
		"body_sha256": hex.EncodeToString(h.Sum(nil)),
	}
	JWT, err := token.SignedString(id.Key)

	// setup a http client
	client := &http.Client{
		Timeout: time.Second * time.Duration(30),
	}
	r, err := http.NewRequest("POST", ramblerIDURL, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-RID-Signature", JWT)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	return b
}

func (id Client) RamblerIdGetProfileInfo(rsid, ip, ua string) (*RamblerProfile, error) {
	rpc := RPCRequest{
		RPC:    "2.0",
		Method: "Rambler::Id::get_profile_info",
		Params: []map[string]string{
			map[string]string{
				"rsid":       rsid,
				"__clientIp": ip,
				"__clientUa": ua,
			},
		},
	}
	body, _ := json.Marshal(rpc)
	profileInfo := id.post(body)

	// parse response
	resp := GetProfileInfoResponse{}
	err := json.Unmarshal(profileInfo, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Result.Status != strings.ToLower("ok") {
		return nil, errors.New("not authenticated")
	}

	return &RamblerProfile{
		Email:     resp.Result.Info.Email,
		Gender:    resp.Result.Profile.Gender,
		Birthdate: resp.Result.Profile.Birthdate,
	}, nil
}
