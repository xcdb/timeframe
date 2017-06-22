package timeframe

import (
	"testing"
	"time"
)

//TestBadAbsolute verifies unrecognised patterns result in a zero range and a non-nil error
func TestBadAbsolute(t *testing.T) {
	patterns := []string{
		"00001332", "00001300", "00001301",
		"00000032", "00000000", "00000001",
		"00000132", "00000100",
		"0000-13-32", "0000-13-00", "0000-13-01",
		"0000-00-32", "0000-00-00", "0000-00-01",
		"0000-01-32", "0000-01-00",
		"0000-00", "0000-13", "000000", "000013", "0000",
		"0000W00", "0000-W00", "0000W54", "0000-W54",
		"0000W230", "0000-W23-0", "0000W238", "0000-W23-8",
		"0000W23 ", "0000-W-23 ",
		"0000000", "0000-000", "0000367", "0000-367",
		"20151332", "20151300", "20151301",
		"20150032", "20150000", "20150001",
		"20150132", "20150100",
		"2015-13-32", "2015-13-00", "2015-13-01",
		"2015-00-32", "2015-00-00", "2015-00-01",
		"2015-01-32", "2015-01-00",
		"2015-00", "2015-13", "201500", "201513",
		"2015W00", "2015-W00", "2015W54", "2015-W54",
		"2015W230", "2015-W23-0", "2015W238", "2015-W23-8",
		"2015W23 ", "2015-W-23 ",
		"2015000", "2015-000", "2015367", "2015-367",
		"2017--03", "2017--03--18",
		"2017-03x18", "2017-W11x6", "2017x03x18", "2017xW11x6",
		"20170318x22", "20170318x2250", "20170318T2", "20170318T225",
		"2017-03-18x22", "2017-03-18x22:50", "2017-03-18T2", "2017-03-18T22:5",
		"2017-03-18T22x50", "2017-03-18T22:50x42", "2017-03-18T22:50:42x000",
		"20170318T225042x000", "20170318T225042000",
		"2017-03-18T22:50:42.0", "2017-03-18T22:50:42.00", "2017-03-18T22:50:42.0000",

		//TODO: tests for a year that has only 52 weeks in it?
		//TODO: tests for a year that has only 365 days in it?
	}
	for _, p := range patterns {
		t.Run(p, func(t *testing.T) {
			r, err := Absolute(p, nil)
			if !r.IsZero() || err == nil {
				t.Fail()
			}
		})
	}
}

//TestAbsoluteLocation verifies the returned Range maintains the provided location
func TestAbsoluteLocation(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	locations := []*time.Location{
		time.Local,
		time.UTC,
		ny,
	}
	for _, loc := range locations {
		t.Run(loc.String(), func(t *testing.T) {
			r, err := Absolute("2017-03-18", loc)
			if err != nil || loc != r.LowerInc.Location() || loc != r.UpperExc.Location() {
				t.Fail()
			}
		})
	}
}

func TestAbsolute(t *testing.T) {
	const format = "2006-01-02"
	loc := time.Local
	testCases := []struct {
		pat   string
		lower string
		upper string
	}{
		{"2017", "2017-01-01", "2018-01-01"},
		{"2017-03", "2017-03-01", "2017-04-01"},
		{"201703", "2017-03-01", "2017-04-01"}, //extension to ISO 8601
		{"2017-03-18", "2017-03-18", "2017-03-19"},
		{"20170318", "2017-03-18", "2017-03-19"},
		{"2015-W01", "2014-12-29", "2015-01-05"}, //boosts coverage
		{"2017-W11", "2017-03-13", "2017-03-20"},
		{"2017W11", "2017-03-13", "2017-03-20"},
		{"2017-W11-6", "2017-03-18", "2017-03-19"},
		{"2017W116", "2017-03-18", "2017-03-19"},
		{"2017-077", "2017-03-18", "2017-03-19"},
		{"2017077", "2017-03-18", "2017-03-19"},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			r, err := Absolute(tc.pat, loc)
			if err != nil {
				t.Fail()
			}
			if inc, _ := time.ParseInLocation(format, tc.lower, loc); inc != r.LowerInc {
				t.Errorf("L %s %s", inc, r.LowerInc)
			}
			if exc, _ := time.ParseInLocation(format, tc.upper, loc); exc != r.UpperExc {
				t.Errorf("U %s %s", exc, r.UpperExc)
			}
		})
	}
}

func TestAbsoluteTime2(t *testing.T) {
	const format = "2006-01-02 15:04:05.000"
	loc := time.Local
	testCases := []struct {
		pat   string
		lower string
		upper string
	}{
		{"2017-03-18T22:51:42", "2017-03-18 22:51:42.000", "2017-03-18 22:51:43.000"},
		{"2017-03-18T22:51:42.123", "2017-03-18 22:51:42.123", "2017-03-18 22:51:42.124"},
		{"20170318T225142", "2017-03-18 22:51:42.000", "2017-03-18 22:51:43.000"},
		{"20170318T225142.123", "2017-03-18 22:51:42.123", "2017-03-18 22:51:42.124"},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			r, err := Absolute(tc.pat, loc)
			if err != nil {
				t.Fail()
			}
			if inc, _ := time.ParseInLocation(format, tc.lower, loc); inc != r.LowerInc {
				t.Errorf("L %s %s", inc, r.LowerInc)
			}
			if exc, _ := time.ParseInLocation(format, tc.upper, loc); exc != r.UpperExc {
				t.Errorf("U %s %s", exc, r.UpperExc)
			}
		})
	}
}

func TestAbsoluteTime(t *testing.T) {
	const format = "2006-01-02 15:04"
	loc := time.Local
	testCases := []struct {
		pat   string
		lower string
		upper string
	}{
		{"2017-03-18T22", "2017-03-18 22:00", "2017-03-18 23:00"},
		{"2017-03-18T22:50", "2017-03-18 22:50", "2017-03-18 22:51"},
		{"20170318T22", "2017-03-18 22:00", "2017-03-18 23:00"},
		{"20170318T2250", "2017-03-18 22:50", "2017-03-18 22:51"},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			r, err := Absolute(tc.pat, loc)
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

func TestAbsoluteDaylightSavingTime(t *testing.T) {
	const format = "2006-01-02 15:04 MST"
	loc, _ := time.LoadLocation("Europe/London")
	testCases := []struct {
		pat   string
		lower string
		upper string
	}{
		{"2017-03-26", "2017-03-26 00:00 GMT", "2017-03-27 00:00 BST"},
		{"2017-03-26T00", "2017-03-26 00:00 GMT", "2017-03-26 02:00 BST"},
		{"2017-03-26T01", "2017-03-26 02:00 BST", "2017-03-26 03:00 BST"},
		{"2017-03-26T01:30", "2017-03-26 02:30 BST", "2017-03-26 02:31 BST"},
		{"20170326T01", "2017-03-26 02:00 BST", "2017-03-26 03:00 BST"},
		{"20170326T0130", "2017-03-26 02:30 BST", "2017-03-26 02:31 BST"},
		{"20170326T02", "2017-03-26 02:00 BST", "2017-03-26 03:00 BST"},
		{"20170326T0230", "2017-03-26 02:30 BST", "2017-03-26 02:31 BST"},

		{"2017-10-29", "2017-10-29 00:00 BST", "2017-10-30 00:00 GMT"},
		{"2017-10-29T00", "2017-10-29 00:00 BST", "2017-10-29 01:00 BST"},
		//{"2017-10-29T01", "2017-10-29 01:00 BST", "2017-10-29 01:00 GMT"},    /* When the clock rolls back, T01 is ambiguous.  */
		//{"2017-10-29T01:30", "2017-10-29 01:30 BST", "2017-10-29 01:31 BST"}, /* We choose to go with the std lib interpretation...  */
		{"2017-10-29T01", "2017-10-29 01:00 GMT", "2017-10-29 02:00 GMT"},
		{"2017-10-29T01:30", "2017-10-29 01:30 GMT", "2017-10-29 01:31 GMT"},
		{"2017-10-29T02", "2017-10-29 02:00 GMT", "2017-10-29 03:00 GMT"},
		{"2017-10-29T02:30", "2017-10-29 02:30 GMT", "2017-10-29 02:31 GMT"},
	}
	for _, tc := range testCases {
		t.Run(tc.pat, func(t *testing.T) {
			r, err := Absolute(tc.pat, loc)
			if err != nil {
				t.Fail()
			}
			if inc, _ := time.ParseInLocation(format, tc.lower, loc); inc != r.LowerInc {
				t.Errorf("L %s %s", inc, r.LowerInc)
			}
			if exc, _ := time.ParseInLocation(format, tc.upper, loc); exc != r.UpperExc {
				t.Errorf("U %s %s", exc, r.UpperExc)
			}
		})
	}
}
