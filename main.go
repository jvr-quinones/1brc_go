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
	stationSamples map[string]*station.StationFloat
	stationNames   []string
)

func main() {
	// Welcome to the One Billion Row Challenge in GO

	fileName := readFlags()
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)

	stationSamples = make(map[string]*station.StationFloat)
	stationNames = make([]string, 0, 1000)
	t0 := time.Now()

	for line := 0; buffer.Scan(); line++ {
		name, val, err := readline.ReadAsFloat(buffer.Text())
		if err != nil {
			fmt.Println("Error reading line", line)
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
	fmt.Printf("%v\n", time.Since(t0))

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

	if *large {
		return largeFile
	} else if *mid {
		return midFile
	} else {
		return smallFile
	}
}
