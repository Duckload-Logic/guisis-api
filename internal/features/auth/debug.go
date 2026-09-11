package auth

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterDebugRoutes(
	rg *gin.RouterGroup,
	redis *datastore.RedisClient,
) {
	// Strictly disabled in production and staging environments
	if os.Getenv("IS_PRODUCTION") == "true" ||
		os.Getenv("IS_STAGING") == "true" {
		return
	}

	debugGroup := rg.Group("/debug")
	{
		debugGroup.GET("/redis", func(c *gin.Context) {
			secret := c.GetHeader("X-Debug-Secret")
			expectedSecret := os.Getenv("REDIS_DEBUG_SECRET")

			if expectedSecret == "" || secret != expectedSecret {
				c.JSON(
					http.StatusUnauthorized,
					gin.H{"error": "Unauthorized debug access"},
				)
				return
			}

			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			// Safe non-blocking key inspection: scan up to 50 keys max
			keys, _, err := redis.Client.Scan(ctx, 0, "*", 50).Result()
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": err.Error()},
				)
				return
			}

			details := make(map[string]interface{})
			for _, key := range keys {
				val, _ := redis.Client.Get(ctx, key).Result()
				details[key] = val
			}

			c.JSON(http.StatusOK, gin.H{
				"totalKeys": len(keys),
				"keys":      keys,
				"data":      details,
				"server":    "Development Redis",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		})
	}
}
