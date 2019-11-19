package rambler

import (
	"net/url"
	"testing"
)

func TestIframeSignature(t *testing.T) {
	u, err := url.Parse("https://localhost/path?user_id=1683086&slug=protocolone&game_id=1079&timestamp=1573646629488&sig=89bda2b84d827138e934c529384a7b6c")
	if err != nil {
		t.Fatal(err)
	}
	ok := VerifySignature(u.Query(), "6f12ff821d49e386c0918415322d0b74",
		"user_id", "game_id", "slug", "timestamp")
	if !ok {
		t.Error("signature not match")
	}
}

func TestIframeSignatureWithoutKeys(t *testing.T) {
	u, err := url.Parse("https://localhost/path?user_id=1683086&slug=protocolone&game_id=1079&timestamp=1573646629488&sig=89bda2b84d827138e934c529384a7b6c")
	if err != nil {
		t.Fatal(err)
	}
	ok := VerifySignature(u.Query(), "6f12ff821d49e386c0918415322d0b74")
	if !ok {
		t.Error("signature not match")
	}
}

func TestGetItemSignature(t *testing.T) {
	u, err := url.Parse("https://localhost/some?item=100500&app_id=1079&user_id=1683086&receiver_id=1683086&lang=ru_RU&notification_type=get_item&sig=58b41bb2565f7e2b9626bd21f3691e4d")
	if err != nil {
		t.Fatal(err)
	}
	ok := VerifySignature(u.Query(), "6f12ff821d49e386c0918415322d0b74",
		"item", "app_id", "user_id", "receiver_id", "lang", "notification_type")
	if !ok {
		t.Error("signature not match")
	}
}
