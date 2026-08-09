package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/dotagiftx"
)

// reportService provides access to report service methods used by http handlers.
type reportService interface {
	// Reports returns a list of reports.
	Reports(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Report, *dotagiftx.FindMetadata, error)

	// Report returns report details by id.
	Report(ctx context.Context, id string) (*dotagiftx.Report, error)

	// Create saves new report details.
	Create(context.Context, *dotagiftx.Report) error
}

func handleReportList(svc reportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &dotagiftx.Report{})
		if err != nil {
			respondError(w, err)
			return
		}

		list, md, err := svc.Reports(r.Context(), opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dotagiftx.Report{}
		}

		o := newDataWithMeta(list, md)
		respondOK(w, o)
	}
}

func handleReportDetail(svc reportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep, err := svc.Report(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, rep)
	}
}

func handleReportCreate(svc reportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep := new(dotagiftx.Report)
		if err := parseForm(r, rep); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), rep); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, rep)
	}
}
