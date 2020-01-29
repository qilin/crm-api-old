package auth

import (
	"context"
	"testing"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
)

func TestInitDisabled(t *testing.T) {
	// TODO implement noop awareset
	set := provider.AwareSet{
		Logger: logger.NewMock(context.TODO(), &logger.Config{}, true),
	}
	_, err := New(nil, set, AppSet{}, &Config{Enabled: false})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestInitNoKeys(t *testing.T) {
	// TODO implement noop awareset
	set := provider.AwareSet{
		Logger: logger.NewMock(context.TODO(), &logger.Config{}, true),
	}
	_, err := New(nil, set, AppSet{}, &Config{Enabled: true})
	if err == nil {
		t.Error("must fails when no jwt keys configured")
		t.FailNow()
	}
}
