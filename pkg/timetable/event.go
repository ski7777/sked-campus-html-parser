package timetable

import (
	"errors"
	"github.com/ski7777/sked-campus-html-parser/internal/tools"
	"github.com/ski7777/sked-campus-html-parser/pkg/htmltable"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
)

var fulltimematcher = regexp.MustCompile(`(?P<begin>[[:digit:]]{1,2}:[[:digit:]]{2}) - (?P<end>[[:digit:]]{1,2}:[[:digit:]]{2}) Uhr`)

type Event struct {
	Date
	Begin, End Time
	Text       []string
}

func EventFromElement(date Date, e htmltable.DimensionalField) (event Event, id string, err error) {
	event.Date = date
	var ok bool
	id, ok = e.Element.Attr("id")
	if !ok {
		err = errors.New("failed to retrieve id")
		return
	}
	children := tools.GetHTMLChildren(e.Element.Nodes[0], func(node *html.Node) bool {
		return node.Type == html.TextNode
	})
	if len(children) == 0 {
		err = errors.New("len(children) == 0: " + strconv.Itoa(len(children)))
	}
	if timematches := tools.FindNamedMatches(fulltimematcher, children[0].Data); timematches == nil {
		err = TimeParseError
		return
	} else {
		if begin, ok := timematches["begin"]; !ok {
			err = TimeParseError
			return
		} else {
			event.Begin, err = TimeFromString(begin)
			if err != nil {
				return
			}
		}
		if end, ok := timematches["end"]; !ok {
			err = TimeParseError
			return
		} else {
			event.End, err = TimeFromString(end)
			if err != nil {
				return
			}
		}
	}
	for _, t := range children[1:] {
		event.Text = append(event.Text, t.Data)
	}
	return
}
