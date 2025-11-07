package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pubudu2003060/go-proxy-prototype/captain/models"
	"github.com/pubudu2003060/go-proxy-prototype/captain/storage"
)

func CreatePool(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreatePoolRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		pool := &models.Pool{
			Name:      req.Name,
			Continent: req.Continent,
			Tag:       req.Tag,
			Subdomain: req.Subdomain,
			CC3:       req.CC3,
			PortStart: req.PortStart,
			PortEnd:   req.PortEnd,
			Outs:      req.Outs,
		}
		
		if err := storage.CreatePool(pool); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, pool)
	}
}

func ListPools(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		pools, err := storage.ListPools()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, pools)
	}
}

func GetPool(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		
		pool, err := storage.GetPool(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, pool)
	}
}

func UpdatePool(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		var pool models.Pool
		if err := c.ShouldBindJSON(&pool); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		pool.Name = name // Ensure name consistency
		
		if err := storage.UpdatePool(name, &pool); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, pool)
	}
}

func DeletePool(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		
		if err := storage.DeletePool(name); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Pool deleted"})
	}
}