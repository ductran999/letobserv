package handler

import (
	"time"

	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/transport/mapper"
	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
)

type inventoryHandler struct {
	productUC usecase.InventoryUsecase
	startUpAt time.Time
}

func NewInventoryHandler(productUC usecase.InventoryUsecase) *inventoryHandler {
	return &inventoryHandler{
		productUC: productUC,
		startUpAt: time.Now(),
	}
}

func (hdl *inventoryHandler) HealthCheck(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": "Healthy",
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *inventoryHandler) ListProducts(c *gin.Context) {
	result, err := hdl.productUC.ListProducts(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}

	response.OK(c, mapper.ToListProductInfoOpenAPI(result), "List products successfully!")
}
