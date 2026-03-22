package producthttp

import (
	"github.com/ductran999/letobserv/pkg/response"
	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	ListProducts(c *gin.Context)
	GetProduct(c *gin.Context, id string)
}

type handler struct {
	productUC product.Usecase
}

func New(pu product.Usecase) ProductHandler {
	if pu == nil {
		panic("product usecase is nil")
	}

	return &handler{
		productUC: pu,
	}
}

func (hdl *handler) ListProducts(c *gin.Context) {
	products, err := hdl.productUC.List(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}

	response.OK(c, ToListProductInfoOpenAPI(products), "List  product successfully!")
}

func (hdl *handler) GetProduct(c *gin.Context, id string) {
	// result, err := hdl.inventoryUC.GetProduct(c.Request.Context(), id)
	// if err != nil {
	// 	response.InternalServerError(c, "VIEW_PRODUCT_ERROR", err)
	// }

	response.OK(c, nil, "inventory reserve successfully!")
}
