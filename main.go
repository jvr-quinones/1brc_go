package main

import (
	"1brc_go/station"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var stationSamples sync.Map

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var (
		// stationNames   = make(station.StringHeap, 0, 1000)
		line      uint64
		prevLine  uint64 = 0
		workersCh chan []string
		wg        = &sync.WaitGroup{}
		cores     = runtime.NumCPU() - 1
	)
	const (
		bufSize uint64 = 3e3
	)

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)
	defer file.Close()

	workersCh = make(chan []string, cores)
	for range cores {
		wg.Add(1)
		go parseSample(workersCh, wg)
	}

	t0 := time.Now()
	tick1s := time.NewTicker(time.Second)
	strArray := make([]string, bufSize)
	nextScan := buffer.Scan()

	for line = 0; nextScan; line++ {
		if (line%bufSize == 0 && line > 0) || !nextScan {
			workersCh <- strArray
		}
		strArray[line%bufSize] = buffer.Text()
		nextScan = buffer.Scan()

		select {
		case <-tick1s.C:
			fmt.Fprintf(os.Stderr, "\r%d lines/sec", line-prevLine)
			prevLine = line
		default:
		}
	}

	tick1s.Stop()
	close(workersCh)
	wg.Wait()
	fmt.Fprintf(os.Stderr, "\r%d lines read in %.3fs", line, time.Since(t0).Seconds())

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
	large := flag.Bool("large", false, "Program will use the small sample file")
	flag.Parse()

	if *large {
		return "samples_1B.txt"
	} else {
		return "samples_100M.txt"
	}
}

func parseSample(strCh <-chan []string, wg *sync.WaitGroup) {
	for array := range strCh {
		for _, sample := range array {
			name, val, err := station.ParseLineFloat(sample)
			if err != nil {
				// fmt.Println("Error reading line", sample.number)
				continue
			}

			pStation, _ := stationSamples.Load(name)
			if pStation != nil {
				pStation.(*station.StationFloat).AddSample(val)
			} else {
				stationSamples.Store(name, station.NewStationFloat(val))
			}
		}
	}
	wg.Done()
}
