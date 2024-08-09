package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

// NewReport returns new Report service.
func NewReport(rs dgx.ReportStorage, wp webhookPoster) dgx.ReportService {
	return &reportService{rs, wp}
}

type reportService struct {
	reportStg     dgx.ReportStorage
	webhookPoster webhookPoster
}

func (s *reportService) Reports(opts dgx.FindOpts) ([]dgx.Report, *dgx.FindMetadata, error) {
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

	return res, &dgx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *reportService) Report(id string) (*dgx.Report, error) {
	return s.reportStg.Get(id)
}

func (s *reportService) CreateSurvey(ctx context.Context, rep *dgx.Report) error {
	rep.Type = dgx.ReportTypeSurvey
	return s.Create(ctx, rep)
}

func (s *reportService) Create(ctx context.Context, rep *dgx.Report) error {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}
	rep.UserID = au.UserID

	rep.Label = strings.TrimSpace(rep.Label)
	rep.Text = strings.TrimSpace(rep.Text)
	if err := rep.CheckCreate(); err != nil {
		return errors.New(dgx.ReportErrRequiredFields, err)
	}

	if err := s.reportStg.Create(rep); err != nil {
		return err
	}

	go func() {
		if err := s.shootToDiscord(rep.ID); err != nil {
			log.Println("could not shoot to discord:", err)
		}
	}()

	return nil
}

func (s *reportService) shootToDiscord(reportID string) error {
	reps, _, err := s.Reports(dgx.FindOpts{Filter: dgx.Report{ID: reportID}})
	if err != nil {
		return err
	}
	if len(reps) == 0 {
		return nil
	}

	rep := reps[0]
	username := fmt.Sprintf("%s (%s)", rep.User.Name, rep.User.SteamID)
	content := fmt.Sprintf("[%s] %s", rep.Type, rep.Text)
	if err = s.webhookPoster.PostWebhook(username, content); err != nil {
		return err
	}

	return nil
}

type webhookPoster interface {
	PostWebhook(username, content string) error
}
