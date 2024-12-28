package main

import (
	"1brc_go/station"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
)

type dataProcessor struct {
	raw     chan *io.SectionReader
	refined chan map[string]*station.AccumulatorFloat
	group   *sync.WaitGroup
}

const MB = 1048576

// var stationSamples sync.Map

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var (
		start  uint64
		offset uint64
		proc   dataProcessor
	)

	fileName, bufSize, threads := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("Error reading %q statistics", fileName)
	}

	fileSize := uint64(stat.Size())
	bufSize = max(bufSize*MB, MB)
	buf := make([]byte, bufSize)

	loops := fileSize / bufSize
	if fileSize%bufSize > 0 {
		loops++
	}
	proc = dataProcessor{
		raw:     make(chan *io.SectionReader, loops),
		refined: make(chan map[string]*station.AccumulatorFloat, loops),
		group:   &sync.WaitGroup{},
	}

	log.Println("File:", fileName)
	log.Println("Number of cores:", threads)
	log.Println("Buffer size:", bufSize)

	t0 := time.Now()
	for range threads {
		go parseSample(bufSize, proc)
		proc.group.Add(1)
	}
	for range loops {
		_, err := file.ReadAt(buf, int64(start))
		if err != nil && err != io.EOF {
			log.Fatal("Error reading buffer:", err)
		}
		offset = uint64(bytes.LastIndexByte(buf, '\n'))
		proc.raw <- io.NewSectionReader(file, int64(start), int64(offset))
		start += offset + 1
		perc := start * 100 / fileSize
		fmt.Fprintf(os.Stderr, "\t\r%d bytes read, %d%%", start, perc)
	}

	close(proc.raw)
	fmt.Fprint(os.Stderr, "\n")
	proc.group.Wait()
	close(proc.refined)
	res := make(map[string]*station.AccumulatorFloat, 1000)
	names := make([]string, 0, 1000)
	for s := range proc.refined {
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

func readFlags() (string, uint64, uint64) {
	var filePath string
	mid := flag.Bool("mid", false, "Program will use the 100M sample file")
	bufSize := flag.Uint64("buffer-size", 0, "Maximum buffer size per thread in MB")
	threads := flag.Uint64("threads", uint64(runtime.NumCPU()), "Maximum number of threads used by the app")
	flag.Parse()

	if *mid {
		filePath = "samples_100M.txt"
	} else {
		filePath = "samples_1B.txt"
	}
	return filePath, *bufSize, max(*threads, 1)
}

func parseSample(bufSize uint64, dp dataProcessor) {
	buf := make([]byte, bufSize)

	for sec := range dp.raw {
		sta := make(map[string]*station.AccumulatorFloat, 1000)
		scn := bufio.NewScanner(sec)
		scn.Buffer(buf, bufio.MaxScanTokenSize)
		scn.Split(bufio.ScanLines)

		for scn.Scan() {
			name, val, err := station.ParseLineFloat(scn.Text())
			if err != nil {
				continue
			}

			if sta[name] == nil {
				sta[name] = station.NewAccumulatorFloat(val)
			} else {
				sta[name].AddSample(val)
			}
		}

		dp.refined <- sta
	}
	dp.group.Done()
}
