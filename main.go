package main

import (
	"1brc_go/station"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
)

// var stationSamples sync.Map

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var (
		bufSize   int64
		cores     = int64(runtime.NumCPU())
		ptr       int64
		ptrOffset int64
	)

	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("Error reading %q statistics", fileName)
	}
	bufSize = stat.Size() / cores
	log.Println("Number of cores:", cores)
	log.Println("Buffer size:", bufSize)

	buf := make([]byte, bufSize)
	samples := make(chan map[string]*station.AccumulatorFloat, cores)
	wg := &sync.WaitGroup{}
	t0 := time.Now()
	for range cores {
		_, err := file.ReadAt(buf, ptr)
		if err != nil {
			log.Fatal("Error reading buffer:", err)
		}
		ptrOffset = int64(bytes.LastIndexByte(buf, '\n'))
		go parseSample(bytes.Clone(buf[:ptrOffset]), samples, wg)
		wg.Add(1)
		ptr += ptrOffset + 1
		log.Printf("%d bytes read, %d%%", ptrOffset, ptr*100/stat.Size())
	}

	wg.Wait()
	close(samples)
	res := make(map[string]*station.AccumulatorFloat, 1000)
	names := make([]string, 0, 1000)
	for s := range samples {
		for k, v := range s {
			if res[k] == nil {
				res[k] = v
				names = append(names, k)
			} else {
				res[k].MergeAccumulator(v)
			}
		}
	}
	slices.Sort(names)
	log.Printf("Time processing data: %v", time.Since(t0))

	for _, val := range names {
		details, err := res[val].PrintDetails()
		if err != nil {
			fmt.Printf("Error getting details for station %q\n", val)
		}
		fmt.Printf("%s: %v\n", val, details)
	}
}

func readFlags() string {
	mid := flag.Bool("large", false, "Program will use the 100M sample file")
	flag.Parse()

	if *mid {
		return "samples_100M.txt"
	} else {
		return "samples_1B.txt"
	}
}

func parseSample(
	sl []byte, ch chan<- map[string]*station.AccumulatorFloat, wg *sync.WaitGroup,
) {
	m := make(map[string]*station.AccumulatorFloat, 1000)

	buf := bufio.NewScanner(bytes.NewReader(sl))
	buf.Split(bufio.ScanLines)
	for buf.Scan() {
		name, val, err := station.ParseLineFloat(buf.Text())
		if err != nil {
			continue
		}

		if m[name] == nil {
			m[name] = station.NewAccumulatorFloat(val)
		} else {
			m[name].AddSample(val)
		}
	}

	ch <- m
	wg.Done()
}
