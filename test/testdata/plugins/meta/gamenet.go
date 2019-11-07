package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pascaldekloe/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type plugin struct {
	//
}

var (
	Plugin plugin
)

func (p *plugin) Name() string {
	return "gamenet.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		// qilinProductUUID, ok := token.String("qilinProductUUID")
		// if !ok {
		// 	qilinProductUUID = "unknown"
		// }
		// userID, ok := token.String("userID")
		// if !ok {
		// 	userID = "unknown"
		// }

		meta := map[string]string{}
		meta["mode"] = "gamenet"
		// meta["qilinProductUUID"] = qilinProductUUID
		// meta["userID"] = userID

		//if authenticate == nil {
		//	return authenticate(ctx, request, token, log)
		//}
		return common.AuthResponse{
			Meta: meta,
		}, nil
	}
}

func init() {
	// starting proxy server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", root())
		if err := http.ListenAndServe(":1443", mux); err != nil {
			fmt.Println("gamenet plugin:", err.Error())
		}
	}()
}

// eyJhbGciOiJFUzUxMiJ9.eyJleHAiOjE1NzM0ODM2NjUsImlzcyI6IlFpbGluIiwicWlsaW5Qcm9kdWN0VVVJRCI6IjNkNGZmNWY5LTg2MTQtNDUyNC1iYTRiLTM3OGE5ZmRiNDU5NCIsInN1YiI6IlFpbGluU3ViamVjdCIsInVzZXJJRCI6IjEwMDUwMCJ9.QDwRpjt93j0oFdHUq9MZEQ8RBJ01QdFeCUz3qppb61b60qq0g_gOQCd-8NuwADtgwUfC4IRwMVfzCixXpJ5ug83lHTprQmXfyyUsSg-nlZ89CFuiCC_PuZkH2CJQKqU5
func root() http.HandlerFunc {
	target, _ := url.Parse("https://gamenet.ru")
	targetQuery := target.RawQuery
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = "gamenet.ru"
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
		ModifyResponse: func(res *http.Response) error {
			for _, v := range res.Header {
				for i := range v {
					if strings.Contains(v[i], "https://gamenet.ru") {
						v[i] = strings.Replace(v[i], "https://gamenet.ru", "http://localhost:1443", -1)
					}
				}
			}
			return nil
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			jwt := r.URL.Query().Get("jwt")
			_ = jwt

			http.SetCookie(w, &http.Cookie{
				Name:  "PHPSESSID",
				Value: "gr0k682k9906qq9s95u3iuqqa5",
			})
			w.WriteHeader(http.StatusOK)
			w.Write(index)
			return
		}
		proxy.ServeHTTP(w, r)
		return
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

var index = []byte(`
	<iframe id="iFrameClient" name="iFrameClient" frameborder="0" width="100%" height="100%" allowfullscreen="allowfullscreen" wmode="Opaque" data-bind="attr: {src: frameSrc}, style: { visibility: hideClient() ? 'hidden' : 'visible' }" src="http://localhost:1443/games/khanwars/iframe?wmode=opaque" style="visibility: visible;"></iframe>
`)
