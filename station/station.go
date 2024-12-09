package station

import (
	"errors"
	"fmt"
)

var ErrCountZero = errors.New("no elements counted")

// An accumulator object that stores the values (except the Count) in float64
type AccumulatorFloat struct {
	// Accumulated Total
	Acc float64
	// Number of Samples
	Count uint64
	// Smallest Sample
	Min float64
	// Largest Sample
	Max float64
}

// Creates a new accumulator object
func NewAccumulatorFloat(val float64) *AccumulatorFloat {
	return &AccumulatorFloat{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

// Adds one sample to an existing station accumulator
func (s *AccumulatorFloat) AddSample(val float64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

// Calculates the average of all the samples accumulated
func (s AccumulatorFloat) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrCountZero
	}
	return s.Acc / float64(s.Count), nil
}

// Combines two accumulators into one
func (s1 *AccumulatorFloat) MergeAccumulator(s2 *AccumulatorFloat) {
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
func (s AccumulatorFloat) PrintDetails() (string, error) {
	var str string
	avg, err := s.CalcAvg()
	if err == nil {
		str = fmt.Sprintf(
			"%.1f/%.1f/%.1f",
			avg,
			s.Min,
			s.Max,
		)
	}
	return str, err
}

// An accumulator object that stores the values (except the Count) in int64
type AccumulatorInt struct {
	Acc   int64
	Count uint64
	Min   int64
	Max   int64
}

// Creates a new accumulator object
func NewAccumulatorInt(val int64) *AccumulatorInt {
	return &AccumulatorInt{
		Acc:   val,
		Count: 1,
		Min:   val,
		Max:   val,
	}
}

// Adds one sample to an existing station accumulator
func (s *AccumulatorInt) AddSample(val int64) {
	s.Acc += val
	s.Count++
	if val < s.Min {
		s.Min = val
	} else if val > s.Max {
		s.Max = val
	}
}

// Calculates the average of all the samples accumulated
func (s AccumulatorInt) CalcAvg() (float64, error) {
	if s.Count == 0 {
		return 0.0, ErrCountZero
	}
	return float64(s.Acc) / float64(s.Count), nil
}

// Combines two accumulators into one
func (s1 *AccumulatorInt) MergeAccumulator(s2 *AccumulatorInt) {
	if s2 == nil {
		return
	} else if s2.Min < s1.Min {
		s1.Min = s2.Min
	} else if s2.Max > s1.Max {
		s1.Max = s2.Max
	}
	s1.Acc += s2.Acc
	s1.Count += s2.Count
}

// Prints the average, minimum and maximum of the station
func (s AccumulatorInt) PrintDetails() (string, error) {
	var str string
	avg, err := s.CalcAvg()
	if err == nil {
		str = fmt.Sprintf(
			"%.1f/%.1f/%.1f",
			avg/10.0,
			float64(s.Min),
			float64(s.Max),
		)
	}
	return str, err
}
