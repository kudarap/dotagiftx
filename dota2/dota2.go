package dota2

import (
	"errors"
	"time"
)

var (
	treasuresIndex map[string]int
	heroesIndex    map[string]int
)

func init() {
	treasuresIndex = make(map[string]int)
	heroesIndex = make(map[string]int)
	for i, v := range AllTreasures {
		treasuresIndex[v.Slug] = i
	}
	for i, v := range AllHeroes {
		heroesIndex[v.ID] = i
	}
}

func TreasureDetail(slug string) (*Treasure, error) {
	idx, ok := treasuresIndex[slug]
	if !ok {
		return nil, errors.New("treasure not found")
	}
	return &AllTreasures[idx], nil
}

func HeroDetail(id string) (*Hero, error) {
	idx, ok := heroesIndex[id]
	if !ok {
		return nil, errors.New("hero not found")
	}
	return &AllHeroes[idx], nil
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
