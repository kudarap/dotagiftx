package fixes

import (
	"context"
	"log"
	"math/rand"
	"strings"

	"github.com/kudarap/dotagiftx/core"
)

func ReIndexAll(
	itemStg core.ItemStorage,
	catalogStg core.CatalogStorage,
) {
	ii, _ := itemStg.Find(core.FindOpts{})
	for _, item := range ii {
		if _, err := catalogStg.Index(item.ID); err != nil {
			log.Println("err", err)
		}

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

func ImageFileExtJPG(itemStg core.ItemStorage, userStg core.UserStorage) {
	const suf = ".jpe"
	const newSuf = ".jpg"

	items, _ := itemStg.Find(core.FindOpts{})
	for _, ii := range items {
		if strings.HasSuffix(ii.Image, suf) {
			ii.Image = strings.Replace(ii.Image, suf, newSuf, 1)
			log.Println("new item image", ii.Image)
			itemStg.Update(&ii)
		}
	}

	users, _ := userStg.Find(core.FindOpts{})
	for _, uu := range users {
		if strings.HasSuffix(uu.Avatar, suf) {
			uu.Avatar = strings.Replace(uu.Avatar, suf, newSuf, 1)
			log.Println("new user avatr", uu.Avatar)
			userStg.Update(&uu)
		}
	}
}
