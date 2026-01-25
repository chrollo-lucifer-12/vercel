package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	PORT      = ":8000"
	BASE_HOST = "https://swcxwcivbezgmunayqlf.supabase.co"
	BASE_PATH = "/storage/v1/object/public/builds"
)

func main() {
	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Content-Security-Policy")

		path := resp.Request.URL.Path

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

		return nil
	}

	proxy.Director = func(req *http.Request) {

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

		log.Println("→", req.URL.String())
	}

	log.Println("Reverse Proxy running on", PORT)
	log.Fatal(http.ListenAndServe(PORT, proxy))
}
