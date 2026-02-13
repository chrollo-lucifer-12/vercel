package server

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
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
	db    *db.DB
}

func NewServerClient(cache *cache.CacheStore, db *db.DB) *ServerClient {
	return &ServerClient{cache: cache, db: db}
}

func (s *ServerClient) trackRequest(subdomain, path, method string, statusCode int, responseTimeMs int, userAgent, ipAddress, referer string) {
	request := db.WebsiteAnalytics{
		Subdomain:      subdomain,
		Path:           path,
		Method:         method,
		StatusCode:     statusCode,
		ResponseTimeMs: responseTimeMs,
		UserAgent:      userAgent,
		IPAddress:      ipAddress,
		Referer:        referer,
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.db.CreateAnalytics(bgCtx, &request); err != nil {
			log.Printf("Failed to track request for %s%s: %v", subdomain, path, err)
		}
	}()
}

func (s *ServerClient) Run(ctx context.Context) error {

	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {

		startTime, ok := resp.Request.Context().Value("startTime").(time.Time)
		if !ok {
			startTime = time.Now()
		}
		responseTimeMs := int(time.Since(startTime).Milliseconds())
		originalHost, _ := resp.Request.Context().Value("originalHost").(string)
		host := strings.Split(originalHost, ":")[0]
		subdomain := strings.TrimSuffix(host, ".localhost")

		s.trackRequest(
			subdomain,
			resp.Request.URL.Path,
			resp.Request.Method,
			resp.StatusCode,
			responseTimeMs,
			resp.Request.UserAgent(),
			strings.Split(resp.Request.RemoteAddr, ":")[0],
			resp.Request.Referer(),
		)

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

		return nil
	}

	proxy.Director = func(req *http.Request) {
		ctx := context.WithValue(req.Context(), "startTime", time.Now())
		ctx = context.WithValue(ctx, "originalHost", req.Host)
		*req = *req.WithContext(ctx)
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
