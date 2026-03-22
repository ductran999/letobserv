package commonhttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	HealthyState = "heathy"
)

type CommonHandler interface {
	HealthCheck(c *gin.Context)
	ReadinessCheck(c *gin.Context)
}

type handler struct {
	db        *gorm.DB
	startupAt time.Time

	// cache readiness result
	lastCheck     time.Time
	lastReady     bool
	checkMu       sync.Mutex
	checkInterval time.Duration
}

func New(db *gorm.DB) CommonHandler {
	return &handler{
		db:            db,
		startupAt:     time.Now(),
		checkInterval: 5 * time.Second,
	}
}

func (hdl *handler) HealthCheck(c *gin.Context) {
	uptime := int64(time.Since(hdl.startupAt).Seconds())
	resp := gin.H{
		"status": HealthyState,
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (h *handler) ReadinessCheck(c *gin.Context) {
	h.checkMu.Lock()
	defer h.checkMu.Unlock()

	if time.Since(h.lastCheck) < h.checkInterval {
		if !h.lastReady {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
		return
	}

	sqlDB, err := h.db.DB()
	h.lastReady = err == nil && sqlDB.Ping() == nil
	h.lastCheck = time.Now()

	if !h.lastReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
