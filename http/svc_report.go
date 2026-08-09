package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
)

// reportService provides access to report service methods used by http handlers.
type reportService interface {
	// Reports returns a list of reports.
	Reports(ctx context.Context, opts dotagiftx2.FindOpts) ([]dotagiftx2.Report, *dotagiftx2.FindMetadata, error)

	// Report returns report details by id.
	Report(ctx context.Context, id string) (*dotagiftx2.Report, error)

	// Create saves new report details.
	Create(context.Context, *dotagiftx2.Report) error
}

func handleReportList(svc reportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &dotagiftx2.Report{})
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
			list = []dotagiftx2.Report{}
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
		rep := new(dotagiftx2.Report)
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
