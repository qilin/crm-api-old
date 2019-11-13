package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
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
		meta := map[string]string{
			"mode": "gamenet",
			"url":  "/games/khanwars/iframe?wmode=opaque",
		}

		return common.AuthResponse{
			Meta: meta,
		}, nil
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	cfg, ok := ctx.Value("config").(map[string]string)
	if !ok {
		log.Error("plugin: can not cast context config to map[string]string")
	}

	port, ok := cfg["port"]
	if !ok {
		port = "1443"
	}

	url, ok := cfg["entryurl"]
	if !ok {
		log.Error("plugin: can not find entryURL in config")
	}

	run(port, url)

}

func run(port, url string) {
	// starting proxy server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/sdk/v1/auth", func(w http.ResponseWriter, r *http.Request) {
			var rs = common.AuthResponse{
				Meta: map[string]string{
					"url": "/games/khanwars/iframe?wmode=opaque",
				},
			}
			data, err := json.Marshal(&rs)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		})
		mux.HandleFunc("/", root(url))
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			fmt.Println("gamenet plugin:", err.Error())
		}
	}()
}

// eyJhbGciOiJFUzUxMiJ9.eyJleHAiOjE1NzM0ODM2NjUsImlzcyI6IlFpbGluIiwicWlsaW5Qcm9kdWN0VVVJRCI6IjNkNGZmNWY5LTg2MTQtNDUyNC1iYTRiLTM3OGE5ZmRiNDU5NCIsInN1YiI6IlFpbGluU3ViamVjdCIsInVzZXJJRCI6IjEwMDUwMCJ9.QDwRpjt93j0oFdHUq9MZEQ8RBJ01QdFeCUz3qppb61b60qq0g_gOQCd-8NuwADtgwUfC4IRwMVfzCixXpJ5ug83lHTprQmXfyyUsSg-nlZ89CFuiCC_PuZkH2CJQKqU5
func root(proxyurl string) http.HandlerFunc {
	target, _ := url.Parse("https://gamenet.ru")
	targetQuery := target.RawQuery
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = "gamenet.ru"
			// follow lines of code copied from httputil.NewReverseProxy
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
			fmt.Println(res.Request.RequestURI)
			for _, v := range res.Header {
				for i := range v {
					if strings.Contains(v[i], "https://gamenet.ru") {
						v[i] = strings.Replace(v[i], "https://gamenet.ru", proxyurl, -1)
					}
				}
			}
			res.Header.Del("X-Frame-Options")

			if res.Header.Get("Content-Encoding") == "gzip" {

				gr, err := gzip.NewReader(res.Body)
				if err != nil {
					return err
				}

				data, err := ioutil.ReadAll(gr)
				if err != nil {
					return err
				}
				res.Body.Close()

				if bytes.Contains(data, []byte("window.top")) {
					// fmt.Println("replace in", res.Request.RequestURI)
					data = bytes.ReplaceAll(data, []byte("window.top"), []byte("window.parent"))
				}
				var buf bytes.Buffer
				w := gzip.NewWriter(&buf)
				w.Write(data)
				w.Close()

				res.Body = ioutil.NopCloser(&buf)
			}
			return nil
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			jwt := r.URL.Query().Get("jwt")
			_ = jwt

			// TODO user from jws token
			session, err := auth("aleksandr.barsukov@protocol.one", "password")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "PHPSESSID",
				Value: session,
			})
			w.WriteHeader(http.StatusOK)
			w.Write(index)
			return
		}
		proxy.ServeHTTP(w, r)
		return
	}
}

func auth(login, password string) (string, error) {
	var c = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	form := make(url.Values)
	form.Add("login", login)
	form.Add("password", password)
	form.Add("captcha", "")
	form.Add("mid", "175419290")
	form.Add("browserhwid", "1ba8e830bc511735b9676f94455feab9")
	form.Add("rp", "")
	form.Add("trustedLocation", "1")
	form.Add("code2fa", "")
	form.Add("json", "1")
	form.Add("2fa", "1")

	req, err := http.NewRequest(http.MethodPost, "https://gnlogin.ru/", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("not authorized")
	}
	return ssid(resp)
}

func ssid(resp *http.Response) (string, error) {
	for _, c := range resp.Cookies() {
		if c.Name == "PHPSESSID" {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("session not found")
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

var index = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
<script type="text/javascript">
parcelRequire=function(e,r,t,n){var i,o="function"==typeof parcelRequire&&parcelRequire,u="function"==typeof require&&require;function f(t,n){if(!r[t]){if(!e[t]){var i="function"==typeof parcelRequire&&parcelRequire;if(!n&&i)return i(t,!0);if(o)return o(t,!0);if(u&&"string"==typeof t)return u(t);var c=new Error("Cannot find module '"+t+"'");throw c.code="MODULE_NOT_FOUND",c}p.resolve=function(r){return e[t][1][r]||r},p.cache={};var l=r[t]=new f.Module(t);e[t][0].call(l.exports,p,l,l.exports,this)}return r[t].exports;function p(e){return f(p.resolve(e))}}f.isParcelRequire=!0,f.Module=function(e){this.id=e,this.bundle=f,this.exports={}},f.modules=e,f.cache=r,f.parent=o,f.register=function(r,t){e[r]=[function(e,r){r.exports=t},{}]};for(var c=0;c<t.length;c++)try{f(t[c])}catch(e){i||(i=e)}if(t.length){var l=f(t[t.length-1]);"object"==typeof exports&&"undefined"!=typeof module?module.exports=l:"function"==typeof define&&define.amd?define(function(){return l}):n&&(this[n]=l)}if(parcelRequire=f,i)throw i;return f}({"nnXf":[function(require,module,exports) {
"use strict";Object.defineProperty(exports,"__esModule",{value:!0}),exports.default=void 0;var e=function(e){var t=document.createElement("iframe");return t.src=e,t.style.width="100%",t.style.height="100%",t.width="100%",t.height="100%",document.body.prepend(t),t};exports.default=e;
},{}],"eKDL":[function(require,module,exports) {
"use strict";Object.defineProperty(exports,"__esModule",{value:!0}),exports.FULLSCREEN_MODE_CHANGED=exports.ENABLE_FULLSCREEN=exports.PAYMENT_FORM_CLOSED=exports.SHOW_PAYMENT_FORM=void 0;var E="SHOW_PAYMENT_FORM";exports.SHOW_PAYMENT_FORM=E;var _="PAYMENT_FORM_CLOSED";exports.PAYMENT_FORM_CLOSED=_;var L="ENABLE_FULLSCREEN";exports.ENABLE_FULLSCREEN=L;var N="FULLSCREEN_MODE_CHANGED";exports.FULLSCREEN_MODE_CHANGED=N;
},{}],"MddR":[function(require,module,exports) {
"use strict";Object.defineProperty(exports,"__esModule",{value:!0}),exports.default=void 0;var e=function(e,t,n,r){return new(n||(n=Promise))(function(o,a){function u(e){try{c(r.next(e))}catch(t){a(t)}}function i(e){try{c(r.throw(e))}catch(t){a(t)}}function c(e){var t;e.done?o(e.value):(t=e.value,t instanceof n?t:new n(function(e){e(t)})).then(u,i)}c((r=r.apply(e,t||[])).next())})},t=function(e,t){var n,r,o,a,u={label:0,sent:function(){if(1&o[0])throw o[1];return o[1]},trys:[],ops:[]};return a={next:i(0),throw:i(1),return:i(2)},"function"==typeof Symbol&&(a[Symbol.iterator]=function(){return this}),a;function i(a){return function(i){return function(a){if(n)throw new TypeError("Generator is already executing.");for(;u;)try{if(n=1,r&&(o=2&a[0]?r.return:a[0]?r.throw||((o=r.return)&&o.call(r),0):r.next)&&!(o=o.call(r,a[1])).done)return o;switch(r=0,o&&(a=[2&a[0],o.value]),a[0]){case 0:case 1:o=a;break;case 4:return u.label++,{value:a[1],done:!1};case 5:u.label++,r=a[1],a=[0];continue;case 7:a=u.ops.pop(),u.trys.pop();continue;default:if(!(o=(o=u.trys).length>0&&o[o.length-1])&&(6===a[0]||2===a[0])){u=0;continue}if(3===a[0]&&(!o||a[1]>o[0]&&a[1]<o[3])){u.label=a[1];break}if(6===a[0]&&u.label<o[1]){u.label=o[1],o=a;break}if(o&&u.label<o[2]){u.label=o[2],u.ops.push(a);break}o[2]&&u.ops.pop(),u.trys.pop();continue}a=t.call(e,u)}catch(i){a=[6,i],r=0}finally{n=o=0}if(5&a[0])throw a[1];return{value:a[0]?a[1]:void 0,done:!0}}([a,i])}}},n=function(n){return function(r,o,a){return e(void 0,void 0,void 0,function(){var e,u;return t(this,function(t){switch(t.label){case 0:e={meta:r,url:o},a&&(e.qilinProductUUID=a),t.label=1;case 1:return t.trys.push([1,4,,5]),[4,fetch(n+"/auth",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(e)})];case 2:return[4,t.sent().json()];case 3:return[2,t.sent().meta];case 4:throw u=t.sent(),console.error(u),u;case 5:return[2]}})})}};exports.default=n;
},{}],"nxQv":[function(require,module,exports) {
"use strict";Object.defineProperty(exports,"__esModule",{value:!0}),exports.default=void 0;var e=require("./constants"),t=n(require("./getAuthFunction"));function n(e){return e&&e.__esModule?e:{default:e}}var r=function(e,t,n,r){return new(n||(n=Promise))(function(o,i){function a(e){try{l(r.next(e))}catch(t){i(t)}}function u(e){try{l(r.throw(e))}catch(t){i(t)}}function l(e){var t;e.done?o(e.value):(t=e.value,t instanceof n?t:new n(function(e){e(t)})).then(a,u)}l((r=r.apply(e,t||[])).next())})},o=function(e,t){var n,r,o,i,a={label:0,sent:function(){if(1&o[0])throw o[1];return o[1]},trys:[],ops:[]};return i={next:u(0),throw:u(1),return:u(2)},"function"==typeof Symbol&&(i[Symbol.iterator]=function(){return this}),i;function u(i){return function(u){return function(i){if(n)throw new TypeError("Generator is already executing.");for(;a;)try{if(n=1,r&&(o=2&i[0]?r.return:i[0]?r.throw||((o=r.return)&&o.call(r),0):r.next)&&!(o=o.call(r,i[1])).done)return o;switch(r=0,o&&(i=[2&i[0],o.value]),i[0]){case 0:case 1:o=i;break;case 4:return a.label++,{value:i[1],done:!1};case 5:a.label++,r=i[1],i=[0];continue;case 7:i=a.ops.pop(),a.trys.pop();continue;default:if(!(o=(o=a.trys).length>0&&o[o.length-1])&&(6===i[0]||2===i[0])){a=0;continue}if(3===i[0]&&(!o||i[1]>o[0]&&i[1]<o[3])){a.label=i[1];break}if(6===i[0]&&a.label<o[1]){a.label=o[1],o=i;break}if(o&&a.label<o[2]){a.label=o[2],a.ops.push(i);break}o[2]&&a.ops.pop(),a.trys.pop();continue}i=t.call(e,a)}catch(u){i=[6,u],r=0}finally{n=o=0}if(5&i[0])throw i[1];return{value:i[0]?i[1]:void 0,done:!0}}([i,u])}}},i=function(n,i,a){if(!n)throw new Error("Game UID is required, but not provided");if(!i)throw new Error("Api URL is required, but not provided");var u=window.location.href,l={},c=!1,s=(0,t.default)(i);return window.addEventListener("message",function(e){if(c){var t,n,r=e.data,o=void 0===r?{}:r,i=o.type,a=o.payload;if(i)t=a,(n=l[i])&&n.forEach(function(e){return e(t)})}}),{enableFullscreenMode:function(){c?window.parent.postMessage({type:e.ENABLE_FULLSCREEN},"*"):alert("Game is not initialized!")},showPaymentForm:function(t,r,o){if(c){var i={type:e.SHOW_PAYMENT_FORM,payload:{qilinProductUUID:o||n,userId:r,itemId:t}};window.parent.postMessage(i,"*")}else alert("Game is not initialized!")},addCallback:function(e,t){var n=l[e]||[];n.push(t),l[e]=n},init:function(e){return r(void 0,void 0,void 0,function(){var t,n;return o(this,function(r){switch(r.label){case 0:return r.trys.push([0,2,,3]),[4,(a||s)(e,u)];case 1:return t=r.sent(),c=!0,[2,t];case 2:throw n=r.sent(),console.error(n),n;case 3:return[2]}})})}}};exports.default=i;
},{"./constants":"eKDL","./getAuthFunction":"MddR"}],"iWQp":[function(require,module,exports) {
"use strict";Object.defineProperty(exports,"__esModule",{value:!0}),exports.default=void 0;var e=n(require("../src/openIframe")),t=n(require("../src/qilinGameFrame")),r=require("../src/constants");function n(e){return e&&e.__esModule?e:{default:e}}var o=function(e,t,r,n){return new(r||(r=Promise))(function(o,a){function i(e){try{c(n.next(e))}catch(t){a(t)}}function u(e){try{c(n.throw(e))}catch(t){a(t)}}function c(e){var t;e.done?o(e.value):(t=e.value,t instanceof r?t:new r(function(e){e(t)})).then(i,u)}c((n=n.apply(e,t||[])).next())})},a=function(e,t){var r,n,o,a,i={label:0,sent:function(){if(1&o[0])throw o[1];return o[1]},trys:[],ops:[]};return a={next:u(0),throw:u(1),return:u(2)},"function"==typeof Symbol&&(a[Symbol.iterator]=function(){return this}),a;function u(a){return function(u){return function(a){if(r)throw new TypeError("Generator is already executing.");for(;i;)try{if(r=1,n&&(o=2&a[0]?n.return:a[0]?n.throw||((o=n.return)&&o.call(n),0):n.next)&&!(o=o.call(n,a[1])).done)return o;switch(n=0,o&&(a=[2&a[0],o.value]),a[0]){case 0:case 1:o=a;break;case 4:return i.label++,{value:a[1],done:!1};case 5:i.label++,n=a[1],a=[0];continue;case 7:a=i.ops.pop(),i.trys.pop();continue;default:if(!(o=(o=i.trys).length>0&&o[o.length-1])&&(6===a[0]||2===a[0])){i=0;continue}if(3===a[0]&&(!o||a[1]>o[0]&&a[1]<o[3])){i.label=a[1];break}if(6===a[0]&&i.label<o[1]){i.label=o[1],o=a;break}if(o&&i.label<o[2]){i.label=o[2],i.ops.push(a);break}o[2]&&i.ops.pop(),i.trys.pop();continue}a=t.call(e,i)}catch(u){a=[6,u],n=0}finally{r=o=0}if(5&a[0])throw a[1];return{value:a[0]?a[1]:void 0,done:!0}}([a,u])}}},i=function(e,t){var r={};for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&t.indexOf(n)<0&&(r[n]=e[n]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols){var o=0;for(n=Object.getOwnPropertySymbols(e);o<n.length;o++)t.indexOf(n[o])<0&&Object.prototype.propertyIsEnumerable.call(e,n[o])&&(r[n[o]]=e[n[o]])}return r},u="buyItem",c="onSuccessBuy",s="onCancelBuy",l=function(n){if(!n)throw new Error("Api URL is required, but not provided");var l,f=(0,t.default)("PROXY",n),p=!1,d=function(e){if(l){var t=e.status;console.log("Gamenet Adapter Payform status: ",t);var r=JSON.stringify({method:t?c:s});l.postMessage(r,"*")}};return{init:function(t){return o(void 0,void 0,void 0,function(){var n,o,c;return a(this,function(a){switch(a.label){case 0:return a.trys.push([0,2,,3]),[4,f.init(t)];case 1:return n=a.sent(),o=n.url,p=!0,(0,e.default)(o),window.addEventListener("message",function(e){if(p){var t=e.data,r=void 0===t?{}:t;if("string"==typeof r){var n=JSON.parse(r),o=n.method,a=i(n,["method"]);if(!o)return;if(l||(l=e.source),o===u){var c=a.itemId;f.showPaymentForm(c,"","")}}}}),f.addCallback(r.PAYMENT_FORM_CLOSED,d),[3,3];case 2:throw c=a.sent(),console.error(c),c;case 3:return[2]}})})}}};exports.default=l;
},{"../src/openIframe":"nnXf","../src/qilinGameFrame":"nxQv","../src/constants":"eKDL"}],"QCba":[function(require,module,exports) {
"use strict";var r=e(require("./qilinGameProxy"));function e(r){return r&&r.__esModule?r:{default:r}}if(!window)throw new Error("SDK only work on web browser");window.qilinGameProxy=r.default;
},{"./qilinGameProxy":"iWQp"}]},{},["QCba"], null)
//# sourceMappingURL=/gamenet-adapter.js.map
</script>
  <style>
    html, body {
      height: 100%;
      overflow: hidden;
    }
  </style>
  <title>Document</title>
</head>
<body>
  <script>
    const helper = qilinGameProxy('/sdk/v1');  
    helper.init()
      .then(() => console.log('Adapter was started'))
      .catch(err => console.log(err));
  </script>
<!-- <iframe id="iFrameClient" name="iFrameClient" frameborder="0" width="100%" height="100%" allowfullscreen="allowfullscreen" wmode="Opaque" data-bind="attr: {src: frameSrc}, style: { visibility: hideClient() ? "hidden" : "visible" }" src="http://localhost:1443/games/khanwars/iframe?wmode=opaque" style="visibility: visible;"></iframe> -->

</body>
</html>
`)
