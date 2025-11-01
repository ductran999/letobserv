package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ductran999/letobserv/services/products/domain"
	"github.com/ductran999/letobserv/services/products/port"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productUC port.ProductUseCase
	startUpAt time.Time
}

func NewProductHandler(uc port.ProductUseCase) *productHandler {
	return &productHandler{
		productUC: uc,
		startUpAt: time.Now(),
	}
}

func (hdl *productHandler) CheckHealth(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"uptime":  uptime,
		"message": "OK",
	})
}

func (hdl *productHandler) ReduceProductStock(c *gin.Context) {
	id := c.Param("id")

	if err := hdl.productUC.ReduceProductStock(c.Request.Context(), id); err != nil {
		log.Printf("[ERR] Failed to reduce stock for product %s: %v", id, err)
		if errors.Is(err, domain.ErrProductInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, domain.ErrOutOfStock) {
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "OK"})
}
