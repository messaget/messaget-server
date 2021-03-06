package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func registerRoutes(cnf *config, r *gin.Engine) {
	var limiter = NewIPRateLimiter(rate.Limit(cnf.Auth.ConnectionsPerSecond), 5)

	// default route to catch root
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "http://www.github.com/messaget")
	})

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	// groups
	authGroup := r.Group("/api")
	publicGroup := r.Group("/public")

	// middleware registry
	publicGroup.Use(func(c *gin.Context) {
		ip, _ := c.RemoteIP()
		if ip != nil {
			limiter := limiter.GetLimiter(ip.String())
			if !limiter.Allow() {
				err := c.AbortWithError(429, RateLimitError)
				if err != nil {
					errorLogger.Println(err)
					c.Abort()
				}
				return
			}
		}

		if cnf.Auth.UseClientPassword {
			p := c.Query("password")
			if p == cnf.Auth.ClientPassword {
				return
			}
			err := c.AbortWithError(401, BadAuthError)
			if err != nil {
				errorLogger.Println(err)
				c.Abort()
			}
		}
	})

	authGroup.Use(func(c *gin.Context) {
		p := c.Query("password")
		if p == cnf.Auth.Password {
			return
		}
		err := c.AbortWithError(401, BadAuthError)
		if err != nil {
			errorLogger.Println(err)
			c.Abort()
		}
	})

	setupMelody()
	setupAdminMelody()
	registerAdminIntents()

	// endpoints
	publicGroup.GET("/attach", handleClientEndpoint)
	authGroup.POST("/intent", handleIntentEndpoint)
	authGroup.GET("/attach", handleControllerSocketEndpoint)
}
