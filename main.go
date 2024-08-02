package main

import (
	"1brc_go/station"
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type StationSample struct {
	name string
	val  float32
}

var stationSamples map[string]*station.Station
var stationNames []string

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var t0 time.Time

	file := openFile("sample_big.txt")
	stationSamples = make(map[string]*station.Station)
	stationNames = make([]string, 0, 1000)
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)

	t0 = time.Now()
	for buffer.Scan() {
		sample := readFileLineV1(buffer)

		if stationSamples[sample.name] != nil {
			stationSamples[sample.name].AddSample(sample.val)
		} else {
			stationSamples[sample.name] = station.NewStation(sample.val)
		}

		ind, exist := slices.BinarySearch(stationNames, sample.name)
		if len(stationNames) == 0 {
			stationNames = append(stationNames, sample.name)
		} else if !exist {
			stationNames = slices.Insert(stationNames, ind, sample.name)
		}
	}
	fmt.Printf("%v\n", time.Since(t0))

	for _, val := range stationNames {
		fmt.Printf("%q: %v\n", val, stationSamples[val].PrintDetails())
	}
}

func readFileLineV1(buffer *bufio.Scanner) *StationSample {
	text := buffer.Text()
	// fmt.Println(string(text))

	textArr := strings.Split(text, ";")
	val, err := strconv.ParseFloat(textArr[1], 32)
	if err != nil {
		log.Fatalln("Error parsing line")
	}

	content := StationSample{
		name: textArr[0],
		val:  float32(val),
	}

	return &content
}

func openFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("Error opening %q", name)
	}
	return file
}
