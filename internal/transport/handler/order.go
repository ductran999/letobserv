package handler

import (
	"net/http"
	"time"

	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUsecase usecase.OrderUseCase
	startUpAt    time.Time
}

func NewOrderHandler(orderUsecase usecase.OrderUseCase) *orderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
		startUpAt:    time.Now(),
	}
}

func (hdl *orderHandler) CheckHealth(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": consts.HealthyState,
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *orderHandler) PlaceOrder(c *gin.Context) {
	err := hdl.orderUsecase.PlacePOrder(c.Request.Context())
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
