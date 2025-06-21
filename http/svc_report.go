package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx"
)

func handleReportList(svc dotagiftx.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &dotagiftx.Report{})
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
			list = []dotagiftx.Report{}
		}

		o := newDataWithMeta(list, md)
		respondOK(w, o)
	}
}

func handleReportDetail(svc dotagiftx.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep, err := svc.Report(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, rep)
	}
}

func handleReportCreate(svc dotagiftx.ReportService) http.HandlerFunc {
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
