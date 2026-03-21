package inventory

import (
	"time"

	"github.com/ductran999/letobserv/pkg/request"
	"github.com/ductran999/letobserv/pkg/response"
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventoryuc "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/usecase"
	"github.com/gin-gonic/gin"
)

const (
	ReservationErrorCode = "RESERVED_FAILED"
)

type inventoryHandler struct {
	inventoryUC inventoryuc.InventoryUsecase
	startUpAt   time.Time
}

func NewInventoryHandler(uc inventoryuc.InventoryUsecase) *inventoryHandler {
	if uc == nil {
		panic("inventory cannot be nil")
	}

	return &inventoryHandler{
		inventoryUC: uc,
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
	// result, err := hdl.inventoryUC.ListProducts(c.Request.Context())
	// if err != nil {
	// 	_ = c.Error(err)
	// }

	response.OK(c, nil, "List  product successfully!")
}

func (hdl *inventoryHandler) InventoryReserve(c *gin.Context) {
	body, err := request.ParseBody[gen.InventoryReserveRequest](c)
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

func (hdl *inventoryHandler) GetProduct(c *gin.Context, id string) {
	// result, err := hdl.inventoryUC.GetProduct(c.Request.Context(), id)
	// if err != nil {
	// 	response.InternalServerError(c, "VIEW_PRODUCT_ERROR", err)
	// }

	response.OK(c, nil, "inventory reserve successfully!")
}
