package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type CorsOption struct {
	AllowOrigins []string
	MaxAge       time.Duration
}

func Cors(option *CorsOption) gin.HandlerFunc {
	var (
		allowOrigins []string
		maxAge       = time.Hour * 2
	)

	if option != nil {
		allowOrigins = option.AllowOrigins

		if option.MaxAge != 0 {
			maxAge = option.MaxAge
		}
	}

	return cors.New(
		cors.Config{
			AllowOriginFunc: func(origin string) bool {
				if len(allowOrigins) == 0 {
					return true
				} else {
					for _, allowOrigin := range allowOrigins {
						if allowOrigin == origin {
							return true
						}
					}
					return false
				}
			},
			//"Access-Control-Allow-Origin"は今のchromeのcookieのheaderの仕様により、"*"を設定するとcookieが送信されない
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type", "X-Request-Id"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowCredentials: true,
			MaxAge:           maxAge,
		},
	)
}
