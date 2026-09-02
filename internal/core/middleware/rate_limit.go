package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(key string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[key]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = i.ips[key]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[key] = limiter
	}

	return limiter
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
					userEmail := ""
					if val, exists := c.Get("userEmail"); exists {
						if email, ok := val.(string); ok {
							userEmail = email
						}
					} else if val, exists := c.Get("userID"); exists {
						if id, ok := val.(string); ok {
							userEmail = id
						}
					}

					actorMsg := userEmail
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
						userEmail,
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
