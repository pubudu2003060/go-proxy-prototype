package main

import (
	"log"
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

	proxy := proxy.NewHTTPProxy(configManager, authClient, usageReporter)

	go configManager.StartSync(30 * time.Second)

	log.Println("Enhanced HTTP proxy starting on :33080")
	if err := proxy.Start(":33080"); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
}