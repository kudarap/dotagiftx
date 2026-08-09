package dotagiftx

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Report error types.
const (
	ReportErrNotFound Errors = iota + reportErrorIndex
	ReportErrRequiredID
	ReportErrRequiredFields
)

// sets error text definition.
func init() {
	appErrorText[ReportErrNotFound] = "report not found"
	appErrorText[ReportErrRequiredID] = "report id is required"
	appErrorText[ReportErrRequiredFields] = "report fields are required"
}

// Report types.
const (
	ReportTypeFeedback     ReportType = 10
	ReportTypeSurvey       ReportType = 20
	ReportTypeBug          ReportType = 30
	ReportTypeScamAlert    ReportType = 40
	ReportTypeScamIncident ReportType = 50
)

// Report available labels.
const (
	ReportLabelSurveyNext = "community-whats-next"
)

type (
	// ReportType report types.
	ReportType uint

	// Report represents feedback from user or system that can be used on survey and bug reporting.
	Report struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		UserID    string     `json:"user_id"    db:"user_id,omitempty"`
		Type      ReportType `json:"type"       db:"type,omitempty,indexed"   valid:"required"`
		Label     string     `json:"label"      db:"label,omitempty,indexed"`
		Text      string     `json:"text"       db:"text,omitempty"           valid:"required"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"user,omitempty"`
	}

	// reportRepository defines operation for report records.
	reportRepository interface {
		// Find returns a list of reports from the data store.
		Find(ctx context.Context, opts FindOpts) ([]Report, error)

		// Count returns number of reports from data store.
		Count(ctx context.Context, opts FindOpts) (int, error)

		// Get returns report details by id from data store.
		Get(ctx context.Context, id string) (*Report, error)

		// Create persists a new report to data store.
		Create(context.Context, *Report) error
	}
)

// CheckCreate validates field on creating a new report.
func (r Report) CheckCreate() error {
	// Check the required fields.
	if err := validator.Struct(r); err != nil {
		return err
	}

	return nil
}

var ReportTypeTexts = map[ReportType]string{
	ReportTypeFeedback:     "Feedback",
	ReportTypeSurvey:       "Survey",
	ReportTypeBug:          "Bug",
	ReportTypeScamAlert:    "ScamAlert",
	ReportTypeScamIncident: "ScamIncident",
}

func (t ReportType) String() string {
	s, ok := ReportTypeTexts[t]
	if !ok {
		return strconv.Itoa(int(t))
	}

	return s
}

// NewReportService returns new report service.
func NewReportService(rs reportRepository, wp webhookPoster) *ReportService {
	return &ReportService{rs, wp}
}

type ReportService struct {
	reportRepo    reportRepository
	webhookPoster webhookPoster
}

func (s *ReportService) Reports(ctx context.Context, opts FindOpts) ([]Report, *FindMetadata, error) {
	res, err := s.reportRepo.Find(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get a result and total count for metadata.
	tc, err := s.reportRepo.Count(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *ReportService) Report(ctx context.Context, id string) (*Report, error) {
	return s.reportRepo.Get(ctx, id)
}

func (s *ReportService) CreateSurvey(ctx context.Context, rep *Report) error {
	rep.Type = ReportTypeSurvey
	return s.Create(ctx, rep)
}

func (s *ReportService) Create(ctx context.Context, rep *Report) error {
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	rep.UserID = au.UserID

	rep.Label = strings.TrimSpace(rep.Label)
	rep.Text = strings.TrimSpace(rep.Text)
	if err := rep.CheckCreate(); err != nil {
		return NewXError(ReportErrRequiredFields, err)
	}

	if err := s.reportRepo.Create(ctx, rep); err != nil {
		return err
	}

	go func() {
		if err := s.shootToDiscord(ctx, rep.ID); err != nil {
			log.Println("could not shoot to discord:", err)
		}
	}()

	return nil
}

func (s *ReportService) shootToDiscord(ctx context.Context, reportID string) error {
	reps, _, err := s.Reports(ctx, FindOpts{Filter: Report{ID: reportID}})
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
