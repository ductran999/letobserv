package handler

import (
	"time"

	"github.com/ductran999/letobserv/pkg/request"
	"github.com/ductran999/letobserv/pkg/response"
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	"github.com/ductran999/letobserv/services/placement/internal/consts"
	order "github.com/ductran999/letobserv/services/placement/internal/modules/order/usecase"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUsecase order.OrderUseCase
	startUpAt    time.Time
}

func NewOrderHandler(orderUsecase order.OrderUseCase) *orderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
		startUpAt:    time.Now(),
	}
}

func (hdl *orderHandler) HealthCheck(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": "healthy",
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *orderHandler) PlaceOrder(c *gin.Context) {
	body, err := request.ParseBody[gen.PlaceOrderJSONRequestBody](c)
	if err != nil {
		response.BadRequest(c, consts.BadRequest, err)
		return
	}

	input := FromPlaceOrderOpenAPIRequest(body)
	placedOrder, err := hdl.orderUsecase.PlacePOrder(c.Request.Context(), *input)
	if err != nil {
		response.InternalServerError(c, consts.InternalServerError, err)
		return
	}

	response.Created(c, ToPlacedOrderInfoOpenAPI(placedOrder), "Order placed successfully!")
}

func (hdl *orderHandler) GetOrder(c *gin.Context, id string) {
	result, err := hdl.orderUsecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, consts.InternalServerError, err)
		return
	}

	response.OK(c, result, "get detail placed order successfully!")
}

func ToPlacedOrderDetailsOpenAPI(order order.PlacedOrder) gen.PlacedOrderDetails {
	orderItems := make([]any, len(order.Items))
	for i, item := range order.Items {
		orderItems[i] = item
	}

	return gen.PlacedOrderDetails{
		OrderId:         order.ID,
		Money:           order.Money,
		OrderItems:      orderItems,
		OrderNumber:     order.OrderNumber,
		ShippingAddress: order.ShippingAddress,
		ShippingPhone:   order.ShippingPhone,
		UserId:          order.UserID,
	}
}
