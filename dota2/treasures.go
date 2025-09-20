package dota2

import (
	"time"
)

type Treasure struct {
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Rarity      string     `json:"rarity"`
	Items       int        `json:"items"`
	ReleaseDate *time.Time `json:"release_date"`
}

func releaseDate(v string) *time.Time {
	if v == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}
	return &t
}

var AllTreasures = []Treasure{
	{
		"cosmic-2025-heroes-hoard",
		"Cosmic 2025 Heroes' Hoard",
		"cosmic_2025_heroes_hoard.png",
		"mythical",
		17,
		releaseDate("2025-09-17"),
	},
	{
		"spring-2024-heroes-hoard",
		"Spring 2025 Heroes' Hoard",
		"spring_2024_heroes_hoard.png",
		"mythical",
		16,
		releaseDate("2025-04-23"),
	},
	{
		"the-charms-of-the-snake",
		"The Charms of the Snake",
		"the_charms_of_the_snake.png",
		"mythical",
		9,
		nil,
	},
	{
		"winter-2024-heroes-hoard",
		"Winter 2024 Heroes' Hoard",
		"winter_2024_heroes_hoard.png",
		"mythical",
		17,
		nil,
	},
	{
		"crownfall-2024-collects-cache-ii",
		"Crownfall 2024 Collector's Cache II",
		"crownfall_2024_collect_s_cache_ii.png",
		"mythical",
		16,
		nil,
	},
	{
		"crownfall-2024-collects-cache",
		"Crownfall 2024 Collector's Cache",
		"crownfall_2024_collect_s_cache.png",
		"mythical",
		16,
		nil,
	},
	{
		"crownfall-treasure-i",
		"Crownfall Treasure I",
		"crownfall_treasure_i.png",
		"mythical",
		12,
		nil,
	},
	{
		"crownfall-treasure-ii",
		"Crownfall Treasure II",
		"crownfall_treasure_ii.png",
		"mythical",
		12,
		nil,
	},
	{
		"crownfall-treasure-iii",
		"Crownfall Treasure III",
		"crownfall_treasure_iii.png",
		"mythical",
		11,
		nil,
	},
	{
		"dead-reckoning-chest",
		"Dead Reckoning Chest",
		"dead_reckoning_chest.png",
		"mythical",
		14,
		nil,
	},
	{
		"august-2023-collectors-cache",
		"August 2023 Collector's Cache",
		"august_2023_collector_s_cache.png",
		"mythical",
		16,
		nil,
	},
	{
		"diretide-2022-collectors-cache-ii",
		"Diretide 2022 Collector's Cache II",
		"diretide_2022_collector_s_cache_ii.png",
		"immortal",
		17,
		nil,
	},
	{
		"diretide-2022-collectors-cache",
		"Diretide 2022 Collector's Cache",
		"diretide_2022_collector_s_cache.png",
		"immortal",
		18,
		nil,
	},
	{
		"immortal-treasure-i-2022",
		"Immortal Treasure I 2022",
		"immortal_treasure_i_2022.png",
		"immortal",
		9,
		nil,
	},
	{
		"immortal-treasure-ii-2022",
		"Immortal Treasure II 2022",
		"immortal_treasure_ii_2022.png",
		"immortal",
		9,
		nil,
	},
	{
		"the-battle-pass-collection-2022",
		"The Battle Pass Collection 2022",
		"the_battle_pass_collection_2022.png",
		"immortal",
		8,
		nil,
	},
	{
		"ageless-heirlooms-2022",
		"Ageless Heirlooms 2022",
		"ageless_heirlooms_2022.png",
		"immortal",
		10,
		nil,
	},
	{
		"aghanims-2021-collectors-cache",
		"Aghanim's 2021 Collector's Cache",
		"aghanim_s_2021_collector_s_cache.webp",
		"mythical",
		17,
		nil,
	},
	{
		"aghanims-2021-ageless-heirlooms",
		"Aghanim's 2021 Ageless Heirlooms",
		"aghanim_s_2021_ageless_heirlooms.webp",
		"mythical",
		10,
		nil,
	},
	{
		"aghanims-2021-continuum-collection",
		"Aghanim's 2021 Continuum Collection",
		"aghanim_s_2021_continuum_collection.webp",
		"mythical",
		7,
		nil,
	},
	{
		"aghanims-2021-immortal-treasure",
		"Aghanim's 2021 Immortal Treasure",
		"aghanim_s_2021_immortal_treasure.webp",
		"immortal",
		9,
		nil,
	},
	{
		"nemestice-2021-collectors-cache",
		"Nemestice 2021 Collector's Cache",
		"nemestice_2021_collector_s_cache.webp",
		"mythical",
		15,
		nil,
	},
	{
		"nemestice-2021-immortal-treasure",
		"Nemestice 2021 Immortal Treasure",
		"nemestice_2021_immortal_treasure.webp",
		"mythical",
		9,
		nil,
	},
	{
		"nemestice-2021-themed-treasure",
		"Nemestice 2021 Themed Treasure",
		"nemestice_2021_themed_treasure.webp",
		"mythical",
		11,
		nil,
	},
	{
		"immortal-treasure-i-2020",
		"Immortal Treasure I 2020",
		"immortal_treasure_i_2020.webp",
		"immortal",
		10,
		nil,
	},
	{
		"immortal-treasure-ii-2020",
		"Immortal Treasure II 2020",
		"immortal_treasure_ii_2020.webp",
		"immortal",
		10,
		nil,
	},
	{
		"immortal-treasure-iii-2020",
		"Immortal Treasure III 2020",
		"immortal_treasure_iii_2020.webp",
		"immortal",
		8,
		nil,
	},
	{
		"the-international-2020-collectors-cache",
		"The International 2020 Collector's Cache",
		"the_international_2020_collector_s_cache.webp",
		"mythical",
		18,
		nil,
	},
	{
		"the-international-2020-collectors-cache-ii",
		"The International 2020 Collector's Cache II",
		"the_international_2020_collector_s_cache_ii.webp",
		"mythical",
		17,
		nil,
	},
	{
		"the-international-2019-collectors-cache",
		"The International 2019 Collector's Cache",
		"the_international_2019_collector_s_cache.webp",
		"mythical",
		18,
		nil,
	},
	{
		"the-international-2019-collectors-cache-ii",
		"The International 2019 Collector's Cache II",
		"the_international_2019_collector_s_cache_ii.webp",
		"mythical",
		16,
		nil,
	},
	{
		"the-international-2018-collectors-cache",
		"The International 2018 Collector's Cache",
		"the_international_2018_collector_s_cache.webp",
		"mythical",
		17,
		nil,
	},
	{
		"the-international-2018-collectors-cache-ii",
		"The International 2018 Collector's Cache II",
		"the_international_2018_collector_s_cache_ii.webp",
		"mythical",
		14,
		nil,
	},
	{
		"the-international-2017-collectors-cache",
		"The International 2017 Collector's Cache",
		"the_international_2017_collector_s_cache.webp",
		"mythical",
		22,
		nil,
	},
	{
		"the-international-2016-collectors-cache",
		"The International 2016 Collector's Cache",
		"the_international_2016_collector_s_cache.webp",
		"mythical",
		15,
		nil,
	},
	{
		"the-international-2015-collectors-cache",
		"The International 2015 Collector's Cache",
		"the_international_2015_collector_s_cache.webp",
		"mythical",
		11,
		nil,
	},
	{
		"treasure-of-the-cryptic-beacon",
		"Treasure of the Cryptic Beacon",
		"treasure_of_the_cryptic_beacon.webp",
		"mythical",
		6,
		nil,
	},
}
