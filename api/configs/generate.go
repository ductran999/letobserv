package generated

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config orders.server.yaml ../openapi/orders.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config orders.types.yaml ../schemas/orders.yaml

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config inventory.server.yaml ../openapi/inventory.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config inventory.types.yaml ../schemas/inventory.yaml
