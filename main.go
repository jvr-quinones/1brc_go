package main

import (
	"1brc_go/station"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

var (
	stationSamples = make(map[string]*station.StationFloat, 1000)
	stationNames   = make([]string, 0, 1000)
)

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var line int

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)
	defer file.Close()

	t0 := time.Now()
	tick1s := time.NewTicker(time.Second)
	for line = 0; buffer.Scan(); line++ {
		name, val, readErr := station.ParseLineFloat(buffer.Text())
		if readErr != nil {
			fmt.Println("Error reading line", line+1)
			continue
		} else if stationSamples[name] != nil {
			stationSamples[name].AddSample(val)
		} else {
			stationSamples[name] = station.NewStationFloat(val)
			stationNames = append(stationNames, name)
			slices.Sort(stationNames)
		}

		select {
		case <-tick1s.C:
			fmt.Fprintf(os.Stderr, "\r%d lines read", line)
		default:
		}
	}
	tick1s.Stop()
	fmt.Fprintf(os.Stderr, "\r%d lines read in %.3fs", line, time.Since(t0).Seconds())

	for _, val := range stationNames {
		details, err := stationSamples[val].PrintDetails()
		if err != nil {
			fmt.Printf("Error getting details for station %q\n", val)
		}
		fmt.Printf("%s: %v\n", val, details)
	}
}

func readFlags() string {
	mid := flag.Bool("mid", false, "Program will use the 10% of the big file")
	large := flag.Bool("large", false, "Program will use the small sample file")
	flag.Parse()

	if *large {
		return "samples_1B.txt"
	} else if *mid {
		return "samples_100M.txt"
	} else {
		return "samples_100K.txt"
	}
}
