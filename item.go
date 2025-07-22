package dotagiftx

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultItemRarity = "regular"

// Item error types.
const (
	ItemErrNotFound Errors = iota + itemErrorIndex
	ItemErrRequiredID
	ItemErrRequiredFields
	ItemErrCreateItemExists
	ItemErrImport
)

// sets error text definition.
func init() {
	appErrorText[ItemErrNotFound] = "item not found"
	appErrorText[ItemErrRequiredID] = "item id is required"
	appErrorText[ItemErrRequiredFields] = "item fields are required"
	appErrorText[ItemErrCreateItemExists] = "item already exists"
	appErrorText[ItemErrImport] = "item import error"
}

type (
	// ItemStatus represents item status.
	ItemStatus uint

	// Item represents item information.
	Item struct {
		ID           string     `json:"id"           db:"id,omitempty"`
		Slug         string     `json:"slug"         db:"slug,omitempty"        valid:"required"`
		Name         string     `json:"name"         db:"name,omitempty"        valid:"required"`
		Hero         string     `json:"hero"         db:"hero,omitempty"        valid:"required"`
		Image        string     `json:"image"        db:"image,omitempty"`
		Origin       string     `json:"origin"       db:"origin,omitempty"`
		Rarity       string     `json:"rarity"       db:"rarity,omitempty"`
		Contributors []string   `json:"-"            db:"contributors,omitempty"`
		Active       *bool      `json:"active"       db:"active,omitempty"`
		ViewCount    int        `json:"view_count"   db:"view_count,omitempty"`
		CreatedAt    *time.Time `json:"created_at"   db:"created_at,omitempty"`
		UpdatedAt    *time.Time `json:"updated_at"   db:"updated_at,omitempty"`
	}

	// ItemImportResult represents import process result.
	ItemImportResult struct {
		Created int `json:"created"`
		Updated int `json:"updated"`
		Error   int `json:"error"`
		Total   int `json:"total"`
	}

	// ItemService provides access to item service.
	ItemService interface {
		// Items returns a list of items.
		Items(opts FindOpts) ([]Item, *FindMetadata, error)

		// Item returns item details by id.
		Item(id string) (*Item, error)

		// Create saves new item details.
		Create(context.Context, *Item) error

		// Update saves item details changes.
		Update(context.Context, *Item) error

		// Import creates new item from yaml format.
		Import(ctx context.Context, f io.Reader) (ItemImportResult, error)

		// TopOrigins returns a list of top origin/treasure base on view count.
		TopOrigins() ([]string, error)

		// TopHeroes returns a list of top heroes base on view count.
		TopHeroes() ([]string, error)
	}

	// ItemStorage defines operation for item records.
	ItemStorage interface {
		// Find returns a list of items from data store.
		Find(opts FindOpts) ([]Item, error)

		// Count returns number of items from data store.
		Count(FindOpts) (int, error)

		// Get returns item details by id from data store.
		Get(id string) (*Item, error)

		// GetBySlug returns item details slug id from data store.
		GetBySlug(slug string) (*Item, error)

		// Create persists a new item to data store.
		Create(*Item) error

		// Update persists item changes to data store.
		Update(*Item) error

		// IsItemExist returns an error if item already exists by name.
		IsItemExist(name string) error

		// AddViewCount increments item view count to data store.
		AddViewCount(id string) error
	}
)

// CheckCreate validates field on creating new item.
func (i Item) CheckCreate() error {
	// Check the required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	return nil
}

// MakeSlug generates item slug.
func (i Item) MakeSlug() string {
	return makeSlug(i.Name + " " + i.Hero)
}

// IsActive determines item is giftable.
func (i Item) IsActive() bool {
	return *i.Active
}

// SetDefaults sets default values for a new item.
func (i Item) SetDefaults() *Item {
	if i.Rarity == "" {
		i.Rarity = defaultItemRarity
	}

	i.Slug = i.MakeSlug()
	return &i
}

func (i Item) ToCatalog() Catalog {
	return Catalog{
		ID:           i.ID,
		Slug:         i.Slug,
		Name:         i.Name,
		Hero:         i.Hero,
		Image:        i.Image,
		Origin:       i.Origin,
		Rarity:       i.Rarity,
		Contributors: i.Contributors,
		ViewCount:    i.ViewCount,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

// NewItemService returns new Item service.
func NewItemService(allowedDomains []string, is ItemStorage, fm FileManager) ItemService {
	return &itemService{is, fm, allowedDomains}
}

type itemService struct {
	itemStg ItemStorage
	fileMgr FileManager

	allowedDomains []string
}

func (s *itemService) Items(opts FindOpts) ([]Item, *FindMetadata, error) {
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

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *itemService) Item(id string) (*Item, error) {
	return s.itemStg.Get(id)
}

func (s *itemService) TopOrigins() ([]string, error) {
	items, err := s.itemStg.Find(FindOpts{})
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
	items, err := s.itemStg.Find(FindOpts{})
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

func (s *itemService) Create(ctx context.Context, itm *Item) error {
	// TODO check moderator/contributors
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	itm.Contributors = []string{au.UserID}

	itm.Name = strings.TrimSpace(itm.Name)
	itm.Hero = strings.TrimSpace(itm.Hero)
	itm.Rarity = strings.ToLower(itm.Rarity)
	itm = itm.SetDefaults()
	if err := itm.CheckCreate(); err != nil {
		return NewXError(ItemErrRequiredFields, err)
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

	go func() {
		if err := pingGoogleSitemap(); err != nil {
			log.Println("ping google sitemap", err)
		}
	}()

	return s.itemStg.Create(itm)
}

func (s *itemService) Update(ctx context.Context, itm *Item) error {
	// TODO check moderator/contributors
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}

	if itm.ID == "" {
		return ItemErrRequiredID
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

func (s *itemService) Import(ctx context.Context, f io.Reader) (ItemImportResult, error) {
	res := ItemImportResult{}

	b, err := io.ReadAll(f)
	if err != nil {
		return res, NewXError(ItemErrImport, err)
	}

	yf := &yamlFile{}
	if err := yaml.Unmarshal(b, yf); err != nil {
		return res, NewXError(ItemErrImport, err)
	}

	res.Total = len(yf.Items)
	for _, ii := range yf.Items {
		itm := &Item{
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

func (s *itemService) getItemByName(name string) (*Item, error) {
	itm, err := s.itemStg.Find(FindOpts{Filter: Item{Name: name}})
	if err != nil {
		return nil, err
	}

	if len(itm) == 0 {
		return nil, ItemErrNotFound
	}

	return &itm[0], nil
}

// downloadItemImage saves an image file from a url.
func (s *itemService) downloadItemImage(baseName, url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	if !slices.Contains(s.allowedDomains, u.Hostname()) {
		return url, fmt.Errorf("item image download for %s is not allowed", url)
	}

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

type yamlFile struct {
	Origin string `yaml:"origin"`
	Items  []struct {
		Name   string `yaml:"name"`
		Hero   string `yaml:"hero"`
		Image  string `yaml:"image"`
		Rarity string `yaml:"rarity"`
	} `yaml:"items"`
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

var slugRE = regexp.MustCompile("[^a-z0-9]+")

// makeSlug creates a URL friendly string base on input.
func makeSlug(s string) string {
	return strings.Trim(slugRE.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
