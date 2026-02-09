package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/chrollo-lucifer-12/shared/cache"
)

const (
	PORT      = ":8000"
	BASE_HOST = "https://swcxwcivbezgmunayqlf.supabase.co"
	BASE_PATH = "/storage/v1/object/public/builds"
)

type CacheEntry struct {
	Status int         `json:"status"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

type ServerClient struct {
	cache *cache.CacheStore
}

func NewServerClient(cache *cache.CacheStore) *ServerClient {
	return &ServerClient{cache: cache}
}

func (s *ServerClient) Run(ctx context.Context) error {

	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {

		path := resp.Request.URL.Path
		resp.Header.Del("Set-Cookie")
		resp.Header.Del("Content-Security-Policy")

		resp.Header.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:;")

		switch {
		case strings.HasSuffix(path, ".html"):
			resp.Header.Set("Content-Type", "text/html; charset=utf-8")
		case strings.HasSuffix(path, ".js"):
			resp.Header.Set("Content-Type", "application/javascript")
		case strings.HasSuffix(path, ".css"):
			resp.Header.Set("Content-Type", "text/css")
		case strings.HasSuffix(path, ".svg"):
			resp.Header.Set("Content-Type", "image/svg+xml")
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("body", string(bodyBytes))

		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		return nil
	}

	proxy.Director = func(req *http.Request) {
		req.Header.Set("X-Original-Path", req.URL.Path)

		host := strings.Split(req.Host, ":")[0]
		subdomain := strings.Split(host, ".")[0]
		path := req.URL.Path

		if path == "/" {
			path = "/index.html"
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path = BASE_PATH + "/" + subdomain + path
	}

	log.Println("Reverse Proxy running on", PORT)
	err = http.ListenAndServe(PORT, proxy)
	return err
}
