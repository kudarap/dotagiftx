package fixes

import (
	"context"
	"log"
	"math/rand"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
)

func ReIndexAll(
	itemStg dgx.ItemStorage,
	catalogStg dgx.CatalogStorage,
) {
	ii, _ := itemStg.Find(dgx.FindOpts{})
	for _, item := range ii {
		if _, err := catalogStg.Index(item.ID); err != nil {
			log.Println("err", err)
		}

	}
}

func GenerateFakeMarket(
	itemStg dgx.ItemStorage,
	userStg dgx.UserStorage,
	marketSvc dgx.MarketService,
) {

	ctx := context.Background()
	ii, _ := itemStg.Find(dgx.FindOpts{})
	uu, _ := userStg.Find(dgx.FindOpts{})
	for _, item := range ii {
		for _, user := range uu {
			m := &dgx.Market{
				ItemID: item.ID,
				Price:  float64(rand.Intn(1000)) / 10,
			}
			auc := dgx.AuthToContext(ctx, &dgx.Auth{UserID: user.ID})
			marketSvc.Create(auc, m)
		}
	}
}

func ImageFileExtJPG(itemStg dgx.ItemStorage, userStg dgx.UserStorage) {
	const suf = ".jpe"
	const newSuf = ".jpg"

	items, _ := itemStg.Find(dgx.FindOpts{})
	for _, ii := range items {
		if strings.HasSuffix(ii.Image, suf) {
			ii.Image = strings.Replace(ii.Image, suf, newSuf, 1)
			log.Println("new item image", ii.Image)
			itemStg.Update(&ii)
		}
	}

	users, _ := userStg.Find(dgx.FindOpts{})
	for _, uu := range users {
		if strings.HasSuffix(uu.Avatar, suf) {
			uu.Avatar = strings.Replace(uu.Avatar, suf, newSuf, 1)
			log.Println("new user avatr", uu.Avatar)
			userStg.Update(&uu)
		}
	}
}
