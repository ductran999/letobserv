package orderhttp

import (
	"time"

	"github.com/ductran999/letobserv/pkg/errs"
	"github.com/ductran999/letobserv/pkg/request"
	"github.com/ductran999/letobserv/pkg/response"
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	order "github.com/ductran999/letobserv/services/placement/internal/modules/order/domain"
	"github.com/gin-gonic/gin"
)

const (
	HealthyState = "heathy"
)

type OrderHandler interface {
	PlaceOrder(c *gin.Context)
	GetOrder(c *gin.Context, id string)
}

type orderHandler struct {
	orderUsecase order.UseCase
	startUpAt    time.Time
}

func NewOrderHandler(orderUsecase order.UseCase) OrderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
		startUpAt:    time.Now(),
	}
}

func (hdl *orderHandler) HealthCheck(c *gin.Context) {
	uptime := int64(time.Since(hdl.startUpAt).Seconds())
	resp := gin.H{
		"status": HealthyState,
		"uptime": uptime,
	}
	response.OK(c, resp, "OK")
}

func (hdl *orderHandler) PlaceOrder(c *gin.Context) {
	body, err := request.ParseBody[gen.PlaceOrderJSONRequestBody](c)
	if err != nil {
		response.BadRequest(c, errs.BadRequest, err)
		return
	}

	input := toPlaceOrderRequest(body)
	placedOrder, err := hdl.orderUsecase.PlacePOrder(c.Request.Context(), *input)
	if err != nil {
		response.InternalServerError(c, errs.InternalServerError, err)
		return
	}

	response.Created(c, toPlacedOrderInfoOpenAPI(placedOrder), "Order placed successfully!")
}

func (hdl *orderHandler) GetOrder(c *gin.Context, id string) {
	result, err := hdl.orderUsecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, errs.InternalServerError, err)
		return
	}

	response.OK(c, toPlacedOrderDetailsOpenAPI(result), "get detail placed order successfully!")
}
