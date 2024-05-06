package service

import (
	"net/http"
	"strings"

	"github.com/kudarap/dotagiftx"
)

// NewTrack returns new track service.
func NewTrack(ts dotagiftx.TrackStorage, ps dotagiftx.ItemStorage) dotagiftx.TrackService {
	return &trackService{ts, ps}
}

type trackService struct {
	trackStg dotagiftx.TrackStorage
	itemStg  dotagiftx.ItemStorage
}

func (s *trackService) Tracks(opts dotagiftx.FindOpts) ([]dotagiftx.Track, *dotagiftx.FindMetadata, error) {
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

	return res, &dotagiftx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  total,
	}, nil
}

func (s *trackService) Track(id string) (*dotagiftx.Track, error) {
	return s.trackStg.Get(id)
}

func (s *trackService) CreateFromRequest(r *http.Request) error {
	t := new(dotagiftx.Track)
	t.SetDefaults(r)

	// Track post view.
	if t.Type == dotagiftx.TrackTypeView && t.ItemID != "" {
		if err := s.itemStg.AddViewCount(t.ItemID); err != nil {
			return err
		}
	}

	return s.trackStg.Create(t)
}

func (s *trackService) CreateSearchKeyword(r *http.Request, keyword string) error {
	if r.Method != http.MethodGet {
		return nil
	}

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	t := new(dotagiftx.Track)
	t.SetDefaults(r)
	t.Type = dotagiftx.TrackTypeSearch
	t.Keyword = keyword
	return s.trackStg.Create(t)
}
