package main

import (
	"1brc_go/station"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	stationSamples sync.Map
	stationNames   = make(station.StringHeap, 0, 1000)
)

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var line uint64

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)
	defer file.Close()

	cores := 8
	t0 := time.Now()
	tick1s := time.NewTicker(time.Second)
	sampleCh := make(chan string, cores*1e6)
	fmt.Fprintf(os.Stderr, "Using %d core(s)\n", cores)

	for range cores {
		go worker(sampleCh, &stationSamples)
	}

	for line = 0; buffer.Scan(); line++ {
		sampleCh <- buffer.Text()
		select {
		case <-tick1s.C:
			fmt.Fprintf(os.Stderr, "\rlen: %d, line: %d", len(sampleCh), line)
		default:
		}
	}
	fmt.Fprintf(os.Stderr, "\r%d lines read in %.3fs", line, time.Since(t0).Seconds())
	tick1s.Stop()
	close(sampleCh)

	stationSamples.Range(func(k any, v any) bool {
		details, err := v.(*station.StationFloat).PrintDetails()
		if err != nil {
			fmt.Printf("Could not print details for stations %q", k.(string))
		} else {
			fmt.Println(k.(string), details)
		}
		return true
	})
	// for _, val := range stationNames {
	// 	details, err := stationSamples[val].PrintDetails()
	// 	if err != nil {
	// 		fmt.Printf("Error getting details for station %q\n", val)
	// 	}
	// 	fmt.Printf("%s: %v\n", val, details)
	// }
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

type StationSample struct {
	number uint64
	text   string
}

func worker(sampleCh <-chan string, pSyncMap *sync.Map) {
	for sample := range sampleCh {
		name, val, err := station.ParseLineFloat(sample)
		if err != nil {
			// fmt.Println("Error reading line", sample.number)
			continue
		}

		pStation, _ := pSyncMap.Load(name)
		if pStation != nil {
			pStation.(*station.StationFloat).AddSample(val)
		} else {
			pSyncMap.Store(name, station.NewStationFloat(val))
		}
	}
}
