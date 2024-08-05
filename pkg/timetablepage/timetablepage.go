package timetablepage

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ski7777/sked-campus-html-parser/pkg/timetable"
	"io"
	"net/http"
)

type TimeTablePage struct {
	TimeTables []*timetable.TimeTable
	AllEvents  map[string]timetable.Event
}

func (ttp *TimeTablePage) addTimeTable(tt *timetable.TimeTable) {
	ttp.TimeTables = append(ttp.TimeTables, tt)
}

func (ttp *TimeTablePage) parseTimeTable(dom *goquery.Selection) (err error) {
	tt, err := timetable.ParseTimeTable(dom)
	if err != nil {
		return
	}
	ttp.addTimeTable(tt)
	return
}

func (ttp *TimeTablePage) GetAllEvents() map[string]timetable.Event {
	events := make(map[string]timetable.Event)
	for _, tt := range ttp.TimeTables {
		for id, e := range tt.Events {
			if _, ok := events[id]; ok {
				err := errors.New("duplicate event id")
				panic(err)
			}
			events[id] = e
		}
	}
	return events
}

func ParseHTML(reader io.Reader) (ttp *TimeTablePage, err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}
	ttp = new(TimeTablePage)
	doc.Find("table").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		err = ttp.parseTimeTable(s)
		return err == nil
	})
	ttp.AllEvents = make(map[string]timetable.Event)
	for _, tt := range ttp.TimeTables {
		for id, e := range tt.Events {
			if _, ok := ttp.AllEvents[id]; ok {
				err = errors.New("duplicate event id")
				return
			}
			ttp.AllEvents[id] = e
		}
	}
	return
}

func ParseHTMLURL(url string) (ttp *TimeTablePage, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("failed fetching %s status code error: %d %s", url, res.StatusCode, res.Status))
		return
	}

	return ParseHTML(res.Body)
}
