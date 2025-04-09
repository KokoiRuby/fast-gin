package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"
)

func LimitMiddleware(limit int) gin.HandlerFunc {
	return NewLimiter(limit, 1*time.Second).Middleware
}

func NewLimiter(limit int, duration time.Duration) *Limiter {
	return &Limiter{
		limit:      limit,
		duration:   duration,
		timestamps: make(map[string][]int64),
	}
}

type Limiter struct {
	limit      int                // Limit
	duration   time.Duration      // Window size
	timestamps map[string][]int64 // Timestamp
}

func (l *Limiter) Middleware(c *gin.Context) {
	ip := c.ClientIP()

	// If timestamp of corresponding ip exists
	if _, ok := l.timestamps[ip]; !ok {
		l.timestamps[ip] = make([]int64, 0)
	}

	now := time.Now().Unix()

	// Remove expired timestamp
	for i := 0; i < len(l.timestamps[ip]); i++ {
		if l.timestamps[ip][i] < now-int64(l.duration.Seconds()) {
			l.timestamps[ip] = append(l.timestamps[ip][:i], l.timestamps[ip][i+1:]...)
			i--
		}
	}

	if len(l.timestamps[ip]) >= l.limit {
		c.JSON(429, gin.H{
			"message": "Too Many Requests",
		})
		c.Abort()
		return
	}

	l.timestamps[ip] = append(l.timestamps[ip], now)

	// Continue
	c.Next()
}
