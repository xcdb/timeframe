package timeframe

import (
	"strconv"
	"strings"
	"time"
)

//Absolute parses a token like 2017-03-18 and returns the Range it represents
func Absolute(s string, loc *time.Location) (Range, error) {
	if len(s) < 4 || len(s) > 23 { //2017, 2017-03-18T22:50:00.000
		return err() //cannot be a valid structure
	}
	if loc == nil {
		loc = time.Local
	}
	i := strings.IndexByte(s, 0x2d /*-*/)
	switch i {
	case -1:
		switch len(s) {
		case 4: //2017
			return parseYear(s, loc)
		case 6: //201703
			if allowYYYYMM {
				return parseMonth(s[:4], s[4:], loc)
			}
		case 7: //2017W11, 2017077
			if s[4] == 0x57 /*W*/ { //2017W11
				return parseISOWeek(s[:4], s[5:], loc)
			}
			return parseOrdinalDate(s[:4], s[4:], loc)
		case 8: //20170318, 2017W116
			if s[4] == 0x57 /*W*/ { //2017W116
				return parseISOWeekDate(s[:4], s[5:7], s[7:], loc)
			}
			return parseDate("20060102", s, loc)
		case 11:
			return parseTime("20060102T15", s, time.Hour, loc)
		case 13:
			return parseTime("20060102T1504", s, time.Minute, loc)
		case 15:
			return parseTime("20060102T150405", s, time.Second, loc)
		case 19:
			return parseTime("20060102T150405.000", s, time.Millisecond, loc)
		}
	case 4:
		switch len(s) {
		case 7: //2017-03
			return parseMonth(s[:4], s[5:], loc)
		case 8: //2017-W11, 2017-077
			if s[5] == 0x57 /*W*/ { //2017W116
				return parseISOWeek(s[:4], s[6:], loc)
			}
			return parseOrdinalDate(s[:4], s[5:], loc)
		case 10: //2017-03-18, 2017-W11-6
			if s[5] == 0x57 /*W*/ && s[8] == 0x2d /*-*/ { //2017-W11-6
				return parseISOWeekDate(s[:4], s[6:8], s[9:], loc)
			}
			return parseDate("2006-01-02", s, loc)
		case 13:
			return parseTime("2006-01-02T15", s, time.Hour, loc)
		case 16:
			return parseTime("2006-01-02T15:04", s, time.Minute, loc)
		case 19:
			return parseTime("2006-01-02T15:04:05", s, time.Second, loc)
		case 23:
			return parseTime("2006-01-02T15:04:05.000", s, time.Millisecond, loc)
		}
	}
	return err()
}

func parse(s string, min, max int) int {
	i, e := strconv.Atoi(s)
	if e != nil || i < min || i > max {
		return -1
	}
	return i
}

func parseYear(sy string, loc *time.Location) (Range, error) {
	y := parse(sy, minYear, maxYear)
	if y == -1 {
		return err()
	}
	return year(y, 1, loc)
}

func parseMonth(sy, sm string, loc *time.Location) (Range, error) {
	y := parse(sy, minYear, maxYear)
	m := parse(sm, 1, 12)
	if y == -1 || m == -1 {
		return err()
	}
	return month(y, m, 1, loc)
}

func parseISOWeek(sy, sw string, loc *time.Location) (Range, error) {
	y := parse(sy, minYear, maxYear)
	w := parse(sw, 1, 53)
	if y == -1 || w == -1 {
		return err()
	}
	return week(y, w, 0, 7, loc)
}

func parseISOWeekDate(sy, sw, swd string, loc *time.Location) (Range, error) {
	y := parse(sy, minYear, maxYear)
	w := parse(sw, 1, 53)
	dw := parse(swd, 1, 7)
	if y == -1 || w == -1 || dw == -1 {
		return err()
	}
	return week(y, w, dw-1, 1, loc)
}

func parseOrdinalDate(sy, sdy string, loc *time.Location) (Range, error) {
	y := parse(sy, minYear, maxYear)
	dy := parse(sdy, 1, 366)
	if y == -1 || dy == -1 {
		return err()
	}
	return day(y, 1, dy, 1, loc)
}

func parseDate(layout, value string, loc *time.Location) (Range, error) {
	t, e := time.ParseInLocation(layout, value, loc)
	if e != nil {
		return err()
	}
	return Range{
		LowerInc: t,
		UpperExc: t.AddDate(0, 0, 1),
	}, nil
}

func parseTime(layout, value string, d time.Duration, loc *time.Location) (Range, error) {
	t, e := time.ParseInLocation(layout, value, loc)
	if e != nil {
		return err()
	}
	return Range{
		LowerInc: t,
		UpperExc: t.Add(d),
	}, nil
}
