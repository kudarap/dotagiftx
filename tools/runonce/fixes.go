package runonce

import (
	"log"
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
