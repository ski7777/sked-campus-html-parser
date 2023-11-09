package timetable

import (
	"github.com/ski7777/sked-campus-html-parser/internal/tools"
	"regexp"
	"strconv"
)

var datematcher = regexp.MustCompile(`(?P<day>[[:digit:]]{2})\.(?P<month>[[:digit:]]{2})\.(?P<year>[[:digit:]]{4})`)

type Date struct {
	Day, Month, Year int
}

func DateFromString(raw string) (date *Date, err error) {
	rawdate := tools.FindNamedMatches(datematcher, raw)
	date = &Date{}
	date.Day, err = strconv.Atoi(rawdate["day"])
	if err != nil {
		return
	}
	date.Month, err = strconv.Atoi(rawdate["month"])
	if err != nil {
		return
	}
	date.Year, err = strconv.Atoi(rawdate["year"])
	if err != nil {
		return
	}
	return
}
