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

var stationSamples sync.Map

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var (
		// stationNames   = make(station.StringHeap, 0, 1000)
		line      uint64
		workersCh []chan string
		wg        = &sync.WaitGroup{}
		cores     = 1 // runtime.NumCPU()
	)
	const (
		bufSize = 1e3
	)

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)
	defer file.Close()

	workersCh = make([]chan string, cores)
	for i := range cores {
		ch := make(chan string, bufSize)
		wg.Add(1)
		go parseSample(ch, wg)
		workersCh[i] = ch
	}

	t0 := time.Now()
	tick1s := time.NewTicker(time.Second)
	for line = 0; buffer.Scan(); line++ {
		workersCh[line%uint64(cores)] <- buffer.Text()
		select {
		case <-tick1s.C:
			fmt.Fprintf(os.Stderr, "\rline: %d", line)
		default:
		}
	}
	tick1s.Stop()
	for i := range cores {
		close(workersCh[i])
	}
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

func parseSample(strCh <-chan string, wg *sync.WaitGroup) {
	for sample := range strCh {
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
	wg.Done()
}
