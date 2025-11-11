package config

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
)

type ConfigManager struct {
	captainURL string
	pools      map[string]*models.Pool
	mu sync.RWMutex
}

func NewConfigManager(captainURL string) *ConfigManager {
	return &ConfigManager{
		captainURL: captainURL,
		pools:      make(map[string]*models.Pool),
	}
}

func (m *ConfigManager) StartSync(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	// Initial sync
	m.syncConfig()
	
	for range ticker.C {
		m.syncConfig()
	}
}

func (m *ConfigManager) syncConfig() {
	resp, err := http.Get(m.captainURL + "/api/v1/config")
	if err != nil {
		log.Printf("Failed to sync config: %v", err)
		return
	}
	defer resp.Body.Close()
	
	var pools map[string]*models.Pool
	if err := json.NewDecoder(resp.Body).Decode(&pools); err != nil {
		log.Printf("Failed to decode config: %v", err)
		return
	}
	
	m.mu.Lock()
	m.pools = pools
	m.mu.Unlock()
	
	log.Printf("Config synced, %d pools loaded", len(pools))
}

func (m *ConfigManager) GetPools() map[string]*models.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy
	result := make(map[string]*models.Pool)
	for k, v := range m.pools {
		result[k] = v
	}
	
	return result
}

func (m *ConfigManager) GetPool(name string) *models.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.pools[name]
}