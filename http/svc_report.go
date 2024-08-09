package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	dgx "github.com/kudarap/dotagiftx"
)

func handleReportList(svc dgx.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &dgx.Report{})
		if err != nil {
			respondError(w, err)
			return
		}

		list, md, err := svc.Reports(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dgx.Report{}
		}

		o := newDataWithMeta(list, md)
		respondOK(w, o)
	}
}

func handleReportDetail(svc dgx.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep, err := svc.Report(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, rep)
	}
}

func handleReportCreate(svc dgx.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep := new(dgx.Report)
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
