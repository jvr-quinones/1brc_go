package station

import (
	"errors"
	"fmt"
)

var ErrNoSep = errors.New("no elements counted")

type StationFloat struct {
	Acc   float64
	Count uint64
	Min   float64
	Max   float64
}

func NewStationFloat(val float64) *StationFloat {
	return &StationFloat{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

func (s *StationFloat) AddSample(val float64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

func (s *StationFloat) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrNoSep
	}
	return s.Acc / float64(s.Count), nil
}

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

func (s *StationFloat) PrintDetails() (string, error) {
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

func NewStationInt(val int64) *StationInt {
	return &StationInt{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

func (s *StationInt) AddSample(val int64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

func (s *StationInt) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrNoSep
	}
	return float64(s.Acc) / float64(s.Count), nil
}

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

func (s *StationInt) PrintDetails() (string, error) {
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
