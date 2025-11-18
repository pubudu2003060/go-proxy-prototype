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
			Region:    req.Region,
			Subdomain: req.Subdomain,
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
		var req models.UpdatePoolRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := storage.UpdatePool(name, func(pool *models.Pool) error {
			if req.Region != nil {
				pool.Region = *req.Region
			}
			if req.Outs != nil {
				pool.Outs = *req.Outs
			}
			if req.PortEnd != nil {
				pool.PortEnd = *req.PortEnd
			}
			if req.PortStart != nil {
				pool.PortStart = *req.PortStart
			}
			if req.Subdomain != nil {
				pool.Subdomain = *req.Subdomain
			}
			return nil
		}); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, struct{}{})
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
