package common

type OAuth2 struct {
	ClientId     string `required:"true"`
	ClientSecret string `required:"true"`
	RedirectUrl  string `required:"true"`
}

type CORS struct {
	Allowed []string
	Headers []string
	Methods []string
}
