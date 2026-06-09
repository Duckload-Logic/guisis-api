package server

import (
	"log"
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
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	g := gin.Default()

	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if cfg.IsProduction {
				target := "dllbsit2027.com"
				parsed, err := url.Parse(origin)
				if err != nil {
					return false
				}
				host := parsed.Hostname()
				return host == target || strings.HasSuffix(host, "."+target)
			}

			return strings.HasPrefix(origin, "http://localhost")
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

	g.Use(cors.New(corsConfig))

	// Security & DoS Protection
	g.Use(
		middleware.BodySizeLimitMiddleware(2 << 20),
	) // 2MB limit for JSON payloads
	g.Use(middleware.SecurityHeadersMiddleware())

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
			method := c.Request.Method
			isLong := strings.Contains(path, "/cors") ||
				strings.Contains(path, "/cor") ||
				strings.Contains(path, "/slips") ||
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
	return g
}
