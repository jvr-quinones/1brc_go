package station

import (
	"errors"
	"fmt"
)

type Station struct {
	Acc   float32
	Count int
	Min   float32
	Max   float32
}

func MergeStations(s1 *Station, s2 *Station) (merged *Station) {
	if s2 == nil {
		return s1
	} else if s1 == nil {
		return s2
	}

	merged = &Station{
		Acc:   s1.Acc + s2.Acc,
		Count: s1.Count + s2.Count,
	}

	if s1.Min <= s2.Min {
		merged.Min = s1.Min
	} else {
		merged.Min = s2.Min
	}

	if s1.Max <= s2.Max {
		merged.Max = s1.Max
	} else {
		merged.Max = s2.Max
	}

	return merged
}

func NewStation(val float32) *Station {
	return &Station{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

func (station *Station) CalcAvg() (float32, error) {
	if station.Count == 0 {
		return 0.0, errors.New("counted elements is zero")
	}
	return station.Acc / float32(station.Count), nil
}

func (station *Station) PrintDetails() string {
	avg, err := station.CalcAvg()
	if err != nil {
		return "No values in station"
	}
	return fmt.Sprintf("Avg: %.1f, Min: %.1f, Max: %.1f", avg, station.Min, station.Max)
}
