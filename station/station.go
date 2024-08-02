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

func NewStation(val float32) *Station {
	return &Station{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

func (s *Station) AddSample(val float32) {
	s.Acc += val
	s.Count++

	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

func (s *Station) CalcAvg() (float32, error) {
	if s.Count == 0 {
		return 0.0, errors.New("counted elements is zero")
	}
	return s.Acc / float32(s.Count), nil
}

func (s1 *Station) MergeStation(s2 *Station) {
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

func (s *Station) PrintDetails() string {
	avg, err := s.CalcAvg()
	if err != nil {
		return fmt.Sprint(err)
	}
	return fmt.Sprintf("Avg: %.1f, Min: %.1f, Max: %.1f", avg, s.Min, s.Max)
}
