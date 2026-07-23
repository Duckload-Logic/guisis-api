package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	docs "github.com/olazo-johnalbert/duckload-api/docs/internal_docs"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/m2mclients"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/integrations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/timeout"
	integrationDocs "github.com/olazo-johnalbert/duckload-api/docs/integrations"
)

func NewRouter(
	db *sqlx.DB,
	handlers *bootstrap.Handlers,
	cfg *config.Config,
) *gin.Engine {
	log.Printf("PRODUCTION MODE: %v", cfg.IsProduction)
	if cfg.IsProduction || cfg.IsStaging {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	g := gin.Default()

	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			target := "dllbsit2027.com"
			parsed, err := url.Parse(origin)
			if err != nil {
				return false
			}
			host := parsed.Hostname()
			if host == target ||
				strings.HasSuffix(host, "."+target) {
				return true
			}
			if !cfg.IsProduction {
				return host == "localhost" ||
					host == "127.0.0.1"
			}
			return false
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Trace-ID",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}

	// Register CORS first to ensure preflights are resolved cleanly
	g.Use(cors.New(corsConfig))

	// Security & DoS Protection
	g.Use(middleware.SecurityHeadersMiddleware())
	g.Use(
		middleware.BodySizeLimitMiddleware(2 << 20),
	) // 2MB limit for JSON payloads

	g.Use(func(c *gin.Context) {
		c.Set(
			middleware.SecurityLoggerContextKey,
			handlers.SystemLogHandler.GetService(),
		)
		c.Next()
	})

	g.Use(middleware.TraceMiddleware())

	limiter := middleware.NewIPRateLimiter(5, 30)
	g.Use(middleware.RateLimitMiddleware(limiter))

	g.Use(func() gin.HandlerFunc {
		shortTimeout := timeout.New(
			timeout.WithTimeout(15*time.Second),
			timeout.WithResponse(func(c *gin.Context) {
				c.JSON(http.StatusGatewayTimeout, gin.H{
					"error": "Request timed out",
				})
			}),
		)
		longTimeout := timeout.New(
			timeout.WithTimeout(3*time.Minute),
			timeout.WithResponse(func(c *gin.Context) {
				c.JSON(http.StatusGatewayTimeout, gin.H{
					"error": "Request timed out",
				})
			}),
		)
		return func(c *gin.Context) {
			path := c.Request.URL.Path
			// Exclude EventSource stream from timeouts
			if strings.Contains(path, "/stream") {
				c.Next()
				return
			}
			method := c.Request.Method
			isLong := strings.Contains(path, "/cors") ||
				strings.Contains(path, "/cor") ||
				strings.Contains(path, "/slips") ||
				strings.Contains(path, "/download") ||
				strings.Contains(path, "/export") ||
				(strings.Contains(path, "/appointments") &&
					(method == http.MethodPost ||
						method == http.MethodPatch))

			if isLong {
				longTimeout(c)
			} else {
				shortTimeout(c)
			}
		}
	}())

	// Stricter limiter for sensitive auth routes
	authLimiter := middleware.NewIPRateLimiter(1, 5)

	apiV1Routes := g.Group("/api/v1")

	apiV1Routes.GET("/docs/internal/*any", func(c *gin.Context) {
		docs.SwaggerInfointernal.Host = c.Request.Host
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.InstanceName("internal"),
		)(c)
	})
	apiV1Routes.GET("/docs/integrations/*any", func(c *gin.Context) {
		integrationDocs.SwaggerInfointegrations.Host = c.Request.Host
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.InstanceName("integrations"),
		)(c)
	})

	files.RegisterRoutes(apiV1Routes, handlers.FileHandler, handlers.Redis)

	apiV1Routes.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{"message": "GuiSIS API version 1.0 initialized"},
		)
	})
	apiV1Routes.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	superadminHealth := apiV1Routes.Group("/logs/system/health")
	superadminHealth.Use(middleware.AuthMiddleware(handlers.Redis))
	superadminHealth.Use(
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
	)
	superadminHealth.GET("", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(), 5*time.Second,
		)
		defer cancel()

		type ServiceHealth struct {
			Name      string `json:"name"`
			Status    string `json:"status"`
			IsHealthy bool   `json:"isHealthy"`
		}

		results := make([]ServiceHealth, 0)

		// 1. API gateway
		results = append(results, ServiceHealth{
			Name:      "API Gateway Server",
			Status:    "Operational",
			IsHealthy: true,
		})

		// 2. MySQL Database
		dbHealthy := true
		dbStatus := "Connected"
		if err := db.PingContext(ctx); err != nil {
			dbHealthy = false
			dbStatus = "Offline"
		}
		results = append(results, ServiceHealth{
			Name:      "MySQL Database",
			Status:    dbStatus,
			IsHealthy: dbHealthy,
		})

		// 3. Redis Cache Store
		redisHealthy := true
		redisStatus := "Connected"
		if err := handlers.Redis.Client.Ping(ctx).Err(); err != nil {
			redisHealthy = false
			redisStatus = "Offline"
		}
		results = append(results, ServiceHealth{
			Name:      "Redis Cache Store",
			Status:    redisStatus,
			IsHealthy: redisHealthy,
		})

		// 4. AI FastAPI Service
		aiHealthy := true
		aiStatus := "Operational"
		aiClient := &http.Client{Timeout: 2 * time.Second}
		aiHealthURL := cfg.AIBaseUrl + "/health"
		if u, err := url.Parse(cfg.AIBaseUrl); err == nil {
			u.Path = "/health"
			aiHealthURL = u.String()
		}
		aiReq, err := http.NewRequestWithContext(
			ctx, "GET", aiHealthURL, nil,
		)
		if err != nil {
			aiHealthy = false
			aiStatus = "Offline"
		} else {
			aiResp, err := aiClient.Do(aiReq)
			if err != nil || aiResp.StatusCode != http.StatusOK {
				aiHealthy = false
				aiStatus = "Offline"
			} else {
				aiResp.Body.Close()
			}
		}
		results = append(results, ServiceHealth{
			Name:      "AI FastAPI Service",
			Status:    aiStatus,
			IsHealthy: aiHealthy,
		})

		// 5. Notification SMTP
		smtpHealthy := true
		smtpStatus := "Active"
		var dialAddr string
		if cfg.IsProduction || cfg.IsStaging {
			dialAddr = fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
		} else {
			dialAddr = fmt.Sprintf(
				"%s:%d", cfg.MailPitHost, cfg.MailPitPort,
			)
		}
		d := net.Dialer{Timeout: 2 * time.Second}
		conn, err := d.DialContext(ctx, "tcp", dialAddr)
		if err != nil {
			smtpHealthy = false
			smtpStatus = "Degraded"
		} else {
			conn.Close()
		}
		results = append(results, ServiceHealth{
			Name:      "Notification SMTP",
			Status:    smtpStatus,
			IsHealthy: smtpHealthy,
		})

		// 6. Identity Provider (IDP)
		idpHealthy := handlers.AuthHandler.IsIDPUp(ctx)
		idpStatus := "Operational"
		if !idpHealthy {
			idpStatus = "Offline"
		}
		results = append(results, ServiceHealth{
			Name:      "Identity Provider (IDP)",
			Status:    idpStatus,
			IsHealthy: idpHealthy,
		})

		c.JSON(http.StatusOK, results)
	})

	authGroup := apiV1Routes.Group("")
	authGroup.Use(middleware.RateLimitMiddleware(authLimiter))
	auth.RegisterRoutes(authGroup, handlers.AuthHandler, handlers.Redis)
	users.RegisterRoutes(db, apiV1Routes, handlers.UserHandler, handlers.Redis)
	locations.RegisterRoutes(
		apiV1Routes,
		handlers.LocationsHandler,
		handlers.Redis,
	)
	students.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.StudentHandler,
		handlers.Redis,
	)
	appointments.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.AppointmentHandler,
		handlers.Redis,
	)
	slips.RegisterRoutes(db, apiV1Routes, handlers.SlipHandler, handlers.Redis)
	analytics.RegisterRoutes(
		apiV1Routes,
		handlers.AnalyticsHandler,
		handlers.Redis,
	)
	m2mclients.RegisterRoutes(
		apiV1Routes,
		handlers.M2MClientHandler,
		handlers.Redis,
	)
	notifications.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.NotificationsHandler,
		handlers.Redis,
	)
	logs.RegisterRoutes(apiV1Routes, handlers.SystemLogHandler, handlers.Redis)
	notes.RegisterRoutes(db, apiV1Routes, handlers.NoteHandler, handlers.Redis)

	integrations.RegisterRoutes(
		apiV1Routes,
		handlers.IntegrationStudentHandler,
		handlers.Redis,
	)

	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Endpoint not found",
		})
	})

	g.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Method not allowed",
		})
	})

	return g
}
