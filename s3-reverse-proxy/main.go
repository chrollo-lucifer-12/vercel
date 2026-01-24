package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

const BASE_PATH = "https://swcxwcivbezgmunayqlf.supabase.co/storage/v1/object/public/builds"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		hostName := strings.Split(host, ":")[0]
		subdomain := strings.Split(hostName, ".")[0]

		target, err := url.Parse(BASE_PATH)
		if err != nil {
			log.Printf("Error parsing URL: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		// Add error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			cleanPath := path.Clean(req.URL.Path)
			if cleanPath == "/" || cleanPath == "." {
				cleanPath = "/index.html"
			}

			req.URL.Path = path.Join(target.Path, subdomain, cleanPath)

			// Log the request
			log.Printf("Proxying: %s -> %s", r.URL.Path, req.URL.String())
		}

		// Remove CSP headers and add no-cache headers
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Remove all CSP headers
			resp.Header.Del("Content-Security-Policy")
			resp.Header.Del("Content-Security-Policy-Report-Only")
			resp.Header.Del("X-Content-Security-Policy")

			// Prevent caching during development
			resp.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			resp.Header.Set("Pragma", "no-cache")
			resp.Header.Set("Expires", "0")

			return nil
		}

		proxy.ServeHTTP(w, r)
	})

	log.Println("Server starting on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
