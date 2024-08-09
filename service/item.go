package service

import (
	"context"
	"io"
	"net/http"
	"sort"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	"gopkg.in/yaml.v3"
)

// NewItem returns new Item service.
func NewItem(is dgx.ItemStorage, fm dgx.FileManager) dgx.ItemService {
	return &itemService{is, fm}
}

type itemService struct {
	itemStg dgx.ItemStorage
	fileMgr dgx.FileManager
}

func (s *itemService) Items(opts dgx.FindOpts) ([]dgx.Item, *dgx.FindMetadata, error) {
	res, err := s.itemStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.itemStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &dgx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *itemService) Item(id string) (*dgx.Item, error) {
	return s.itemStg.Get(id)
}

func (s *itemService) TopOrigins() ([]string, error) {
	items, err := s.itemStg.Find(dgx.FindOpts{})
	if err != nil {
		return nil, err
	}

	col := map[string]int{}
	for _, ii := range items {
		col[ii.Origin] += ii.ViewCount
	}

	var pt []string
	pt = append(pt, sortedKeys(col)...)

	return pt, nil
}

func (s *itemService) TopHeroes() ([]string, error) {
	items, err := s.itemStg.Find(dgx.FindOpts{})
	if err != nil {
		return nil, err
	}

	col := map[string]int{}
	for _, ii := range items {
		col[ii.Hero] += ii.ViewCount
	}

	var ph []string
	for _, s := range sortedKeys(col) {
		ph = append(ph, s)
	}

	return ph, nil
}

func (s *itemService) Create(ctx context.Context, itm *dgx.Item) error {
	// TODO check moderator/contributors
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}
	itm.Contributors = []string{au.UserID}

	itm.Name = strings.TrimSpace(itm.Name)
	itm.Hero = strings.TrimSpace(itm.Hero)
	itm.Rarity = strings.ToLower(itm.Rarity)
	itm = itm.SetDefaults()
	if err := itm.CheckCreate(); err != nil {
		return errors.New(dgx.ItemErrRequiredFields, err)
	}

	if err := s.itemStg.IsItemExist(itm.Name); err != nil {
		return err
	}

	// Download image when available.
	if itm.Image != "" {
		img, err := s.downloadItemImage(itm.MakeSlug(), itm.Image)
		if err != nil {
			return err
		}
		itm.Image = img
	}

	go pingGoogleSitemap()

	return s.itemStg.Create(itm)
}

func (s *itemService) Update(ctx context.Context, itm *dgx.Item) error {
	// TODO check moderator/contributors
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}

	if itm.ID == "" {
		return dgx.ItemErrRequiredID
	}

	itm.Name = strings.TrimSpace(itm.Name)
	itm.Hero = strings.TrimSpace(itm.Hero)
	itm.Rarity = strings.ToLower(itm.Rarity)

	// Download image when available.
	if itm.Image != "" {
		img, err := s.downloadItemImage(itm.MakeSlug(), itm.Image)
		if err != nil {
			return err
		}
		itm.Image = img
	}

	return s.itemStg.Update(itm)
}

type yamlFile struct {
	Origin string `yaml:"origin"`
	Items  []struct {
		Name   string `yaml:"name"`
		Hero   string `yaml:"hero"`
		Image  string `yaml:"image"`
		Rarity string `yaml:"rarity"`
	} `yaml:"items"`
}

func (s *itemService) Import(ctx context.Context, f io.Reader) (dgx.ItemImportResult, error) {
	res := dgx.ItemImportResult{}

	b, err := io.ReadAll(f)
	if err != nil {
		return res, errors.New(dgx.ItemErrImport, err)
	}

	yf := &yamlFile{}
	if err := yaml.Unmarshal(b, yf); err != nil {
		return res, errors.New(dgx.ItemErrImport, err)
	}

	res.Total = len(yf.Items)
	for _, ii := range yf.Items {
		itm := &dgx.Item{
			Origin: yf.Origin,
			Name:   ii.Name,
			Hero:   ii.Hero,
			Image:  ii.Image,
			Rarity: ii.Rarity,
		}

		// Update current item if exists.
		if cur, _ := s.getItemByName(ii.Name); cur != nil {
			itm.ID = cur.ID
			if err := s.Update(ctx, itm); err != nil {
				res.Error++
				continue
			}
			res.Updated++
			continue
		}

		if err := s.Create(ctx, itm); err != nil {
			res.Error++
			continue
		}
		res.Created++
	}

	return res, nil
}

func (s *itemService) getItemByName(name string) (*dgx.Item, error) {
	itm, err := s.itemStg.Find(dgx.FindOpts{Filter: dgx.Item{Name: name}})
	if err != nil {
		return nil, err
	}

	if len(itm) == 0 {
		return nil, dgx.ItemErrNotFound
	}

	return &itm[0], nil
}

// downloadItemImage saves image file from a url.
func (s *itemService) downloadItemImage(baseName, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	n, err := s.fileMgr.SaveWithName(resp.Body, baseName)
	if err != nil {
		return "", err
	}

	return n, nil
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

// sortedKeys sorts map string by int value.
func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
