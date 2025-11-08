package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pubudu2003060/go-proxy-prototype/captain/handlers"
	"github.com/pubudu2003060/go-proxy-prototype/captain/models"
	"github.com/pubudu2003060/go-proxy-prototype/captain/storage"
)

func main() {

	storage := storage.NewMemoryStorage()
	
	initSampleData(storage)

	r := gin.Default()

	// User management
	r.POST("/api/v1/users", handlers.CreateUser(storage))
	r.GET("/api/v1/users", handlers.ListUsers(storage))
	r.GET("/api/v1/users/:id", handlers.GetUser(storage))
	r.PUT("/api/v1/users/:id", handlers.UpdateUser(storage))
	r.DELETE("/api/v1/users/:id", handlers.DeleteUser(storage))
	
	// Pool management
	r.POST("/api/v1/pools", handlers.CreatePool(storage))
	r.GET("/api/v1/pools", handlers.ListPools(storage))
	r.GET("/api/v1/pools/:name", handlers.GetPool(storage))
	r.PUT("/api/v1/pools/:name", handlers.UpdatePool(storage))
	r.DELETE("/api/v1/pools/:name", handlers.DeletePool(storage))
	
	// Worker endpoints
	r.GET("/api/v1/config", handlers.GetConfig(storage))
	r.POST("/api/v1/auth", handlers.AuthenticateUser(storage))
	r.POST("/api/v1/usage", handlers.ReportUsage(storage))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Captain API server starting on :8080")
	log.Fatal(r.Run(":8080"))

}

func initSampleData(storage *storage.MemoryStorage) {
	// Create sample users
	storage.CreateUser(&models.User{
		Id:           "user1",
		Username:     "testuser",
		Password:     "testpass", // In real implementation, this should be hashed
		DataLimit:    1000000000, // 1GB
		DataUsed:     0,
		AllowedPools: []string{"smart-usa", "mixed-pool"},
		Status:       "active",
		CreatedAt:    time.Now(),
	})

	// Create sample postorage
	storage.CreatePool(&models.Pool{
		Name:      "smart-usa",
		Continent: "usa",
		Tag:       "usa",
		Subdomain: "smart.us",
		CC3:       "SMRT",
		PortStart: 7000,
		PortEnd:   7999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "user-%s",
				UpstreamPort: 8000,
				Domain:       "smartproxy.com",
				Weight:       50,
			},
		},
	})

	storage.CreatePool(&models.Pool{
		Name:      "mixed-pool",
		Continent: "global",
		Tag:       "global",
		Subdomain: "mixed.global",
		CC3:       "MIX",
		PortStart: 8000,
		PortEnd:   8999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "user-%s-country-us",
				UpstreamPort: 8000,
				Domain:       "brightdata.com",
				Weight:       30,
			},
			{
				Format:       "user-%s-country-eu",
				UpstreamPort: 8000,
				Domain:       "oxylabs.io",
				Weight:       30,
			},
			{
				Format:       "user-%s-session-%s",
				UpstreamPort: 8000,
				Domain:       "netnut.io",
				Weight:       40,
			},
		},
	})
}