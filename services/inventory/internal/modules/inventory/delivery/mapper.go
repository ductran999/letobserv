package inventory

import (
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventoryuc "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/usecase"
	"github.com/google/uuid"
)

func FromOpenApiInventoryReserveRequest(body *gen.InventoryReserveRequest) *inventoryuc.InventoryReserveInput {
	orderItems := make([]inventoryuc.ReservationItem, 0, len(body.Items))
	for _, item := range body.Items {
		orderItem := inventoryuc.ReservationItem{
			ProductID: item.ProductId.String(),
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	return &inventoryuc.InventoryReserveInput{
		OrderID: body.OrderId.String(),
		Items:   orderItems,
		TTL:     body.TtlSeconds,
	}
}

// func ToListProductInfoOpenAPI(output *inventory.ListProductsOutput) []gen.ProductInfo {
// 	if output == nil {
// 		return []gen.ProductInfo{}
// 	}

// 	resp := make([]gen.ProductInfo, len(output.Products))
// 	for i, p := range output.Products {
// 		resp[i] = gen.ProductInfo{
// 			Id:          p.ID,
// 			Name:        p.Name,
// 			Description: p.Description,
// 			Price:       p.Price,
// 		}
// 	}

// 	return resp
// }

// func ToProductInfoOpenAPI(output *product.Product) gen.ProductInfo {
// 	if output == nil {
// 		return gen.ProductInfo{}
// 	}

// 	resp := gen.ProductInfo{
// 		Id:          output.ID,
// 		Name:        output.Name,
// 		Description: output.Description,
// 		Price:       output.Price,
// 	}

// 	return resp
// }

func ToReservationResponseOpenAPI(reservation *inventoryuc.InventoryReservationView) *gen.ReservationResponse {
	reservationItems := make([]gen.ReservationItem, len(reservation.Items))
	for i, rev := range reservation.Items {
		reservationItems[i] = gen.ReservationItem{
			ProductId: uuid.MustParse(rev.ProductID),
			Quantity:  rev.Quantity,
		}
	}

	return &gen.ReservationResponse{
		ReservationId: reservation.ID,
		OrderId:       reservation.OrderID,
		Items:         reservationItems,
		ExpiredAt:     reservation.TTL,
	}
}
