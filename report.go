package dgx

import (
	"context"
	"strconv"
	"time"
)

// Report error types.
const (
	ReportErrNotFound Errors = iota + 5000
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

// Report pre-defined labels.
const (
	ReportLabelSurveyNext = "community-whats-next"
)

type (
	// Represents report types.
	ReportType uint

	// Report represents feedback from user or system that can
	// be use on survey and bug reporting.
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

	// ReportService provides access to report service.
	ReportService interface {
		// Reports returns a list of reports.
		Reports(opts FindOpts) ([]Report, *FindMetadata, error)

		// Report returns report details by id.
		Report(id string) (*Report, error)

		// Create saves new report details.
		Create(context.Context, *Report) error
	}

	// ReportStorage defines operation for report records.
	ReportStorage interface {
		// Find returns a list of reports from data store.
		Find(opts FindOpts) ([]Report, error)

		// Count returns number of reports from data store.
		Count(FindOpts) (int, error)

		// Get returns report details by id from data store.
		Get(id string) (*Report, error)

		// Create persists a new report to data store.
		Create(*Report) error
	}
)

// CheckCreate validates field on creating new report.
func (r Report) CheckCreate() error {
	// Check required fields.
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
