package timetable

import (
	"errors"
	"github.com/ski7777/sked-campus-html-parser/internal/tools"
	"regexp"
	"strconv"
	"time"
)

var timematcher = regexp.MustCompile(`(?P<hours>[[:digit:]]{1,2}):(?P<minutes>[[:digit:]]{2})`)

var TimeParseError = errors.New("failed to parse time")

type Time struct {
	Hours, Minutes int
}

func (t Time) ToTime(date Date, timezone *time.Location) time.Time {
	return time.Date(
		date.Year,
		time.Month(date.Month),
		date.Day,
		t.Hours,
		t.Minutes,
		0,
		0,
		timezone,
	)
}

func TimeFromString(raw string) (time Time, err error) {
	matches := tools.FindNamedMatches(timematcher, raw)
	if matches == nil {
		err = TimeParseError
		return
	}
	if hour, ok := matches["hours"]; !ok {
		err = TimeParseError
		return
	} else {
		time.Hours, err = strconv.Atoi(hour)
		if err != nil {
			return
		}
	}
	if minutes, ok := matches["minutes"]; !ok {
		err = TimeParseError
		return
	} else {
		time.Minutes, err = strconv.Atoi(minutes)
		if err != nil {
			return
		}
	}
	return
}
