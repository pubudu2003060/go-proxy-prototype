package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/proxy"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

func main() {
	addr:= ":8081"

	configManager := config.NewConfigManager("http://localhost:8080")

	authClient := auth.NewAuthClient("http://localhost:8080")

	usageReporter := usage.NewUsageReporter("http://localhost:8080")

	proxy := proxy.NewHTTPProxy(configManager, authClient, usageReporter)

	go configManager.StartSync(30 * time.Second)

	srv := &http.Server{
		Addr:         addr,
		Handler:      http.HandlerFunc(proxy.HandleConnection),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("proxy listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}