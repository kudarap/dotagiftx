package http

import (
	"net/url"

	"github.com/kudarap/dota2giftables/core"
)

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

func newDataWithMeta(data interface{}, md *core.FindMetadata) dataWithMeta {
	return dataWithMeta{data, md.ResultCount, md.TotalCount}
}

func hasQueryField(url *url.URL, key string) bool {
	_, ok := url.Query()[key]
	return ok
}
