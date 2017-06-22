package timeframe

import (
	"strconv"
	"strings"
	"time"
)

//Relative parses a token like 3_days_ago and returns the Range it represents
func Relative(s string, t *time.Time) (Range, error) {
	if len(s) < 5 || len(s) > 20 { //today, previous_999_minutes
		return err() //cannot be a valid structure
	}
	if t == nil {
		now := time.Now()
		t = &now
	}
	i := strings.IndexByte(s, 0x5f /*_*/)
	if i == -1 {
		switch s {
		case "yesterday":
			return newRange("day", -1, -1, t)
		case "today":
			return newRange("day", 0, 0, t)
		case "tomorrow":
			return newRange("day", +1, +1, t)
		default:
			return err()
		}
	}
	j := strings.LastIndexByte(s, 0x5f /*_*/)
	if i == j { //prev_day
		vs := s[:i]
		dp := s[j+1:]
		if len(vs) == 0 || len(dp) == 0 || dp[len(dp)-1] == 0x73 /*s*/ || vs == "last" {
			return err() //disallow s suffix, last_xxx
		}
		return slice(dp, vs, 1, t)
	}
	{ //1_day_ago, prev_1_day
		ns := s[:i]
		dp := s[i+1 : j]
		vs := s[j+1:]
		if len(ns) == 0 || len(dp) == 0 || len(vs) == 0 {
			return err()
		}
		if vs[0] != 0x61 /*a*/ { //prev_1_day
			vs, ns, dp = ns, dp, vs
		}
		n, e := strconv.Atoi(ns)
		if e != nil || n < 0 || n > maxOffset {
			return err() //require n be 0-maxOffset
		}
		if n == 0 && vs[0] != 0x61 /*a*/ {
			return err() //0_days_ago is allowed, last_0_days isn't
		}
		if dp[len(dp)-1] == 0x73 /*s*/ {
			dp = dp[:len(dp)-1] //remove s suffix
		}
		return slice(dp, vs, n, t)
	}
}

func slice(dp, vs string, n int, t *time.Time) (Range, error) {
	var l, u int
	switch vs {
	case "ago":
		l, u = -n, -n
	case "ahead":
		l, u = +n, +n
	case "prev", "previous":
		l, u = -n, -1
	case "last":
		l, u = 1-n, 0
	case "this":
		l, u = 0, n-1
	case "next":
		l, u = 1, n
	default:
		return err()
	}
	return newRange(dp, l, u, t)
}

func newRange(dp string, l, u int, t *time.Time) (Range, error) {
	loc := t.Location()
	switch dp {
	case "day":
		return day(t.Year(), int(t.Month()), t.Day()+l, u-l+1, loc)
	case "hour":
		return minute(t.Year(), int(t.Month()), t.Day(), t.Hour()+l, 0, (u-l+1)*60, loc)
	case "min", "minute":
		return minute(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute()+l, u-l+1, loc)
	case "month":
		return month(t.Year(), int(t.Month()+time.Month(l)), u-l+1, loc)
	case "year":
		return year(t.Year()+l, u-l+1, loc)
	case "week":
		d := t.Weekday() - firstWeekday
		if d < 0 {
			d += 7
		}
		return day(t.Year(), int(t.Month()), -int(d)+t.Day()+l*7, (u-l+1)*7, loc)
	}
	return err()
}
