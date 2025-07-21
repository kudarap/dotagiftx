package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx"
)

var json = jsoniter.ConfigFastest

type httpMsg struct {
	Error bool   `json:"error,omitempty"`
	Typ   string `json:"type,omitempty"`
	Msg   string `json:"msg"`
}

func newMsg(msg string) httpMsg {
	m := httpMsg{}
	m.Msg = msg
	return m
}

func newError(err error) interface{} {
	m := httpMsg{}
	m.Error = true
	m.Msg = err.Error()
	return m
}

type dataWithMeta struct {
	Data        interface{} `json:"data"`
	ResultCount int         `json:"result_count"`
	TotalCount  int         `json:"total_count"`
}

func newDataWithMeta(data interface{}, md *dotagiftx.FindMetadata) dataWithMeta {
	return dataWithMeta{data, md.ResultCount, md.TotalCount}
}

func hasQueryField(url *url.URL, key string) bool {
	_, ok := url.Query()[key]
	return ok
}

func respond(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Check for string response body.
	if s, ok := body.(string); ok {
		_, _ = w.Write([]byte(s))
		return
	}

	// Generate the response
	enc := json.NewEncoder(w)
	if err := enc.Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("could not encode body into JSON: %s", err)))
	}
}

func respondOK(w http.ResponseWriter, body interface{}) {
	respond(w, http.StatusOK, body)
}

func respondError(w http.ResponseWriter, err error) {
	var body interface{}
	status := http.StatusBadRequest

	// Try to parse handled errors.
	xerr, ok := parseXError(err)
	if ok {
		if xerr.Fatal {
			status = http.StatusInternalServerError
		} else if errors.Is(xerr.Type, dotagiftx.AuthErrNoAccess) {
			status = http.StatusUnauthorized
		} else if errors.Is(xerr.Type, dotagiftx.AuthErrForbidden) {
			status = http.StatusForbidden
		}

		body = httpMsg{true, xerr.Type.String(), err.Error()}
	} else {
		// Use a generic error message
		body = newError(err)
	}

	respond(w, status, body)
}

func handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusNotFound, newError(fmt.Errorf("resource not found")))
	}
}

func handle405() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusMethodNotAllowed, newError(fmt.Errorf("method not allowed")))
	}
}

func parseForm(r *http.Request, form interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		return fmt.Errorf("could not parse json form: %s", err)
	}

	return nil
}

// parseXError returns Errors value if available, else returns nil and ok is false.
// When error is a core.Error type will create new error with that type to handle them gracefully.
// Useful when checking errors types on Parse().
func parseXError(err error) (e *dotagiftx.XErrors, ok bool) {
	// Try packaged error assertion.
	if errors.As(err, &e) {
		return
	}

	// Try core error assertion as type and handles unpackaged error with valid type that can be used to
	// check typed errors.
	var t dotagiftx.Errors
	if errors.As(err, &t) {
		// Error with no details.
		return &dotagiftx.XErrors{Type: t, Err: errors.New("")}, true
	}

	// Cant parse the error.
	return nil, false
}
