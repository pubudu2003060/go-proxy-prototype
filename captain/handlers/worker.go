package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pubudu2003060/go-proxy-prototype/captain/models"
	"github.com/pubudu2003060/go-proxy-prototype/captain/storage"
)

func getConfig(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		pools, err := storage.GetAllPools()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, pools)
	}
}

func authenticateUser(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		user, err := storage.GetUserByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusOK, models.AuthResponse{
				Success: false,
				Message: "Invalid credentials",
			})
			return
		}
		
		// Simple password check - in production, use hashed passwords
		if user.Password != req.Password {
			c.JSON(http.StatusOK, models.AuthResponse{
				Success: false,
				Message: "Invalid credentials",
			})
			return
		}
		
		if user.Status != "active" {
			c.JSON(http.StatusOK, models.AuthResponse{
				Success: false,
				Message: "User account is " + user.Status,
			})
			return
		}
		
		c.JSON(http.StatusOK, models.AuthResponse{
			Success:      true,
			UserID:       user.Id,
			AllowedPools: user.AllowedPools,
			DataLimit:    user.DataLimit,
			DataUsed:     user.DataUsed,
		})
	}
}

func reportUsage(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UsageReport
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := storage.UpdateUser(req.UserID, func(user *models.User) error {
			user.DataUsed += req.Bytes
			return nil
		})
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Usage reported"})
	}
}