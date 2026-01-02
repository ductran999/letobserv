package inventory

import (
	"time"

	generated "github.com/ductran999/letobserv/api/generated/inventory"
	"github.com/ductran999/letobserv/internal/application/usecase/inventory"
	"github.com/ductran999/letobserv/pkg/request"
	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	ReservationErrorCode = "RESERVED_FAILED"
)

type inventoryHandler struct {
	inventoryUC inventory.InventoryUsecase
	startUpAt   time.Time
}

func NewInventoryHandler(inventoryUC inventory.InventoryUsecase) *inventoryHandler {
	return &inventoryHandler{
		inventoryUC: inventoryUC,
		startUpAt:   time.Now(),
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
	result, err := hdl.inventoryUC.ListProducts(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}

	response.OK(c, ToListProductInfoOpenAPI(result), "List  product successfully!")
}

func (hdl *inventoryHandler) InventoryReserve(c *gin.Context) {
	body, err := request.ParseBody[generated.InventoryReserveRequest](c)
	if err != nil {
		response.BadRequest(c, ReservationErrorCode, err)
		return
	}

	req := FromOpenApiInventoryReserveRequest(body)

	reservation, err := hdl.inventoryUC.InventoryReserve(c.Request.Context(), *req)
	if err != nil {
		response.InternalServerError(c, ReservationErrorCode, err)
		return
	}

	response.OK(c, ToReservationResponseOpenAPI(reservation), "inventory reserve successfully!")
}
