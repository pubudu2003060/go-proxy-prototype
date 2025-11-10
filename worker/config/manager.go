package config

import (
	"sync"

	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
)

type ConfigManager struct {
	captainURL string
	pools      map[string]*models.Pool
	mu sync.RWMutex
}