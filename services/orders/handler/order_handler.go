package handler

import (
	"net/http"
	"time"

	"github.com/ductran999/letobserv/services/orders/port"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUC   port.OrderUseCase
	startUpAt time.Time
}

func NewOrderHandler(uc port.OrderUseCase) *orderHandler {
	return &orderHandler{
		orderUC:   uc,
		startUpAt: time.Now(),
	}
}

func (hdl *orderHandler) CheckHealth(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"uptime":  uptime,
		"message": "OK",
	})
}

func (hdl *orderHandler) PlaceOrder(c *gin.Context) {
	err := hdl.orderUC.PlacePOrder(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Placed ordered failed.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Place an order successfully.",
	})
}
