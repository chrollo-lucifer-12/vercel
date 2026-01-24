package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host

		hostOnly := strings.Split(host, ":")[0]

		parts := strings.Split(hostOnly, ".")

		if len(parts) > 2 {
			subdomain := parts[0]
			resolvesTo :=
		} else {
			fmt.Println("No subdomain")
		}
	})

	http.ListenAndServe(":9000", nil)
}
