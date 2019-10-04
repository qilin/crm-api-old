package common

type AuthConfig struct {
	Issuer       string `envconfig:"AUTH1_ISSUER" default:"https://dev-auth1.tst.protocol.one"`
	ClientId     string `envconfig:"AUTH1_CLIENTID" required:"true"`
	ClientSecret string `envconfig:"AUTH1_CLIENTSECRET" required:"true"`
	RedirectUrl  string `envconfig:"AUTH1_REDIRECTURL" required:"true"`
}
