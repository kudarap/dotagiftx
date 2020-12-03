package fixes

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steam"
)

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
	for _, m := range append(rsvd, sold...) {
		s := extractProfURLFromNotes(m.Notes)
		if s == "" {
			continue
		}

		// Trim notes line to get only the url.
		vURL := strings.Split(s, " ")[0]
		notes := strings.ReplaceAll(m.Notes, vURL, "")
		notes = strings.Trim(notes, "\n")
		notes = strings.TrimSpace(notes)

		steamID, err := steamClient.ResolveVanityURL(vURL)
		if err != nil {
			fmt.Println("ERR! could not resolve:", err)
			return
		}

		// Queue for market update.
		updates = append(updates, core.Market{
			ID:             m.ID,
			PartnerSteamID: steamID,
			Notes:          notes,
		})

		fmt.Println(m.ID)
		fmt.Println(vURL)
		fmt.Println("added for update!")
		fmt.Println(strings.Repeat("-", 100))
	}

	fmt.Println("starting update...")
	for _, m := range updates {
		if err := marketStg.Update(&m); err != nil {
			fmt.Println("ERR! could not update:", err)
		}

		fmt.Println("OK", m.ID)
	}
	fmt.Println("update done!")
}

func extractProfURLFromNotes(notes string) (url string) {
	if notes == "" {
		return
	}

	for _, n := range strings.Split(notes, "\n") {
		if strings.HasPrefix(n, steam.VanityPrefixID) {
			return n
		} else if strings.HasPrefix(n, steam.VanityPrefixProfile) {
			return n
		}
	}

	return
}
