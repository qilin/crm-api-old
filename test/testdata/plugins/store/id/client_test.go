package id_test

import (
	"crypto/rsa"
	"os"
	"testing"

	"github.com/qilin/crm-api/test/testdata/plugins/store/id"

	"github.com/dgrijalva/jwt-go"
)

func TestPost(t *testing.T) {
	var rsid, ip, ua string
	var key *rsa.PrivateKey

	key, e := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("APP_SDK_PLUGIN_STORE_KEYS_RAMBLERID_RSAPRIVATEKEY")))
	if e != nil {
		t.Error(e.Error())
	}

	// todo: kid & key
	id := id.Client{
		KID: "rambler_games",
		Key: key,
	}

	rsid = os.Getenv("APP_SDK_PLUGIN_STORE_RSID")

	//ip = ""
	ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:70.0) Gecko/20100101 Firefox/70.0"

	info, err := id.RamblerIdGetProfileInfo(rsid, ip, ua)
	//info := id.RamblerMeta(rsid, ip, ua)
	t.Log(err.Error())
	t.Log(info)
}
