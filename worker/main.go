package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/proxy"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

func main() {
	configManager := config.NewConfigManager("http://localhost:8080")
	authClient := auth.NewAuthClient("http://localhost:8080")
	usageReporter := usage.NewUsageReporter("http://localhost:8080")

	wg := sync.WaitGroup{}
	wg.Add(2)
	go startHTTPProxy(&wg, configManager, authClient, usageReporter)
	go startSOCKSProxy(&wg, configManager, authClient, usageReporter)
	wg.Wait()
}

func startHTTPProxy(wg *sync.WaitGroup, configManager *config.ConfigManager, authClient *auth.AuthClient, usageReporter *usage.UsageRepoter) {
	addr := ":8081"

	p := proxy.NewHTTPProxy(configManager, authClient, usageReporter)

	go configManager.StartSync(30 * time.Second)

	srv := &http.Server{
		Addr:         addr,
		Handler:      http.HandlerFunc(p.HandleConnection),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("HTTP/S proxy listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		wg.Done()
		return
	}
}

func startSOCKSProxy(wg *sync.WaitGroup, configManager *config.ConfigManager, authClient *auth.AuthClient, usageReporter *usage.UsageRepoter) {
	addr := ":1080"
	s := proxy.NewSocksProxy(configManager, authClient, usageReporter)
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen :1080: %s", err)
	}
	defer listner.Close()
	defer wg.Done()
	log.Println("SOCKS5 proxy listening on :1080")

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		s.HandleConnection(conn)
	}
}
