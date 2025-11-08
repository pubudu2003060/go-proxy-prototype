package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pubudu2003060/go-proxy-prototype/captain/models"
	"github.com/pubudu2003060/go-proxy-prototype/captain/storage"
)

func CreateUser(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context){
		var req models.CreateUserRequest
		if err := c.ShouldBindJSON(&req);err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}

		user := &models.User{
			Id: uuid.New().String(),
			Username: req.Username,
			Password: req.Password,
			DataLimit: req.DataLimit,
			DataUsed: 0,
			AllowedPools: req.AllowedPools,
			IPWhitelist: req.IPWhitelist,
			Status: "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := storage.CreateUser(user);err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}

		c.JSON(http.StatusCreated,user)
	}
}

func ListUsers (storage *storage.MemoryStorage) gin.HandlerFunc{
	return func(c *gin.Context){
		users,err := storage.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}

		c.JSON(http.StatusOK,users)
	}
}

func GetUser (storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")

		user,err := storage.GetUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}

		user.Password = ""
		c.JSON(http.StatusOK,user)
	}
}

func UpdateUser(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")
		var req models.UpdateUserRequest
		if err := c.ShouldBindJSON(&req);err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}

		err := storage.UpdateUser(id,func (user *models.User) error {
			if req.Password != nil {
				user.Password = *req.Password
			}
			if req.DataLimit != nil {
				user.DataLimit = *req.DataLimit
			}
			if req.AllowedPools != nil {
				user.AllowedPools = *req.AllowedPools
			}
			if req.IPWhitelist != nil {
				user.IPWhitelist = *req.IPWhitelist
			}
			if req.Status != nil {
				user.Status = *req.Status
			}
			return nil
		})

		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return 
		}

		user,_ := storage.GetUser(id)
		user.Password = ""
		c.JSON(http.StatusOK,user)
	}
}

func DeleteUser(storage *storage.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := storage.DeleteUser(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}


