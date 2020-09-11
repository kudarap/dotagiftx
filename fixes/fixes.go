package fixes

import (
	"context"
	"math/rand"

	"github.com/kudarap/dotagiftx/core"
)

func ReIndexAll(
	itemStg core.ItemStorage,
	catalogStg core.CatalogStorage,
) {
	ii, _ := itemStg.Find(core.FindOpts{})
	for _, item := range ii {
		catalogStg.Index(item.ID)
	}
}

func GenerateFakeMarket(
	itemStg core.ItemStorage,
	userStg core.UserStorage,
	marketSvc core.MarketService,
) {

	ctx := context.Background()
	ii, _ := itemStg.Find(core.FindOpts{})
	uu, _ := userStg.Find(core.FindOpts{})
	for _, item := range ii {
		for _, user := range uu {
			m := &core.Market{
				ItemID: item.ID,
				Price:  float64(rand.Intn(1000)) / 10,
			}
			auc := core.AuthToContext(ctx, &core.Auth{UserID: user.ID})
			marketSvc.Create(auc, m)
		}
	}
}
