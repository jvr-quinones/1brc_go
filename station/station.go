package station

import (
	"errors"
	"fmt"
)

var ErrCountZero = errors.New("no elements counted")

// An accumulator object that stores the values (except the Count) in float64
type AccumulatorFloat struct {
	Acc   float64 // Accumulated Total
	Count uint64  // Number of Samples
	Min   float64 // Smallest Sample
	Max   float64 // Largest Sample
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
	s.Min = min(val, s.Min)
	s.Max = max(val, s.Max)
}

// Calculates the average of all the samples accumulated
func (s AccumulatorFloat) CalcAvg() (float64, error) {
	var (
		avg float64
		err error
	)
	if s.Count == 0 {
		err = ErrCountZero
	} else {
		avg = s.Acc / float64(s.Count)
	}
	return avg, err
}

// Combines two accumulators into one
func (s1 *AccumulatorFloat) MergeAccumulator(s2 *AccumulatorFloat) {
	if s2 == nil {
		return
	}
	s1.Min = min(s1.Min, s2.Min)
	s1.Max = max(s1.Max, s2.Max)
	s1.Acc += s2.Acc
	s1.Count += s2.Count
}

// Prints the average, minimum and maximum of the station
func (s AccumulatorFloat) PrintDetails() (string, error) {
	var (
		str string
		err error
	)
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
	Acc   int64  // Accumulated Total
	Count uint64 // Number of Samples
	Min   int64  // Smallest Sample
	Max   int64  // Largest Sample
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
	s.Min = min(val, s.Min)
	s.Max = max(val, s.Max)
}

// Calculates the average of all the samples accumulated
func (s AccumulatorInt) CalcAvg() (int64, error) {
	var (
		avg int64
		err error
	)
	if s.Count == 0 {
		err = ErrCountZero
	} else {
		avg = s.Acc / int64(s.Count)
	}
	return avg, err
}

// Combines two accumulators into one
func (s1 *AccumulatorInt) MergeAccumulator(s2 *AccumulatorInt) {
	if s2 == nil {
		return
	}
	s1.Min = min(s1.Min, s2.Min)
	s1.Max = max(s1.Max, s2.Max)
	s1.Acc += s2.Acc
	s1.Count += s2.Count
}

// Prints the average, minimum and maximum of the station
func (s AccumulatorInt) PrintDetails() (string, error) {
	var (
		str string
		err error
	)
	avg, err := s.CalcAvg()
	if err == nil {
		str = fmt.Sprintf(
			"%.1f/%.1f/%.1f",
			float64(avg/10.0),
			float64(s.Min),
			float64(s.Max),
		)
	}
	return str, err
}
