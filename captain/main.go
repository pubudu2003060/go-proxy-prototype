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
	//r.POST("/api/v1/generate_proxy_string", handlers.Generate(storage))

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
		AllowedPools: []string{"netnut asia", "netnut eu", "netnut america", "iproyal asia", "iproyal eu", "iproyal america"},
		Status:       "active",
		CreatedAt:    time.Now(),
	})

	//create sample countries
	japan := models.Country{
		Code: "JP",
	}
	storage.CreateCountry(&japan)

	india := models.Country{
		Code: "IN",
	}
	storage.CreateCountry(&india)

	uk := models.Country{
		Code: "GB",
	}
	storage.CreateCountry(&uk)

	germany := models.Country{
		Code: "DE",
	}
	storage.CreateCountry(&germany)

	usa := models.Country{
		Code: "US",
	}
	storage.CreateCountry(&usa)

	canada := models.Country{
		Code: "CA",
	}
	storage.CreateCountry(&canada)

	//create sampel pools outs
	//iproyal - username123:password321-country-dk_session-sgn34f3e_lifetime-1h@geo.iproyal.com:12321
	//netnut - USERNAME-res-nl:PASSWORD-sid-947045456@gw.netnut.net:5959
	netnutasia := models.Pool{
		Name:      "netnut asia",
		Continent: "asia",
		Tag:       "asia1",
		Subdomain: "netnutasia.x",
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
	}
	storage.CreatePool(&netnutasia)

	iproyalasia := models.Pool{
		Name:      "iproyal asia",
		Continent: "asia",
		Tag:       "asia2",
		Subdomain: "iproyalasia.x",
		CC3:       "asia",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "otJhMuv0:5uhhT0Ds-%s",
				UpstreamPort: 12322,
				Domain:       "iproyalasia.x.proxiess.com",
				Weight:       100,
			},
		},
	}
	storage.CreatePool(&iproyalasia)

	netnuteu := models.Pool{
		Name:      "netnut eu",
		Continent: "eu",
		Tag:       "eu1",
		Subdomain: "netnuteu.x",
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
	}
	storage.CreatePool(&netnuteu)

	iproyaleu := models.Pool{
		Name:      "iproyal eu",
		Continent: "eu",
		Tag:       "eu2",
		Subdomain: "iproyaleu.x",
		CC3:       "eu",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "otJhMuv0:5uhhT0Ds-%s",
				UpstreamPort: 12323,
				Domain:       "iproyaleu.x.proxiess.com",
				Weight:       100,
			},
		},
	}
	storage.CreatePool(&iproyaleu)

	netnutamerica := models.Pool{
		Name:      "netnut america",
		Continent: "america",
		Tag:       "america1",
		Subdomain: "netnutamerica.x",
		CC3:       "america",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "cFAPhxyG:9dgbjKKV-%s",
				UpstreamPort: 6500,
				Domain:       "netnut.x.proxiess.com",
				Weight:       100,
			},
		},
	}

	storage.CreatePool(&netnutamerica)

	iproyalamerica := models.Pool{
		Name:      "iproyal america",
		Continent: "america",
		Tag:       "america2",
		Subdomain: "iproyalamerica.x",
		CC3:       "america",
		PortStart: 6000,
		PortEnd:   6999,
		Flag:      0,
		Outs: []models.Out{
			{
				Format:       "otJhMuv0:5uhhT0Ds-%s",
				UpstreamPort: 12321,
				Domain:       "iproyal.x.proxiess.com",
				Weight:       100,
			},
		},
	}

	storage.CreatePool(&iproyalamerica)

	worker1 := models.Worker{
		WName: "asia",
		Pools: []*models.Pool{&netnutasia, &iproyalasia},
	}

	worker2 := models.Worker{
		WName: "eu",
		Pools: []*models.Pool{&netnuteu, &iproyaleu},
	}

	worker3 := models.Worker{
		WName: "america",
		Pools: []*models.Pool{&netnutamerica, &iproyalamerica},
	}

	storage.CreateWorker(&worker1)
	storage.CreateWorker(&worker2)
	storage.CreateWorker(&worker3)

	//create sampe regions
	asia := models.Region{
		RName:     "asia",
		Countries: []models.Country{japan, india},
		Workers:   []models.Worker{worker1},
	}
	storage.CreateRegion(&asia)

	eu := models.Region{
		RName:     "eu",
		Countries: []models.Country{uk, germany},
		Workers:   []models.Worker{worker2},
	}
	storage.CreateRegion(&eu)

	america := models.Region{
		RName:     "america",
		Countries: []models.Country{usa, canada},
		Workers:   []models.Worker{worker3},
	}
	storage.CreateRegion(&america)

}
