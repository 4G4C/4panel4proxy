package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kataras/iris/v12"
)

func reverseProxy(ctx iris.Context, domain Domain) bool {
	// Set the target URL of the backend server
	targetURL, err := url.Parse("http://" + domain.ReverseProxy)
	if err != nil {
		return false
	}
	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// Wrap the proxy with a handler that adds a timeout.
	proxy.ErrorHandler = func(http.ResponseWriter, *http.Request, error) {
		ctx.HTML("Tidak dapat menemukan backend")
	}
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	return true
}
