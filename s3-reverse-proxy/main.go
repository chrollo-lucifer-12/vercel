package main

import (
	"context"
	"fmt"
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

	redisUrl := os.Getenv("REDIS_URL")

	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Fatal(err)
	}

	r, _ := redis.NewRedisClient(redisUrl)

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Content-Security-Policy")

		path := resp.Request.URL.Path

		go r.PublishLog(
			context.Background(),
			fmt.Sprintf(
				"RESP %d %s",
				resp.StatusCode,
				path,
			),
			strings.Split(resp.Request.Host, ".")[0],
			"info",
		)

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

		go r.PublishLog(context.Background(), fmt.Sprintf(
			"%s %s | ip=%s | ua=%s",
			req.Method,
			path,
			req.RemoteAddr,
			req.UserAgent(),
		),
			subdomain,
			"info")

		log.Println("→", req.URL.String())
	}

	log.Println("Reverse Proxy running on", PORT)
	log.Fatal(http.ListenAndServe(PORT, proxy))
}
