package service

import (
	"context"

	"github.com/kudarap/dota2giftables/core"
)

// NewPost returns new Item service.
func NewPost(ps core.ItemStorage, us core.UserStorage, fm core.FileManager) core.ItemService {
	return &itemService{ps, us, fm}
}

type itemService struct {
	itemStg core.ItemStorage
	userStg core.UserStorage
	fileMgr core.FileManager
}

func (i *itemService) Items(opts core.FindOpts) ([]core.Item, error) {
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
