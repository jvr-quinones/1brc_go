package main

import (
	"1brc_go/station"
	"bufio"
	"fmt"
	"io"
	"slices"
)

func processBuffer(buf io.Reader) (
	map[string]*station.AccumulatorFloat,
	[]string,
) {
	var (
		line    uint64
		names   = make([]string, 0, 1000)
		samples = make(map[string]*station.AccumulatorFloat, 1000)
	)

	buffer := bufio.NewScanner(buf)
	buffer.Split(bufio.ScanLines)
	for line = 1; buffer.Scan(); line++ {
		name, val, readErr := station.ParseLineFloat(buffer.Text())
		if readErr != nil {
			fmt.Println("Error reading line", line)
			continue
		} else if samples[name] != nil {
			samples[name].AddSample(val)
		} else {
			samples[name] = station.NewAccumulatorFloat(val)
			names = append(names, name)
		}
	}

	slices.Sort(names)
	return samples, names
}
