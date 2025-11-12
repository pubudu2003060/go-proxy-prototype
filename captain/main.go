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
		Password:     "testpass", 
		DataLimit:    1000000000, 
		DataUsed:     0,
		AllowedPools: []string{"asia1,eu1,america1"},
		Status:       "active",
		CreatedAt:    time.Now(),
	})

	// Create sample postorage
	storage.CreatePool(&models.Pool{
		Name:      "asia1",
		Continent: "asia",
		Tag:       "asia1",
		Subdomain: "asia1.x",
		CC3:       "asia",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6502,
				Domain:       "netnutasia.x.proxiess.com",
				Weight:       100,
			},
		},
	})

	storage.CreatePool(&models.Pool{
		Name:      "asia2",
		Continent: "asia",
		Tag:       "asia2",
		Subdomain: "asia2.x",
		CC3:       "asia",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6502,
				Domain:       "netnutasia.x.proxiess.com",
				Weight:       100,
			},
		},
	})

	storage.CreatePool(&models.Pool{
		Name:      "eu1",
		Continent: "eu",
		Tag:       "eu1",
		Subdomain: "eu1.x",
		CC3:       "eu",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6501,
				Domain:       "netnuteu.x.proxiess.com",
				Weight:       100,
			},
		},
	})

	storage.CreatePool(&models.Pool{
		Name:      "eu2",
		Continent: "eu",
		Tag:       "eu2",
		Subdomain: "eu2.x",
		CC3:       "eu",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6501,
				Domain:       "netnuteu.x.proxiess.com",
				Weight:       100,
			},
		},
	})

		storage.CreatePool(&models.Pool{
		Name:      "america1",
		Continent: "america",
		Tag:       "america1",
		Subdomain: "america1.x",
		CC3:       "america",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6500,
				Domain:       "devnetnut.x.proxiess.com",
				Weight:       100,
			},
		},
	})

	storage.CreatePool(&models.Pool{
		Name:      "america2",
		Continent: "america",
		Tag:       "america2",
		Subdomain: "america2.x",
		CC3:       "america",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6500,
				Domain:       "devnetnut.x.proxiess.com",
				Weight:       100,
			},
		},
	})	
}