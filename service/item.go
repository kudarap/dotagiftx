package service

import (
	"context"

	"github.com/kudarap/dota2giftables/core"
)

// NewItem returns new Item service.
func NewItem(ps core.ItemStorage, us core.UserStorage) core.ItemService {
	return &itemService{ps, us}
}

type itemService struct {
	itemStg core.ItemStorage
	userStg core.UserStorage
}

func (i *itemService) Items(opts core.FindOpts) ([]core.Item, core.FindMetadata, error) {
	panic("implement me")
}

func (i *itemService) Item(id string) (*core.Item, error) {
	panic("implement me")
}

func (i *itemService) Create(ctx context.Context, item *core.Item) error {
	panic("implement me")
}

func (i *itemService) Update(ctx context.Context, item *core.Item) error {
	panic("implement me")
}
