package service

import (
	"net/http"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
)

// NewTrack returns new track service.
func NewTrack(ts dgx.TrackStorage, ps dgx.ItemStorage) dgx.TrackService {
	return &trackService{ts, ps}
}

type trackService struct {
	trackStg dgx.TrackStorage
	itemStg  dgx.ItemStorage
}

func (s *trackService) Tracks(opts dgx.FindOpts) ([]dgx.Track, *dgx.FindMetadata, error) {
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

	return res, &dgx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  total,
	}, nil
}

func (s *trackService) Track(id string) (*dgx.Track, error) {
	return s.trackStg.Get(id)
}

func (s *trackService) CreateFromRequest(r *http.Request) error {
	t := new(dgx.Track)
	t.SetDefaults(r)

	// Track post view.
	if t.Type == dgx.TrackTypeView && t.ItemID != "" {
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

	t := new(dgx.Track)
	t.SetDefaults(r)
	t.Type = dgx.TrackTypeSearch
	t.Keyword = keyword
	return s.trackStg.Create(t)
}
