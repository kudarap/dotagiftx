package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
)

type ApiLog struct {
	Endpoint string
	Status   int
	Stamp    time.Time
	Elapse   time.Duration
}

type Report struct {
	Endpoint  string
	Reqs      int
	Oks       int
	Errors    int
	TrendTime []time.Duration
	HighTime  time.Duration
	AvgTime   time.Duration
	LowTime   time.Duration
}

type Reports []Report

func (e Reports) Len() int {
	return len(e)
}

func (e Reports) Less(i, j int) bool {
	// return e[i].AvgTime < e[j].AvgTime
	return e[i].HighTime < e[j].HighTime
}

func (e Reports) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

var bar = strings.Repeat("-", 50)

func main() {
	log.Println("reading logs...")
	fptr := flag.String("file", "testdata", "log file to parse")
	flag.Parse()
	data, err := ioutil.ReadFile(*fptr)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("parsing logs...")
	logs := parseLog(string(data))
	log.Println("generationg report...")
	reports := generateReport(logs)

	sort.Sort(Reports(reports))
	for _, rr := range reports {
		fmt.Println(bar)
		fmt.Println(rr.Endpoint)
		fmt.Println(bar)
		fmt.Printf("total: %d(%d/%d)\t\thi: %s avg: %s lo: %s\n",
			rr.Reqs, rr.Oks, rr.Errors,
			rr.HighTime, rr.AvgTime, rr.LowTime)

		plotChart(rr.TrendTime)
	}
}

func plotChart(td []time.Duration) {
	var plot []float64
	for _, tt := range td {
		if tt < time.Second {
			continue
		}

		sec := tt.Seconds()
		plot = append(plot, sec)
	}
	if len(plot) == 0 {
		return
	}
	graph := asciigraph.Plot(plot, asciigraph.Height(10), asciigraph.Width(180))
	fmt.Println(graph)
}

func generateReport(logs []ApiLog) []Report {
	reports := map[string]Report{}
	for _, log := range logs {
		report, init := reports[log.Endpoint]
		if !init {
			report = Report{
				Endpoint:  log.Endpoint,
				AvgTime:   log.Elapse,
				HighTime:  log.Elapse,
				LowTime:   log.Elapse,
				TrendTime: []time.Duration{log.Elapse},
			}
		}

		report.Reqs++
		if log.Status == http.StatusOK {
			report.Oks++
		} else if log.Status >= 400 || log.Status <= 500 {
			report.Errors++
		}

		if init {
			report.AvgTime = (report.AvgTime + log.Elapse) / 2

			if log.Elapse > report.HighTime {
				report.HighTime = log.Elapse
			} else if log.Elapse < report.LowTime {
				report.LowTime = log.Elapse
			}

			report.TrendTime = append(report.TrendTime, log.Elapse)
		}

		reports[log.Endpoint] = report
	}

	var list []Report
	for _, rr := range reports {
		list = append(list, rr)
	}

	return list
}

func parseLog(raw string) []ApiLog {
	var logs []ApiLog
	for _, raw := range strings.Split(raw, "\n") {
		var entry ApiLog

		body := strings.Split(raw, `" level=info msg="`)
		if len(body) < 2 {
			continue
		}
		head, tail := body[0], body[1]
		entry.Stamp, _ = time.Parse(time.RFC3339, strings.TrimPrefix(head, `time="`))

		ts := strings.Split(tail, `] "`)
		if len(ts) < 2 {
			continue
		}

		tail = strings.TrimRight(ts[1], `"`)
		parts := strings.Split(tail, " ")

		surl := strings.TrimPrefix(parts[1], "http://api.dotagiftx.com")
		purl, _ := url.Parse(surl)

		entry.Endpoint = parts[0] + " " + purl.String()

		entry.Status, _ = strconv.Atoi(parts[6])
		entry.Elapse, _ = time.ParseDuration(parts[9])

		logs = append(logs, entry)
	}

	return logs
}

func parseLog0(raw string) []ApiLog {
	var logs []ApiLog
	for _, i := range strings.Split(raw, "\n") {
		parts := strings.Split(i, " ")
		if len(parts) < 14 {
			continue
		}

		ts := parts[0]
		ts = strings.TrimLeft(ts, "[")
		ts = strings.TrimRight(ts, "]")
		tt, _ := time.Parse(time.RFC3339, ts)
		url := strings.TrimPrefix(parts[5], "http://api.dotagiftx.com")
		stat, _ := strconv.Atoi(parts[10])
		elap, _ := time.ParseDuration(parts[13])

		logs = append(logs, ApiLog{
			Endpoint: url,
			Status:   stat,
			Elapse:   elap,
			Stamp:    tt,
		})
	}

	return logs
}
