package inventoryhttp

import (
	"github.com/ductran999/letobserv/pkg/request"
	"github.com/ductran999/letobserv/pkg/response"
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"
	"github.com/gin-gonic/gin"
)

const (
	ReservationErrorCode = "RESERVED_FAILED"
)

type InventoryHandler interface {
	InventoryReserve(c *gin.Context)
}

type handler struct {
	inventoryUC inventory.Usecase
}

func NewHandler(uc inventory.Usecase) InventoryHandler {
	if uc == nil {
		panic("inventory cannot be nil")
	}

	return &handler{inventoryUC: uc}
}

func (hdl *handler) InventoryReserve(c *gin.Context) {
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

	response.OK(c, toReservationResponseOpenAPI(reservation), "inventory reserve successfully!")
}
