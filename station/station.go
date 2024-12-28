package station

import (
	"errors"
	"fmt"
)

var ErrCountZero = errors.New("no elements counted")

// An accumulator object that stores the values (except the Count) in float64
type AccumulatorFloat struct {
	acc   float64 // Accumulated Total
	count uint64  // Number of Samples
	min   float64 // Smallest Sample
	max   float64 // Largest Sample
}

// Creates a new accumulator object
func NewAccumulatorFloat(val float64) *AccumulatorFloat {
	return &AccumulatorFloat{
		acc:   val,
		count: 1,
		min:   val,
		max:   val,
	}
}

// Adds one sample to an existing station accumulator
func (s *AccumulatorFloat) AddSample(val float64) {
	s.acc += val
	s.count++
	s.min = min(val, s.min)
	s.max = max(val, s.max)
}

// Calculates the average of all the samples accumulated
func (s AccumulatorFloat) CalcAvg() (float64, error) {
	var (
		avg float64
		err error
	)
	if s.count == 0 {
		err = ErrCountZero
	} else {
		avg = s.acc / float64(s.count)
	}
	return avg, err
}

// Combines two accumulators into one
func (s1 *AccumulatorFloat) MergeAccumulator(s2 *AccumulatorFloat) {
	if s2 == nil {
		return
	}
	s1.min = min(s1.min, s2.min)
	s1.max = max(s1.max, s2.max)
	s1.acc += s2.acc
	s1.count += s2.count
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
			s.min,
			s.max,
		)
	}
	return str, err
}

// An accumulator object that stores the values (except the Count) in int64
type AccumulatorInt struct {
	acc   int64  // Accumulated Total
	count uint64 // Number of Samples
	min   int64  // Smallest Sample
	max   int64  // Largest Sample
}

// Creates a new accumulator object
func NewAccumulatorInt(val int64) *AccumulatorInt {
	return &AccumulatorInt{
		acc:   val,
		count: 1,
		min:   val,
		max:   val,
	}
}

// Adds one sample to an existing station accumulator
func (s *AccumulatorInt) AddSample(val int64) {
	s.acc += val
	s.count++
	s.min = min(val, s.min)
	s.max = max(val, s.max)
}

// Calculates the average of all the samples accumulated
func (s AccumulatorInt) CalcAvg() (int64, error) {
	var (
		avg int64
		err error
	)
	if s.count == 0 {
		err = ErrCountZero
	} else {
		avg = s.acc / int64(s.count)
	}
	return avg, err
}

// Combines two accumulators into one
func (s1 *AccumulatorInt) MergeAccumulator(s2 *AccumulatorInt) {
	if s2 == nil {
		return
	}
	s1.min = min(s1.min, s2.min)
	s1.max = max(s1.max, s2.max)
	s1.acc += s2.acc
	s1.count += s2.count
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
			float64(avg)/10.0,
			float64(s.min)/10.0,
			float64(s.max)/10.0,
		)
	}
	return str, err
}
