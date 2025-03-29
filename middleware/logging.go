package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// LoggingMiddleware returns a middleware function that logs HTTP requests
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			path := req.URL.Path
			method := req.Method

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			stop := time.Now()
			latency := stop.Sub(start)
			statusCode := res.Status
			ip := c.RealIP()
			userAgent := req.UserAgent()

			// Get the user ID if available (from authenticated request)
			userID := "anonymous"
			if user, ok := c.Get("user").(*struct{ ID string }); ok && user != nil {
				userID = user.ID
			}

			log.Infof("[%s] %s %s %d %s | Latency: %v | IP: %s | User-Agent: %s | User: %s",
				method, path, req.Proto, statusCode, http.StatusText(statusCode), latency, ip, userAgent, userID)

			return err
		}
	}
}
