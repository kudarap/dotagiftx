package fixes

import (
	"math/rand"

	"github.com/kudarap/dota2giftables/core"
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

	ii, _ := itemStg.Find(core.FindOpts{})
	uu, _ := userStg.Find(core.FindOpts{})
	for _, item := range ii {
		for _, user := range uu {
			m := &core.Market{
				ItemID: item.ID,
				Price:  float64(rand.Intn(1000)) / 10,
			}
			ctx := core.AuthToContext(nil, &core.Auth{UserID: user.ID})
			marketSvc.Create(ctx, m)
		}
	}
}
