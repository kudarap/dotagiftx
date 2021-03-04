package service

import (
	"context"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewReport returns new Report service.
func NewReport(rs core.ReportStorage) core.ReportService {
	return &reportService{rs}
}

type reportService struct {
	reportStg core.ReportStorage
}

func (s *reportService) Reports(opts core.FindOpts) ([]core.Report, *core.FindMetadata, error) {
	res, err := s.reportStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.reportStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *reportService) Report(id string) (*core.Report, error) {
	return s.reportStg.Get(id)
}

func (s *reportService) CreateSurvey(ctx context.Context, rep *core.Report) error {
	rep.Type = core.ReportTypeSurvey
	return s.Create(ctx, rep)
}

func (s *reportService) Create(ctx context.Context, rep *core.Report) error {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}
	rep.UserID = au.UserID

	rep.Label = strings.TrimSpace(rep.Label)
	rep.Text = strings.TrimSpace(rep.Text)
	if err := rep.CheckCreate(); err != nil {
		return errors.New(core.ReportErrRequiredFields, err)
	}

	return s.reportStg.Create(rep)
}
