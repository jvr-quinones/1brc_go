package station

import (
	"errors"
	"fmt"
)

var ErrCountZero = errors.New("no elements counted")

type StationFloat struct {
	// Accumulated Total
	Acc float64
	// Number of Samples
	Count uint64
	// Smallest Sample
	Min float64
	// Largest Sample
	Max float64
}

// Creates a new station object
func NewStationFloat(val float64) *StationFloat {
	return &StationFloat{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

// Adds one sample to an existing station
func (s *StationFloat) AddSample(val float64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

// Calculates the average of all the samples accumulated
func (s StationFloat) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrCountZero
	}
	return s.Acc / float64(s.Count), nil
}

// Combines two stations into one
func (s1 *StationFloat) MergeStation(s2 *StationFloat) {
	if s2 == nil {
		return
	} else if s2.Min < s1.Min {
		s1.Min = s2.Min
	} else if s2.Max > s1.Max {
		s1.Max = s2.Max
	}
	s1.Acc += s2.Max
	s1.Count += s2.Count
}

// Prints the average, minimum and maximum of the station
func (s StationFloat) PrintDetails() (string, error) {
	avg, err := s.CalcAvg()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Avg: %.1f, Min: %.1f, Max: %.1f",
		avg,
		s.Min,
		s.Max,
	), nil
}

type StationInt struct {
	Acc   int64
	Count uint64
	Min   int64
	Max   int64
}

// Creates a new station object
func NewStationInt(val int64) *StationInt {
	return &StationInt{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

// Adds one sample to an existing station
func (s *StationInt) AddSample(val int64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

// Calculates the average of all the samples accumulated
func (s StationInt) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrCountZero
	}
	return float64(s.Acc) / float64(s.Count), nil
}

// Combines two stations into one
func (s1 *StationInt) MergeStation(s2 *StationInt) {
	if s2 == nil {
		return
	} else if s2.Min < s1.Min {
		s1.Min = s2.Min
	} else if s2.Max > s1.Max {
		s1.Max = s2.Max
	}
	s1.Acc += s2.Max
	s1.Count += s2.Count
}

// Prints the average, minimum and maximum of the station
func (s StationInt) PrintDetails() (string, error) {
	avg, err := s.CalcAvg()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Avg: %.1f, Min: %.1f, Max: %.1f",
		avg/10.0,
		float64(s.Min),
		float64(s.Max),
	), nil
}
