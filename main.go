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

type section struct {
	start  uint64
	offset uint64
}

type processor struct {
	in  chan section
	out chan map[string]*station.AccumulatorFloat
	wg  *sync.WaitGroup
}

// var stationSamples sync.Map

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var (
		start     uint64
		end       uint64
		processor processor
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

	if bufSize == 0 {
		bufSize = fileSize / threads
	} else {
		bufSize *= 1048576
	}
	loops := fileSize / bufSize
	if fileSize%bufSize > 0 {
		loops++
	}
	buf := make([]byte, bufSize)
	processor.in = make(chan section, loops)
	processor.out = make(chan map[string]*station.AccumulatorFloat, loops)
	processor.wg = &sync.WaitGroup{}
	log.Println("File:", fileName)
	log.Println("Number of cores:", threads)
	log.Println("Buffer size:", bufSize)

	t0 := time.Now()
	for range threads {
		go parseSample(file, processor)
		processor.wg.Add(1)
	}
	for range loops {
		_, err := file.ReadAt(buf, int64(start))
		if err != nil && err != io.EOF {
			log.Fatal("Error reading buffer:", err)
		}
		end = uint64(bytes.LastIndexByte(buf, '\n'))
		processor.in <- section{
			start:  start,
			offset: end,
		}
		start += end + 1
		perc := start * 100 / uint64(stat.Size())
		fmt.Fprintf(os.Stderr, "\t\r%d bytes read, %d%%", end, perc)
	}

	close(processor.in)
	fmt.Fprint(os.Stderr, "\n")
	processor.wg.Wait()
	close(processor.out)
	res := make(map[string]*station.AccumulatorFloat, 1000)
	names := make([]string, 0, 1000)
	for s := range processor.out {
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
	mid := flag.Bool("large", false, "Program will use the 100M sample file")
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

func parseSample(
	file io.ReaderAt,
	data processor,
) {
	m := make(map[string]*station.AccumulatorFloat, 1000)

	for pointer := range data.in {
		slice := make([]byte, pointer.offset)
		file.ReadAt(slice, int64(pointer.start))
		buf := bufio.NewScanner(bytes.NewReader(slice))
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

		data.out <- m
	}
	data.wg.Done()
}
