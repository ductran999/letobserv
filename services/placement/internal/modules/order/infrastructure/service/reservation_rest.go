package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ductran999/letobserv/internal/application/port/service"
	"github.com/ductran999/letobserv/internal/application/usecase/order"
	"github.com/ductran999/letobserv/pkg/httpclient"
)

const (
	inventoryReservationPath = "/inventory/reserve"
	productDetailPath        = "/products"
)

type inventoryService struct {
	client              httpclient.Client
	inventoryServiceURL string
}

func NewInventoryService(client httpclient.Client, inventoryServiceURL string) *inventoryService {
	return &inventoryService{
		client:              client,
		inventoryServiceURL: inventoryServiceURL,
	}
}

func (svc *inventoryService) Reserve(ctx context.Context, req service.InventoryReserveRequest) error {
	url := svc.inventoryServiceURL + inventoryReservationPath
	resp, err := svc.client.Post(ctx, url, req, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return order.ErrReserveInventory
	}

	return nil
}

func (svc *inventoryService) GetProduct(
	ctx context.Context,
	id string,
) error {
	baseURL := fmt.Sprintf(
		"%s%s/%s",
		svc.inventoryServiceURL,
		productDetailPath,
		id,
	)

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("fake", "true")
	u.RawQuery = q.Encode()

	resp, err := svc.client.Get(ctx, u.String(), httpclient.ReqHeader{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return order.ErrReserveInventory
	}

	return nil
}
