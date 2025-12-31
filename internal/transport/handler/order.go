package handler

import (
	"time"

	"github.com/ductran999/letobserv/api/generated"
	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/transport/mapper"
	"github.com/ductran999/letobserv/pkg/response"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUsecase usecase.OrderUseCase
	startUpAt    time.Time
}

func NewOrderHandler(orderUsecase usecase.OrderUseCase) OrderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
		startUpAt:    time.Now(),
	}
}

func (hdl *orderHandler) HealthCheck(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": consts.HealthyState,
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *orderHandler) PlaceOrder(c *gin.Context) {
	body, err := ParseBody[generated.PlaceOrderJSONRequestBody](c)
	if err != nil {
		response.BadRequest(c, consts.BadRequest, err)
		return
	}

	input := mapper.FromPlaceOrderOpenAPIRequest(body)
	placedOrder, err := hdl.orderUsecase.PlacePOrder(c.Request.Context(), *input)
	if err != nil {
		response.InternalServerError(c, consts.InternalServerError, err)
		return
	}

	response.Created(c, mapper.ToPlacedOrderInfoOpenAPI(placedOrder), "Order placed successfully!")
}
