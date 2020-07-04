package service

import (
	"net/http"

	"github.com/kudarap/dota2giftables/core"
)

// NewTrack returns new track service.
func NewTrack(ts core.TrackStorage, ps core.ItemStorage) core.TrackService {
	return &trackService{ts, ps}
}

type trackService struct {
	trackStg core.TrackStorage
	itemStg  core.ItemStorage
}

func (s *trackService) Tracks(opts core.FindOpts) ([]core.Track, *core.FindMetadata, error) {
	res, err := s.trackStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get total count for metadata.
	total, err := s.trackStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  total,
	}, nil
}

func (s *trackService) Track(id string) (*core.Track, error) {
	return s.trackStg.Get(id)
}

func (s *trackService) CreateFromRequest(r *http.Request) error {
	t := new(core.Track)
	t.SetDefaults(r)

	// Track post view.
	if t.Type == core.TrackTypeView && t.ItemID != "" {
		if err := s.itemStg.AddViewCount(t.ItemID); err != nil {
			return err
		}
	}

	return s.trackStg.Create(t)
}

func (s *trackService) CreateSearchKeyword(r *http.Request, keyword string) error {
	t := new(core.Track)
	t.SetDefaults(r)
	t.Type = core.TrackTypeSearch
	t.Keyword = keyword
	return s.trackStg.Create(t)
}
