package fixes

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steam"
)

func MarketSetRankingScores(userSvc core.UserService, marketSvc core.MarketService) {
	users, err := userSvc.Users(core.FindOpts{})
	if err != nil {
		panic(err)
	}

	for _, uu := range users {
		if err = marketSvc.UpdateUserRankScore(uu.ID); err != nil {
			fmt.Println("market user score error:", err)
			continue
		}

		fmt.Println("market user score ok!:", uu.ID)
	}

	fmt.Println("market user score done!")
}

func MarketIndexRebuild(marketStg core.MarketStorage) {
	res, _ := marketStg.Find(core.FindOpts{})

	for _, rr := range res {
		if _, err := marketStg.Index(rr.ID); err != nil {
			fmt.Println("market index error:", err)
			continue
		}

		fmt.Println("market index ok!:", rr.ID)
	}

	fmt.Println("market index done!")
}

// MarketExtractProfileURLFromNotes WARNING! only use these once and just keeping for reference.
func MarketExtractProfileURLFromNotes(
	marketStg core.MarketStorage,
	steamClient core.SteamClient,
) {
	var updates []core.Market

	rsvd, _ := marketStg.Find(core.FindOpts{
		Filter: core.Market{Status: core.MarketStatusReserved},
	})
	sold, _ := marketStg.Find(core.FindOpts{
		Filter: core.Market{Status: core.MarketStatusSold},
	})
	cnld, _ := marketStg.Find(core.FindOpts{
		Filter: core.Market{Status: core.MarketStatusCancelled},
	})

	for _, m := range append(append(rsvd, sold...), cnld...) {
		vURL, newNotes := extractProfURLFromNotes(m.Notes)
		if vURL == "" {
			fmt.Println("Skipped", m.ID)
			continue
		}

		steamID, err := steamClient.ResolveVanityURL(vURL)
		if err != nil {
			fmt.Println("ERR! could not resolve:", err)
			return
		}

		// Queue for market update.
		updates = append(updates, core.Market{
			ID:             m.ID,
			UpdatedAt:      m.UpdatedAt, // Skips the time update
			PartnerSteamID: steamID,
			Notes:          newNotes,
		})

		fmt.Println(m.ID)
		fmt.Println(vURL)
		fmt.Println("added for update!")
		fmt.Println(strings.Repeat("-", 100))
	}

	fmt.Println("starting update...")
	for _, m := range updates {
		if m.Notes == "" {
			m.Notes = "-"
		}

		if err := marketStg.Update(&m); err != nil {
			fmt.Println("ERR! could not update:", err)
		}

		fmt.Println("OK", m.ID)
	}
	fmt.Println("update done!")
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
