package handler

import (
	"time"

	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productUC usecase.ProductUseCase
	startUpAt time.Time
}

func NewProductHandler(productUC usecase.ProductUseCase) *productHandler {
	return &productHandler{
		productUC: productUC,
		startUpAt: time.Now(),
	}
}

func (hdl *productHandler) CheckHealth(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": "Healthy",
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *productHandler) ReduceProductStock(c *gin.Context) {
	id := c.Param("id")

	if err := hdl.productUC.ReduceProductStock(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
	}

	response.OK(c, nil, "Reduce product stock successfully!")
}
