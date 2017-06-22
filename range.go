package timeframe

import "time"

//Range represents a timespan of LowerInc to UpperExc
type Range struct {
	LowerInc time.Time
	UpperExc time.Time
}

//IsZero reports whether r represents a range that has been initialised explicitly
func (r *Range) IsZero() bool {
	return r.LowerInc.IsZero() && r.UpperExc.IsZero()
}
