package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	if err := s.reportStg.Create(rep); err != nil {
		return err
	}

	go func() {
		if err := s.shootToDiscord(rep.ID); err != nil {
			fmt.Println("could not shoot to discord:", err)
		}
	}()

	return nil
}

const discordURL = "https://discord.com/api/webhooks/856275008867008523/hS3jT4bUyoJbtBMZq106QK24sM2L54Xvyyz1M_hExOu-tQeKyZjmbNIWteg-Yg2sTfvU"

func (s *reportService) shootToDiscord(reportID string) error {
	reps, _, err := s.Reports(core.FindOpts{Filter: core.Report{ID: reportID}})
	if err != nil {
		return err
	}
	if len(reps) == 0 {
		return nil
	}
	rep := reps[0]

	payload := struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}{
		fmt.Sprintf("%s (%s)", rep.User.Name, rep.User.SteamID),
		fmt.Sprintf("[%s] %s", rep.Type, rep.Text),
	}

	b := new(bytes.Buffer)
	if err = json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}
	resp, err := http.Post(discordURL, "application/json", b)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
