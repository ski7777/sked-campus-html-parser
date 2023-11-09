package timetable

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/ski7777/sked-campus-html-parser/pkg/htmltable"
	"strings"
)

type TimeTable struct {
	html   *htmltable.HTMLTable
	Events map[string]Event
}

func ParseTimeTable(dom *goquery.Selection) (tt *TimeTable, err error) {
	tt = &TimeTable{
		Events: make(map[string]Event),
	}
	tt.html, err = htmltable.ParseHTMLTable(dom)
	if err != nil {
		return
	}
	var dateCols []*Date
	for _, df := range tt.html.Fields[0] {
		var date *Date
		if df.HasClass("t") {
			if date, err = DateFromString(df.Text()); err != nil {
				return
			}
		}
		dateCols = append(dateCols, date)
	}
	eventFields := tt.html.FindMatchingFields(
		func(s *goquery.Selection) bool {
			return s.HasClass("v") && strings.HasPrefix(s.AttrOr("id", ""), "z")
		},
	)
	var event Event
	var eventid string
	for _, e := range eventFields {
		if date := dateCols[e.Column]; date == nil {
			err = errors.New("failed to retrieve date")
			return
		} else {
			event, eventid, err = EventFromElement(*date, e)
			if err != nil {
				return
			}
			if _, ok := tt.Events[eventid]; ok {
				err = errors.New("duplicate event id")
				return
			}
			tt.Events[eventid] = event
		}
	}
	return
}
