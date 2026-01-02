package service

import (
	"context"

	"github.com/ductran999/letobserv/internal/application/port/service"
	"github.com/ductran999/letobserv/internal/application/usecase/order"
	"github.com/ductran999/letobserv/pkg/httpclient"
)

const (
	inventoryReservationPath = "/inventory/reserve"
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
	defer resp.Body.Close() //nolint

	if resp.StatusCode != 200 {
		return order.ErrReserveInventory
	}

	return nil
}
