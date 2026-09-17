package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	rateLimiterCleanupInterval = 10 * time.Minute
	rateLimiterIdleTimeout     = 15 * time.Minute
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	clients map[string]*clientLimiter
	mu      *sync.RWMutex
	r       rate.Limit
	b       int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		clients: make(map[string]*clientLimiter),
		mu:      &sync.RWMutex{},
		r:       r,
		b:       b,
	}

	go limiter.cleanupRoutine(
		rateLimiterCleanupInterval,
		rateLimiterIdleTimeout,
	)

	return limiter
}

func (i *IPRateLimiter) GetLimiter(key string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	client, exists := i.clients[key]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.clients[key] = &clientLimiter{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func (i *IPRateLimiter) cleanupRoutine(
	interval time.Duration,
	idleTimeout time.Duration,
) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for key, client := range i.clients {
			if now.Sub(client.lastSeen) > idleTimeout {
				delete(i.clients, key)
			}
		}
		i.mu.Unlock()
	}
}

func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limitKey := clientIP
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			limitKey = authHeader
		}

		if !limiter.GetLimiter(limitKey).Allow() {
			if logSvc, ok := c.Get(SecurityLoggerContextKey); ok {
				if svc, ok := logSvc.(SecurityLogger); ok {
					actor := ""

					if isM2M, ok := c.Get("isM2M"); ok && isM2M == true {
						clientName := ginStringVal(c, "clientName")
						m2mID := ginStringVal(c, "m2mClientID")
						if clientName != "" {
							actor = fmt.Sprintf(
								"M2M Client: %s (%s)",
								clientName,
								m2mID,
							)
						} else if m2mID != "" {
							actor = fmt.Sprintf("M2M Client (%s)", m2mID)
						} else {
							actor = "M2M Client"
						}
					}

					if actor == "" {
						if val, exists := c.Get("userEmail"); exists {
							if email, ok := val.(string); ok && email != "" {
								actor = email
							}
						} else if val, exists := c.Get("userID"); exists {
							if id, ok := val.(string); ok && id != "" {
								actor = id
							}
						}
					}

					actorMsg := actor
					if actorMsg == "" {
						actorMsg = fmt.Sprintf("Anonymous (%s)", clientIP)
					}

					svc.RecordSecurity(
						c.Request.Context(),
						"RATE_LIMIT_EXCEEDED",
						fmt.Sprintf(
							"Rate limit exceeded by %s on %s %s",
							actorMsg,
							c.Request.Method,
							c.Request.URL.Path,
						),
						actor,
						clientIP,
						c.Request.UserAgent(),
					)
				}
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}
		c.Next()
	}
}


