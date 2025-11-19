package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/captain/models"
)

type MemoryStorage struct {
	users   map[string]*models.User
	pools   map[string]*models.Pool
	Workers map[string]*models.Worker
	Region  map[string]*models.Region
	Country map[string]*models.Country
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:   make(map[string]*models.User),
		pools:   make(map[string]*models.Pool),
		Workers: make(map[string]*models.Worker),
		Region:  make(map[string]*models.Region),
		Country: make(map[string]*models.Country),
	}
}

func (s *MemoryStorage) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if u.Username == user.Username {
			return fmt.Errorf("username %s already Exit", user.Username)
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.Id] = user
	fmt.Printf("user created %v \n", s.users[user.Id].Id)

	return nil
}

func (s *MemoryStorage) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *MemoryStorage) GetUserByUsername(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (s *MemoryStorage) ListUsers() ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := []*models.User{}
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

func (s *MemoryStorage) UpdateUser(id string, updateFun func(*models.User) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}

	if err := updateFun(user); err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return nil
}

func (s *MemoryStorage) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := s.users[id]
	if user == nil {
		return fmt.Errorf("user not found")
	}

	delete(s.users, id)
	return nil
}

func (s *MemoryStorage) CreatePool(pool *models.Pool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.pools[pool.Name]; ok {
		return fmt.Errorf("pool with %s already exit", pool.Name)
	}
	for _, v := range s.pools {
		if v.Subdomain == pool.Subdomain {
			return fmt.Errorf("pool with %s already exit", pool.Subdomain)
		}
	}

	s.pools[pool.Name] = pool
	fmt.Printf("pool created %v \n", s.pools[pool.Name].Name)
	return nil
}

func (s *MemoryStorage) GetPool(name string) (*models.Pool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pool, exists := s.pools[name]
	if !exists {
		return nil, fmt.Errorf("pool not found")
	}

	return pool, nil
}

func (s *MemoryStorage) ListPools() ([]*models.Pool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pools := make([]*models.Pool, 0, len(s.pools))
	for _, pool := range s.pools {
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *MemoryStorage) UpdatePool(name string, updateFunc func(model *models.Pool) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pool, exists := s.pools[name]
	if !exists {
		return fmt.Errorf("pool not found")
	}

	if err := updateFunc(pool); err != nil {
		return fmt.Errorf("error in user request")
	}

	return nil
}

func (s *MemoryStorage) DeletePool(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.pools[name]; !ok {
		return fmt.Errorf("pool not found")
	}

	delete(s.pools, name)
	return nil
}

func (s *MemoryStorage) GetAllPools() (map[string]*models.Pool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*models.Pool)
	for k, v := range s.pools {
		result[k] = v
	}

	return result, nil
}

func (s *MemoryStorage) CreateWorker(worker *models.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Workers[worker.Name]; ok {
		return fmt.Errorf("worker with %s already Exit", worker.Name)
	}

	s.Workers[worker.Name] = worker
	fmt.Printf("worker created %v \n", s.Workers[worker.Name])

	return nil
}

func (s *MemoryStorage) CreateRegion(region *models.Region) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Region[region.Name]; ok {
		return fmt.Errorf("region with %s already Exit", region.Name)
	}

	s.Region[region.Name] = region
	fmt.Printf("region created %v \n", s.Region[region.Name].Name)

	return nil
}

func (s *MemoryStorage) CreateCountry(country *models.Country) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Country[country.Name]; ok {
		return fmt.Errorf("cuntry with %s already Exit", country.Code)
	}

	s.Country[country.Code] = country
	fmt.Printf("country created %v \n", s.Country[country.Code])

	return nil
}
