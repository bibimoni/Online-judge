package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/bibimoni/Online-judge/gateway/src/infrastructure/config"
)

func SubmissionApiProxy() http.Handler {
	cfg := config.Load()
	return newProxy(cfg.Endpoints.Submission)
}

func ProblemApiProxy() http.Handler {
	cfg := config.Load()
	return newProxy(cfg.Endpoints.Problem)
}

func LoginApiProxy() http.Handler {
	cfg := config.Load()
	return newProxy(cfg.Endpoints.Auth)
}

func newProxy(endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(endpoint)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		httpProxy := httputil.NewSingleHostReverseProxy(url)
		httpProxy.ServeHTTP(w, r)
	})
}

func WSSubmissionProxy() http.Handler {
	cfg := config.Load()
	return WSProxy(cfg.Endpoints.Submission)
}

func WSProxy(endpoint string) http.Handler {
	target, err := url.Parse(endpoint)
	if err != nil {
		panic("Invalid proxy URL: " + err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: false,
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		if strings.EqualFold(req.Header.Get("Connection"), "upgrade") {
			req.Header.Set("Connection", "upgrade")
		}
		if strings.EqualFold(req.Header.Get("Upgrade"), "websocket") {
			req.Header.Set("Upgrade", "websocket")
		}
	}

	return proxy
}
