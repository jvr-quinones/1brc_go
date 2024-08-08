package main

import (
	"1brc_go/readline"
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

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)

		t0 := time.Now()
	for line := 1; buffer.Scan(); line++ {
		name, val, readErr := readline.ReadAsFloat(buffer.Text())
		if readErr != nil {
			fmt.Println("Error reading line", line)
continue
		}

		if stationSamples[name] != nil {
			stationSamples[name].AddSample(val)
		} else {
			stationSamples[name] = station.NewStationFloat(val)
		}

		ind, exist := slices.BinarySearch(stationNames, name)
		if len(stationNames) == 0 {
			stationNames = append(stationNames, name)
		} else if !exist {
			stationNames = slices.Insert(stationNames, ind, name)
		}
	}
	fmt.Printf("Read %v in %v??\r", fileName, time.Since(t0))

	for _, val := range stationNames {
		details, err := stationSamples[val].PrintDetails()
		if err != nil {
			fmt.Printf("Error getting details for station %q\n", val)
		}
		fmt.Printf("%q: %v\n", val, details)
	}
}

func readFlags() string {
	const largeFile = "samples_1b.txt"
	const midFile = "samples_100M.txt"
	const smallFile = "samples_100K.txt"

	mid := flag.Bool("mid", false, "Program will use the 10% of the big file")
	large := flag.Bool("large", false, "Program will use the small sample file")
	flag.Parse()

	if *large {
		return largeFile
	} else if *mid {
		return midFile
	} else {
		return smallFile
	}
}
