package inventory

import (
	generated "github.com/ductran999/letobserv/api/generated/inventory"
	"github.com/ductran999/letobserv/internal/application/usecase/inventory"
	"github.com/ductran999/letobserv/internal/domain/product"
	"github.com/google/uuid"
)

func FromOpenApiInventoryReserveRequest(body *generated.InventoryReserveRequest) *inventory.InventoryReserveInput {
	orderItems := make([]inventory.ReservationItem, 0, len(body.Items))
	for _, item := range body.Items {
		orderItem := inventory.ReservationItem{
			ProductID: item.ProductId.String(),
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	return &inventory.InventoryReserveInput{
		OrderID: body.OrderId.String(),
		Items:   orderItems,
		TTL:     body.TtlSeconds,
	}
}

func ToListProductInfoOpenAPI(output *inventory.ListProductsOutput) []generated.ProductInfo {
	if output == nil {
		return []generated.ProductInfo{}
	}

	resp := make([]generated.ProductInfo, len(output.Products))
	for i, p := range output.Products {
		resp[i] = generated.ProductInfo{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}

	return resp
}

func ToProductInfoOpenAPI(output *product.Product) generated.ProductInfo {
	if output == nil {
		return generated.ProductInfo{}
	}

	resp := generated.ProductInfo{
		Id:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		Price:       output.Price,
	}

	return resp
}

func ToReservationResponseOpenAPI(reservation *inventory.InventoryReservationView) *generated.ReservationResponse {
	reservationItems := make([]generated.ReservationItem, len(reservation.Items))
	for i, rev := range reservation.Items {
		reservationItems[i] = generated.ReservationItem{
			ProductId: uuid.MustParse(rev.ProductID),
			Quantity:  rev.Quantity,
		}
	}

	return &generated.ReservationResponse{
		ReservationId: reservation.ID,
		OrderId:       reservation.OrderID,
		Items:         reservationItems,
		ExpiredAt:     reservation.TTL,
	}
}
