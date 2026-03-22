package producthttp

import (
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
)

func ToListProductInfoOpenAPI(output *product.ListProductsOutput) []gen.ProductInfo {
	if len(output.Products) == 0 {
		return []gen.ProductInfo{}
	}

	resp := make([]gen.ProductInfo, len(output.Products))
	for i, p := range output.Products {
		resp[i] = gen.ProductInfo{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}

	return resp
}
