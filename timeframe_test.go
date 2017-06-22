package timeframe

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

func TestRandomString(t *testing.T) {
	if !testing.Verbose() {
		t.Skip()
	}
	f := func(s string) bool {
		r, _ := Expand(s, time.Local)
		return r.LowerInc.Unix() <= r.UpperExc.Unix()
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

type asciiString string

func (asciiString) Generate(rand *rand.Rand, size int) reflect.Value {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		c := rand.Intn(127)
		buffer.WriteByte(byte(c))
	}
	s := asciiString(buffer.String())
	return reflect.ValueOf(s)
}

func TestRandomAscii(t *testing.T) {
	f := func(s asciiString) bool {
		r, _ := Expand(string(s), time.Local)
		return r.LowerInc.Unix() <= r.UpperExc.Unix()
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

//TestBadExpand verifies unrecognised patterns result in a zero range and a non-nil error
func TestBadExpand(t *testing.T) {
	r, err := Expand("", nil)
	if !r.IsZero() || err == nil {
		t.Fail()
	}
}

func TestExpandNilLoc(t *testing.T) {
	r, err := Expand("today", nil)
	if err != nil || time.Local != r.LowerInc.Location() || time.Local != r.UpperExc.Location() {
		t.Fail()
	}
}

//TestExpandLocation verifies the returned Range maintains the provided location
func TestExpandLocation(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	locations := []*time.Location{
		time.Local,
		time.UTC,
		ny,
	}
	for _, loc := range locations {
		t.Run(loc.String(), func(t *testing.T) {
			r, err := Expand("today", loc)
			if err != nil || loc != r.LowerInc.Location() || loc != r.UpperExc.Location() {
				t.Fail()
			}
		})
	}
}

func BenchmarkExpand(b *testing.B) {
	benchmarks := []string{
		"previous_28_days", "prev_28_days", "next_3_days",
		"today", "3_weeks_ahead", "prev_month", "prev_day",
		"last_30_mins", "5_hours_ago", "next_48_hours", "0_hours_ago",
		"2017", "2017W116", "20170318", "2017-03-18",
		"2017-03-18T22", "2017-03-18T22:50",
	}
	for _, bm := range benchmarks {
		b.Run(bm, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r, err := Expand(bm, time.Local)
				if r.IsZero() || err != nil {
					b.Fail()
				}
			}
		})
	}
}
