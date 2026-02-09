package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/chrollo-lucifer-12/shared/cache"
	"gorm.io/datatypes"
)

const (
	PORT      = ":8000"
	BASE_HOST = "https://swcxwcivbezgmunayqlf.supabase.co"
	BASE_PATH = "/storage/v1/object/public/builds"
)

type CacheEntry struct {
	Status int         `json:"status"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
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
		originalPath := resp.Request.Header.Get("X-Original-Path")

		resp.Header.Del("Content-Security-Policy")

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		entry := CacheEntry{
			Status: resp.StatusCode,
			Header: resp.Header.Clone(),
			Body:   bodyBytes,
		}

		jsonBytes, err := json.Marshal(entry)
		if err == nil {
			s.cache.Set(ctx, originalPath, datatypes.JSON(jsonBytes))
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

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			value, err := s.cache.Get(ctx, r.URL.Path)
			if err == nil {
				var entry CacheEntry
				if err := json.Unmarshal(value, &entry); err == nil {
					for k, v := range entry.Header {
						for _, vv := range v {
							w.Header().Add(k, vv)
						}
					}

					w.WriteHeader(entry.Status)
					w.Write(entry.Body)
					return
				}
			}
		}

		proxy.ServeHTTP(w, r)
	})

	log.Println("Reverse Proxy running on", PORT)
	err = http.ListenAndServe(PORT, handler)
	return err
}
