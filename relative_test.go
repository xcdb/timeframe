package timeframe

import (
	"testing"
	"time"
)

//TestBadRelative verifies unrecognised patterns result in a zero range and a non-nil error
func TestBadRelative(t *testing.T) {
	patterns := []string{
		"", "_", "__", "___",
		"_day_", "12_days", "12_days_", "12__ago",
		"prev_", "prev__", "prev_13_", "evth_1_day",
		"hello", "something_else_entirely", //not recognised
		"this__day", //multiple underscores
		" today ",   //whitespace
		" today",
		"today ",
		"1234_days_ago", //more than 3 digit number
		"a123_days_ago",
		"123a_hours_ago",
		"last_1234_days",
		"23_da_ago", "12_months_ahe", "this_sa", "nex_2_days", //incomplete
		"-12_days_ago", //negative number
		"last_-12_days",
		"prev_days", //s suffix
		"this_weeks",
		"next_years",
		"prev_0_days", //zero only supported on 'ago' and 'ahead'
		"previous_0_days",
		"last_0_weeks",
		"this_0_months",
		"next_0_years",
		"last_day", //ambiguous with prev_month (would actually be treated same as this_month)
		"last_week",
		"last_month",
		"last_year",
		"last_hour",
		"last_min",
		"last_minute",
	}
	for _, p := range patterns {
		t.Run(p, func(t *testing.T) {
			r, err := Relative(p, nil)
			if !r.IsZero() || err == nil {
				t.Fail()
			}
		})
	}
}

//TestRelativeLocation verifies the returned Range maintains the provided location
func TestRelativeLocation(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	locations := []*time.Location{
		time.Local,
		time.UTC,
		ny,
	}
	for _, loc := range locations {
		t.Run(loc.String(), func(t *testing.T) {
			now := time.Now().In(loc)
			r, err := Relative("today", &now)
			if err != nil || loc != r.LowerInc.Location() || loc != r.UpperExc.Location() {
				t.Fail()
			}
		})
	}
}

func TestRelative(t *testing.T) {
	const format = "2006-01-02"
	loc, _ := time.LoadLocation("Europe/London")
	rel := time.Date(2017, 04, 16, 0, 0, 0, 0, loc)
	testCases := []struct {
		pat   string
		lower string
		upper string
	}{
		{"today", "2017-04-16", "2017-04-17"},
		{"this_week", "2017-04-10", "2017-04-17"}, //boosts coverage
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			r, err := Relative(tc.pat, &rel)
			if err != nil {
				t.Fail()
			}
			if inc, _ := time.ParseInLocation(format, tc.lower, loc); inc != r.LowerInc {
				t.Fail()
			}
			if exc, _ := time.ParseInLocation(format, tc.upper, loc); exc != r.UpperExc {
				t.Fail()
			}
		})
	}
}

func TestRelativeAlias(t *testing.T) {
	testCases := []struct {
		pat      string
		datepart string
		incr     int
	}{
		{"yesterday", "day", -1},
		{"today", "day", 0},
		{"tomorrow", "day", +1},
		{"prev_day", "day", -1},
		{"previous_day", "day", -1},
		{"this_day", "day", 0},
		{"next_day", "day", +1},
		{"prev_week", "week", -1},
		{"previous_week", "week", -1},
		{"this_week", "week", 0},
		{"next_week", "week", +1},
		{"prev_month", "month", -1},
		{"previous_month", "month", -1},
		{"this_month", "month", 0},
		{"next_month", "month", +1},
		{"prev_year", "year", -1},
		{"previous_year", "year", -1},
		{"this_year", "year", 0},
		{"next_year", "year", +1},
		{"prev_hour", "hour", -1},
		{"previous_hour", "hour", -1},
		{"this_hour", "hour", 0},
		{"next_hour", "hour", +1},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			e := expected(tc.datepart, tc.incr, tc.incr)
			r, err := Relative(tc.pat, nil)
			if err != nil || e != r {
				t.Fail()
			}
		})
	}
}

func TestRelativeSlice(t *testing.T) {
	testCases := []struct {
		pat      string
		datepart string
		lower    int
		upper    int
	}{
		{"999_days_ago", "day", -999, -999},
		{"12_months_ago", "month", -12, -12},
		{"4_weeks_ago", "week", -4, -4},
		{"1_hour_ago", "hour", -1, -1},
		{"000_months_ago", "month", 0, 0},

		{"999_minutes_ahead", "minute", 999, 999},
		{"3_days_ahead", "day", 3, 3},
		{"6_month_ahead", "month", 6, 6},
		{"04_years_ahead", "year", 4, 4},
		{"1_week_ahead", "week", 1, 1},
		{"0_hours_ahead", "hour", 0, 0},

		{"prev_999_days", "day", -999, -1},
		{"prev_240_months", "month", -240, -1},
		{"prev_11_weeks", "week", -11, -1},
		{"prev_1_year", "year", -1, -1},

		{"previous_999_days", "day", -999, -1},
		{"previous_240_months", "month", -240, -1},
		{"previous_11_weeks", "week", -11, -1},
		{"previous_1_year", "year", -1, -1},

		{"next_999_days", "day", 1, 999},
		{"next_400_weeks", "week", 1, 400},
		{"next_020_mins", "minute", 1, 20},
		{"next_001_month", "month", 1, 1},

		{"last_999_months", "month", -998, 0},
		{"last_678_years", "year", -677, 0},
		{"last_31_days", "day", -30, 0},
		{"last_01_weeks", "week", 0, 0},

		{"this_999_year", "year", 0, 998},
		{"this_111_weeks", "week", 0, 110},
		{"this_11_months", "month", 0, 10},
		{"this_1_day", "day", 0, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			e := expected(tc.datepart, tc.lower, tc.upper)
			r, err := Relative(tc.pat, nil)
			if err != nil || e != r {
				t.Fail()
			}
		})
	}
}

func expected(datepart string, lower int, upper int) Range {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	switch datepart {
	case "minute":
		minute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)
		return Range{
			LowerInc: minute.Add(time.Minute * time.Duration(lower)),
			UpperExc: minute.Add(time.Minute * time.Duration(upper+1)),
		}
	case "hour":
		hour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local)
		return Range{
			LowerInc: hour.Add(time.Hour * time.Duration(lower)),
			UpperExc: hour.Add(time.Hour * time.Duration(upper+1)),
		}
	case "day":
		return Range{
			LowerInc: today.AddDate(0, 0, lower),
			UpperExc: today.AddDate(0, 0, upper+1),
		}
	case "month":
		month := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)
		return Range{
			LowerInc: month.AddDate(0, lower, 0),
			UpperExc: month.AddDate(0, upper+1, 0),
		}
	case "year":
		year := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		return Range{
			LowerInc: year.AddDate(lower, 0, 0),
			UpperExc: year.AddDate(upper+1, 0, 0),
		}
	case "week":
		diff := today.Weekday() - time.Monday
		if diff < 0 {
			diff += 7
		}
		week := today.AddDate(0, 0, -int(diff))
		return Range{
			LowerInc: week.AddDate(0, 0, lower*7),
			UpperExc: week.AddDate(0, 0, upper*7+7),
		}
	}
	panic("unexpected datepart")
}
