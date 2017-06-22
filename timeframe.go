package timeframe

import (
	"errors"
	"time"
)

const (
	maxOffset    = 999         //arbitrary
	firstWeekday = time.Monday //as per ISO 8601
	allowYYYYMM  = true        //disallowed by the ISO 8601 (to avoid confusion with YYMMDD)
	minYear      = 0001        //years prior to 1583 are not automatically allowed
	maxYear      = 9999
)

//Expand parses a token like 3_days_ago and returns the Range it represents
func Expand(s string, loc *time.Location) (Range, error) {
	if len(s) < 4 { //shortest tokens are "2017" and "today" respectively
		return err()
	}
	if loc == nil {
		loc = time.Local
	}
	if s[3] >= 0x30 /*0*/ && s[3] <= 0x39 /*9*/ {
		return Absolute(s, loc)
	}
	t := time.Now().In(loc)
	return Relative(s, &t)
}

func err() (Range, error) {
	return Range{}, errors.New("Timeframe not recognised")
}

func year(y, l int, loc *time.Location) (Range, error) {
	year := time.Date(y, 1, 1, 0, 0, 0, 0, loc)
	return Range{
		LowerInc: year,
		UpperExc: year.AddDate(l, 0, 0),
	}, nil
}

func month(y, m, l int, loc *time.Location) (Range, error) {
	month := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, loc)
	return Range{
		LowerInc: month,
		UpperExc: month.AddDate(0, l, 0),
	}, nil
}

func day(y, m, d, l int, loc *time.Location) (Range, error) {
	today := time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	return Range{
		LowerInc: today,
		UpperExc: today.AddDate(0, 0, l),
	}, nil
}

func minute(y, m, d, hh, mm, l int, loc *time.Location) (Range, error) {
	t := time.Date(y, time.Month(m), d, hh, mm, 0, 0, loc)
	return Range{
		LowerInc: t,
		UpperExc: t.Add(time.Minute * time.Duration(l)),
	}, nil
}

//week finds start of week number w in year y, skips o days, then returns a range of l days
func week(y, w, o, l int, loc *time.Location) (Range, error) {
	jan4 := time.Date(y, 1, 4, 0, 0, 0, 0, loc)
	d := jan4.Weekday() - firstWeekday
	if d < 0 {
		d += 7
	}
	week := jan4.AddDate(0, 0, -int(d)+(w*7)-7+o)
	return Range{
		LowerInc: week,
		UpperExc: week.AddDate(0, 0, l),
	}, nil
}
