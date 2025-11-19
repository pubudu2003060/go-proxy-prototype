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
	r.POST("/api/v1/users/proxy-string", handlers.Generate(storage))

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
		AllowedPools: []string{"netnutasia", "netnuteu", "netnutamerica", "iproyalasia", "iproyaleu", "iproyalamerica"},
		Status:       "active",
		CreatedAt:    time.Now(),
	})

	//create sample countries
	japan := models.Country{
		Name: "japan",
		Code: "JP",
	}
	storage.CreateCountry(&japan)

	india := models.Country{
		Name: "india",
		Code: "IN",
	}
	storage.CreateCountry(&india)

	uk := models.Country{
		Name: "united kindom",
		Code: "GB",
	}
	storage.CreateCountry(&uk)

	germany := models.Country{
		Name: "germany",
		Code: "DE",
	}
	storage.CreateCountry(&germany)

	usa := models.Country{
		Name: "usa",
		Code: "US",
	}
	storage.CreateCountry(&usa)

	canada := models.Country{
		Name: "canada",
		Code: "CA",
	}
	storage.CreateCountry(&canada)

	//create sampel pools outs
	//iproyal - username123:password321-country-dk_session-sgn34f3e_lifetime-1h@geo.iproyal.com:12321
	//netnut - USERNAME-res-nl:PASSWORD-sid-947045456@gw.netnut.net:5959
	netnutasia := models.Pool{
		Name:      "netnutasia",
		Region:    "asia",
		Subdomain: "netnutasia.x",
		Port:      6000,
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
		Name:      "iproyalasia",
		Subdomain: "iproyalasia.x",
		Port:      6000,
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
		Name:      "netnuteu",
		Subdomain: "netnuteu.x",
		Port:      6000,
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
		Name:      "iproyaleu",
		Subdomain: "iproyaleu.x",
		Port:      6000,
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
		Name:      "netnutamerica",
		Subdomain: "netnutamerica.x",
		Port:      6000,
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
		Name:      "iproyalamerica",
		Subdomain: "iproyalamerica.x",
		Port:      6000,
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
		Name:       "asia",
		SubDomains: []string{"iproyalamerica.x", "netnutamerica.x"},
	}

	worker2 := models.Worker{
		Name:       "eu",
		SubDomains: []string{"iproyalasia.x", "netnutasia.x"},
	}

	worker3 := models.Worker{
		Name:       "america",
		SubDomains: []string{"iproyaleu.x", "netnuteu.x"},
	}

	storage.CreateWorker(&worker1)
	storage.CreateWorker(&worker2)
	storage.CreateWorker(&worker3)

	//create sampe regions
	asia := models.Region{
		Name:      "asia",
		Countries: []models.Country{japan, india},
		Pools:     []models.Pool{iproyalasia, netnutasia},
	}
	storage.CreateRegion(&asia)

	eu := models.Region{
		Name:      "eu",
		Countries: []models.Country{uk, germany},
		Pools:     []models.Pool{iproyaleu, netnuteu},
	}
	storage.CreateRegion(&eu)

	america := models.Region{
		Name:      "america",
		Countries: []models.Country{usa, canada},
		Pools:     []models.Pool{iproyalamerica, netnutamerica},
	}
	storage.CreateRegion(&america)

}
