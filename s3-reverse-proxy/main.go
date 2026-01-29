package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/chrollo-lucider-12/proxy/redis"
	"github.com/joho/godotenv"
)

const (
	PORT      = ":8000"
	BASE_HOST = "https://swcxwcivbezgmunayqlf.supabase.co"
	BASE_PATH = "/storage/v1/object/public/builds"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	redisURL := os.Getenv("REDIS_URL")
	r, err := redis.NewRedisClient(redisURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {
		originalPath := resp.Request.Header.Get("X-Original-Path")
		method := resp.Request.Method
		status := resp.StatusCode
		path := resp.Request.URL.Path
		parts := strings.Split(path, "/")

		r.PublishLog(context.Background(), status, parts[4], originalPath, method)

		resp.Header.Del("Content-Security-Policy")

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
		req.Header.Set("X-Original-Path", req.URL.Path)

		host := strings.Split(req.Host, ":")[0]
		subdomain := strings.Split(host, ".")[0]
		path := req.URL.Path

		var targetPath string

		if strings.Contains(path, ".") {
			targetPath = path
		} else {

			if path == "/" {
				targetPath = "/index.html"
			} else {
				targetPath = "/index.html"
			}
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path = BASE_PATH + "/" + subdomain + targetPath

	}

	log.Println("Reverse Proxy running on", PORT)
	log.Fatal(http.ListenAndServe(PORT, proxy))
}
