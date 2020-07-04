package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
	"gopkg.in/yaml.v3"
)

// NewItem returns new Item service.
func NewItem(is core.ItemStorage) core.ItemService {
	return &itemService{is}
}

type itemService struct {
	itemStg core.ItemStorage
}

func (s *itemService) Items(opts core.FindOpts) ([]core.Item, *core.FindMetadata, error) {
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

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *itemService) Item(id string) (*core.Item, error) {
	return s.itemStg.Get(id)
}

func (s *itemService) Create(ctx context.Context, itm *core.Item) error {
	// TODO check moderator/contributors
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}
	itm.Contributors = []string{au.UserID}

	itm.Name = strings.TrimSpace(itm.Name)
	itm.Hero = strings.TrimSpace(itm.Hero)
	itm.Rarity = strings.ToLower(itm.Rarity)
	itm = itm.SetDefaults()
	if err := itm.CheckCreate(); err != nil {
		return errors.New(core.ItemErrRequiredFields, err)
	}

	if err := s.itemStg.IsItemExist(itm.Name); err != nil {
		return err
	}

	return s.itemStg.Create(itm)
}

func (s *itemService) Update(ctx context.Context, it *core.Item) error {
	panic("implement me")
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

func (s *itemService) Import(ctx context.Context, f io.Reader) error {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New(core.ItemErrImport, err)
	}

	yf := &yamlFile{}
	if err := yaml.Unmarshal(b, yf); err != nil {
		return errors.New(core.ItemErrImport, err)
	}

	fmt.Println(yf)

	return nil
}
