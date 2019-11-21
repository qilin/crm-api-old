package utils_test

import (
	"crypto/rsa"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/test/testdata/plugins/store/utils"
)

var PK = `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQC72hHjg1JZhf08
Fi2kdbP5QUL57LwxUHxndyVn78q/cbaL9NnvZikhRngCgXgP2ujk1WQkvoBk6IN/
ywpTXZ9PnkRdE3E3YDWzyeYE3PGA32L8gygPXstFSLPjq+hi1QqSTlU4kXxGyZWT
4SL4IGKWa/e2bNWsce2FAt6i3fQ23riIc/LFOdtRS8msJYm9It/Y4QM6mVkS3Xao
flzHj1n/0ZIdj4Norl2Cc9M8KnL4N5j0Uo5lplydTQmmuKjcY6NxT63Jvsl7vLBQ
r8jszWVTkxBPTnY0Ozm2EELZK9/Wf/alOogwsSPGKcCMBhK4K8LUgQNVE3KoQTts
SKmsOf0HAgMBAAECggEBAI8DeHcDdWBel+p04A7C5V+wBbOMPcI1imCi3sGAV5Tk
l5t5r6mI12tT4O3Xb3ZyrLf6laE7vzgTpHlYNOY+8piE33sU4C7OelQEM7AkHWCF
sTCZEZiSDKMUtI9yQxtrIf88z7ifWkPyGIRa9Gp/DU+DpzUlKo98tN18z86it9hE
7b3JWCI44m6ipOyzYj2VW6no8AmvB+KEKKZIZDDnL4EjXEHjhFKeGjSv7Iehnhd+
EqzT4TUH330p8wvMecgo3rHJ02jZFnkekTURGBHwPesptWyeL+4D5aZDyGLmGwnJ
mY9ugjEqES8WboHJKv0pAWgwTGQdgNOWpgT+ceIwimECgYEA3fJqkIf6qNSEE0Lt
z0gT51AgreSuFUgQOp5otfKFr5zn9b5AzQAFcDWEHcZ6eUBBr6PGLMny84FN8FRo
12DMIJPsz2IGrSx6na4HPHZs+R4YxJd5+HHayLVVbPnvjk5OuXo3AJNBoEXAepNR
GdudFMt8vS4vuYkYxvo0whxzsvkCgYEA2Kx5Oh41yKkQeHwTxC4xX0ITVKe4YQUe
eIpZ5nLLf23TnNvDYJluiynPSSOP06mNX8BOaQ22sRHZuXpgSmIuLh6t2GNoHQmA
7KNS6JOuwFu8Rp3HBmdduUkro7rF6I4HK0lopAgNzKl4cTlD6hU7qMkL4qGcRATr
4bmRqwjXL/8CgYEAgn+01MJ/SaGa/tBNj6ErwshETrq0+OJkWHMn0kOFA1rYsI9q
/p5SlEWDJxa6kGyNsr4zGcasSSzwLK0U7/6ER2tyxAU5M72BYxEeRBjFvjxKB92g
48neAEFOt0LF7gBxHXGUwYvT/G7G28ue1fthAwcakwmDGi5YGTaoqrGb/dECgYEA
xhR+yuPdVXFBjnQX+ewk1KnqVCT0STXN0nLglu1xHjDAGRFLPt9lkLGLP5jUHrNN
fDCpPh78WkowWgEHUFkLULxZP445GvqaMztoSxjf1BjJOWF6Fl+e7gl3bLoNvXlC
Eo+MqxB11RlE83VfofsBTF9njyshWYmKPxRPmCV/2+8CgYEAuQXX3hqdX/bhFonG
QN/TsGUnVMlEgHskohSfsbhGEHkth5gblZarW6NGs8/KZvrw/gfpiekzJGFJwDmD
uIcX6Ea48NZBqguK/3rQEdELOzzFH8KgrVCaGL/kDX1U4XhdrvO6wPPiw6nN8aTm
/ldNB3H94kmDg7Ik/YVFWysaTt8=
-----END PRIVATE KEY-----`

func TestPost(t *testing.T) {
	var rsid, ip, ua string
	var key *rsa.PrivateKey

	key, e := jwt.ParseRSAPrivateKeyFromPEM([]byte(PK))
	if e != nil {
		t.Error(e.Error())
	}

	// todo: kid & key
	id := utils.IDClient{
		KID: "rambler_games",
		Key: key,
	}

	rsid = "eyJleHRyYSI6eyJkYXRhIjoiSGJYaVoyZGsxdHNoU2RpVWVpSVhhUnZwSm91YXNiN24xaDQ2SXVzTFlBZXZaUStOY2pDaGU0ajMzdGFoNGFQXC9PT29zUDJvU3NSMmtJNHd5XC9BVGpFeG1EdTluZmZ3WjlUQ1FwRjZOWWNJamFiVkR0WWpMZTBBbk00dlhtc0hJY212RGViTFRLY1hVZVNcL2JwVHU5ZktSUHAyNFFnNE1xeVpsSEdMS3hWeStKaWorWXRNb2t0Zlc0Z1NybGY1YmcyRHlSWUNJQjBmUHk5R1lFTXo1UXRIeXhQIiwiZW5jX2tleSI6ImtleTEifSwicnNpZCI6IjUxNjZiZmEyYzUxZWJkYTNkMzNhMGNhZjhmODA5YTAwIn0.v2.x"

	//ip = ""
	ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:70.0) Gecko/20100101 Firefox/70.0"

	//info := id.RamblerIdGetProfileInfo(rsid, ip, ua)
	info := id.RamblerMeta(rsid, ip, ua)

	t.Log(info)
}
