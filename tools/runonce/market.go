package runonce

import (
	"context"
	"log"
	"math/rand"
	"strings"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

func MarketSetRankingScores(userSvc dotagiftx.UserService, marketSvc dotagiftx.MarketService) {
	users, err := userSvc.Users(dotagiftx.FindOpts{})
	if err != nil {
		panic(err)
	}

	for _, uu := range users {
		if err = marketSvc.UpdateUserRankScore(uu.ID); err != nil {
			log.Println("market user score error:", err)
			continue
		}

		log.Println("market user score ok!:", uu.ID)
	}

	log.Println("market user score done!")
}

func MarketIndexRebuild(marketStg dotagiftx.MarketStorage) {
	res, _ := marketStg.Find(dotagiftx.FindOpts{
		Filter: dotagiftx.Market{Status: dotagiftx.MarketStatusSold, Type: dotagiftx.MarketTypeAsk},
	})

	for _, rr := range res {
		if _, err := marketStg.Index(rr.ID); err != nil {
			log.Println("market index error:", err)
			continue
		}

		log.Println("market index ok!:", rr.ID)
	}

	log.Println("market index done!")
}

// MarketExtractProfileURLFromNotes WARNING! only use these once and just keeping for reference.
func MarketExtractProfileURLFromNotes(
	marketStg dotagiftx.MarketStorage,
	steamClient dotagiftx.SteamClient,
) {
	var updates []dotagiftx.Market

	rsvd, _ := marketStg.Find(dotagiftx.FindOpts{
		Filter: dotagiftx.Market{Status: dotagiftx.MarketStatusReserved},
	})
	sold, _ := marketStg.Find(dotagiftx.FindOpts{
		Filter: dotagiftx.Market{Status: dotagiftx.MarketStatusSold},
	})
	cnld, _ := marketStg.Find(dotagiftx.FindOpts{
		Filter: dotagiftx.Market{Status: dotagiftx.MarketStatusCancelled},
	})

	for _, m := range append(append(rsvd, sold...), cnld...) {
		vURL, newNotes := extractProfURLFromNotes(m.Notes)
		if vURL == "" {
			log.Println("Skipped", m.ID)
			continue
		}

		steamID, err := steamClient.ResolveVanityURL(vURL)
		if err != nil {
			log.Println("ERR! could not resolve:", err)
			return
		}

		// Queue for market update.
		updates = append(updates, dotagiftx.Market{
			ID:             m.ID,
			UpdatedAt:      m.UpdatedAt, // Skips the time update
			PartnerSteamID: steamID,
			Notes:          newNotes,
		})

		log.Println(m.ID)
		log.Println(vURL)
		log.Println("added for update!")
		log.Println(strings.Repeat("-", 100))
	}

	log.Println("starting update...")
	for _, m := range updates {
		if m.Notes == "" {
			m.Notes = "-"
		}

		if err := marketStg.Update(&m); err != nil {
			log.Println("ERR! could not update:", err)
		}

		log.Println("OK", m.ID)
	}
	log.Println("update done!")
}

func GenerateFakeMarket(
	itemStg dotagiftx.ItemStorage,
	userStg dotagiftx.UserStorage,
	marketSvc dotagiftx.MarketService,
) {

	ctx := context.Background()
	ii, _ := itemStg.Find(dotagiftx.FindOpts{})
	uu, _ := userStg.Find(dotagiftx.FindOpts{})
	for _, item := range ii {
		for _, user := range uu {
			m := &dotagiftx.Market{
				ItemID: item.ID,
				Price:  float64(rand.Intn(1000)) / 10,
			}
			auc := dotagiftx.AuthToContext(ctx, &dotagiftx.Auth{UserID: user.ID})
			marketSvc.Create(auc, m)
		}
	}
}

const lineBreaker = "\n"

func extractProfURLFromNotes(notes string) (url, newNotes string) {
	if notes == "" {
		return
	}

	ss := strings.Split(notes, lineBreaker)
	for i, n := range ss {
		if strings.HasPrefix(n, steam.VanityPrefixID) || strings.HasPrefix(n, steam.VanityPrefixProfile) {
			// Current data shows all line notes start with the url.
			url = strings.Split(n, " ")[0]
			// Removes the url and additional spaces.
			ss[i] = strings.ReplaceAll(n, url, "")
			ss[i] = strings.TrimSpace(ss[i])

			var nn []string
			for _, s := range ss {
				if strings.TrimSpace(s) == "" {
					continue
				}

				nn = append(nn, s)
			}

			return url, strings.Join(nn, lineBreaker)
		}
	}

	return
}
