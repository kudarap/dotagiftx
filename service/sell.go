package service

import (
	"context"

	"github.com/kudarap/dota2giftables/core"
)

// NewSell returns new Sell service.
func NewSell(is core.ItemStorage, us core.UserStorage, fm core.FileManager) core.SellService {
	return &sellService{is, us, fm}
}

type sellService struct {
	itemStg core.ItemStorage
	userStg core.UserStorage
	fileMgr core.FileManager
}

func (s *sellService) Sells(opts core.FindOpts) ([]core.Sell, core.FindMetadata, error) {
	panic("implement me")
}

func (s *sellService) Sell(id string) (*core.Sell, error) {
	panic("implement me")
}

func (s *sellService) Create(ctx context.Context, sell *core.Sell) error {
	panic("implement me")
}

func (s *sellService) Update(ctx context.Context, sell *core.Sell) error {
	panic("implement me")
}
