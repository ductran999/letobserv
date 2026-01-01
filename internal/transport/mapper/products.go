package mapper

import (
	generated "github.com/ductran999/letobserv/api/generated/inventory"
	"github.com/ductran999/letobserv/internal/application/outputs"
)

// --------- Map Handler to OpenAPI struct ------------------

func ToListProductInfoOpenAPI(output *outputs.ListProductsOutput) []generated.ProductInfo {
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
