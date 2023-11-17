package main

import (
	"context"
	"os"

	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var uriPattern = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}(/.*)?$`)

type HTTPCheck struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name" binding:"required"`
	URI                  string    `json:"uri" binding:"required"`
	IsPaused             bool      `json:"is_paused" binding:"required"`
	NumRetries           int       `json:"num_retries" binding:"required"`
	UptimeSLA            int       `json:"uptime_sla" binding:"required"`
	ResponseTimeSLA      int       `json:"response_time_sla" binding:"required"`
	UseSSL               bool      `json:"use_ssl"`
	ResponseStatusCode   int       `json:"response_status_code"`
	CheckIntervalSeconds int       `json:"check_interval_in_seconds"`
	CheckCreated         time.Time `json:"check_created"`
	CheckUpdated         time.Time `json:"check_updated"`
}

func setupRouter(repo HTTPCheckRepositoryDecorator) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		// Attempt to ping the database
		err := repo.PingDB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Database not available"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
		}
	})

	r.GET("v1/http-checks", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		checksChan := make(chan []HTTPCheck)
		errChan := make(chan error)

		go func() {
			checks, err := repo.GetAll()
			if err != nil {
				errChan <- err
			} else {
				checksChan <- checks
			}
		}()

		select {
		case checks := <-checksChan:
			c.JSON(http.StatusOK, checks)
		case err := <-errChan:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
		}
	})

	r.GET("v1/http-check/:id", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		id := c.Param("id")
		check, err := repo.GetByID(id, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if check == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "HTTP Check not found"})
			return
		}
		c.JSON(http.StatusOK, check)
	})

	r.POST("v1/http-check", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var check HTTPCheck
		if err := c.ShouldBindJSON(&check); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := createCustomResource(check); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Create(&check, ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, check)
	})

	r.PUT("v1/http-check/:id", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		id := c.Param("id")

		var check HTTPCheck
		if err := c.ShouldBindJSON(&check); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := updateCustomResource(check); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}
		if err := repo.Update(id, &check, ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, check)
	})

	r.DELETE("v1/http-check/:id", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		id := c.Param("id")
		check, err := repo.GetByID(id, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if check == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "HTTP Check not found"})
			return
		}
		if err := deleteCustomResource(check.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := repo.Delete(id, ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	})

	return r
}

func main() {
	maxRetries := 5
	retryInterval := 5 * time.Second

	var db *gorm.DB
	var err error
	for i := 1; i <= maxRetries; i++ {
		db, err = openDatabaseConnection()
		if err == nil {
			break
		}

		log.Printf("Failed to connect to the database (attempt %d/%d). Retrying in %s...\n", i, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to connect to the database after %d retries. Error: %s", maxRetries, err)
	}

	db.AutoMigrate(&HTTPCheck{})
	healthCheckRepo := &HTTPCheckGORMRepository{db}
	decoratedRepo := &LoggingDecorator{Repo: healthCheckRepo}

	r := setupRouter(decoratedRepo)
	r.Run(":8080")
}

func openDatabaseConnection() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	// dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1045 || mysqlErr.Number == 1049 {
				http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("Waiting for the database to become available..."))
				})

				go http.ListenAndServe(":8080", nil)
			}
		}
		return nil, err
	}

	return db, nil
}

//testing
