package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/storage"
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
	cache   *cache.CacheStore
	db      *db.DB
	storage *storage.S3Storage
}

func NewServerClient(cache *cache.CacheStore, db *db.DB, storage *storage.S3Storage) *ServerClient {
	return &ServerClient{
		cache:   cache,
		db:      db,
		storage: storage,
	}
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

func (s *ServerClient) handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	host := strings.Split(r.Host, ":")[0]
	subdomain := strings.Split(host, ".")[0]

	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	objectKey := subdomain + path

	ctx := r.Context()

	if cachedJSON, err := s.cache.Get(ctx, objectKey); err == nil && len(cachedJSON) > 0 {

		var cachedEntry CacheEntry
		if err := json.Unmarshal(cachedJSON, &cachedEntry); err == nil {
			for k, vv := range cachedEntry.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(cachedEntry.Status)
			_, _ = w.Write([]byte(cachedEntry.Body))

			fmt.Printf("Served from DB cache: %s/%s\n", subdomain, path)

			responseTime := int(time.Since(start).Milliseconds())
			s.trackRequest(
				subdomain,
				path,
				r.Method,
				cachedEntry.Status,
				responseTime,
				r.UserAgent(),
				strings.Split(r.RemoteAddr, ":")[0],
				r.Referer(),
			)
			return
		}
	}

	reader, err := s.storage.GetObject(ctx, objectKey)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer reader.Close()

	switch {
	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(path, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	}

	w.Header().Set("Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data:;")

	responseBuffer := &strings.Builder{}

	writer := io.MultiWriter(w, responseBuffer)
	_, err = io.Copy(writer, reader)
	if err != nil {
		log.Printf("Failed to write response for %s%s: %v", subdomain, path, err)
		return
	}

	cachedEntry := CacheEntry{
		Status: http.StatusOK,
		Header: w.Header().Clone(),
		Body:   responseBuffer.String(),
	}

	cachedJSON, err := json.Marshal(cachedEntry)
	if err == nil {
		if err := s.cache.Set(ctx, objectKey, cachedJSON); err != nil {
			log.Printf("Failed to cache %s/%s: %v", subdomain, path, err)
		}
	}

	responseTime := int(time.Since(start).Milliseconds())
	s.trackRequest(
		subdomain,
		path,
		r.Method,
		http.StatusOK,
		responseTime,
		r.UserAgent(),
		strings.Split(r.RemoteAddr, ":")[0],
		r.Referer(),
	)
}

func (s *ServerClient) Run(ctx context.Context) error {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	log.Println("Server running on", PORT)

	return http.ListenAndServe(PORT, mux)
}
