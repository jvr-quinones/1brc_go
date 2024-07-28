package main

import (
	"1brc_go/station"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type StationSample struct {
	name string
	val  float32
}

var stationSamples map[string]*station.Station

func main() {
	// Welcome to the One Billion Row Challenge in GO
	var t0 time.Time

	file := openFile("sample_big.txt")
	stationSamples = make(map[string]*station.Station)
	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)

	t0 = time.Now()
	for buffer.Scan() {
		line := readFileLineV1(buffer)
		sta := stationSamples[line.name]

		if sta == nil {
			stationSamples[line.name] = station.NewStation(line.val)
			continue
		} else if line.val < sta.Min {
			sta.Min = line.val
		} else if line.val > sta.Max {
			sta.Max = line.val
		}

		sta.Acc += line.val
		sta.Count++
	}
	fmt.Printf("%v\n", time.Since(t0))

	for k, v := range stationSamples {
		fmt.Printf("%q: %v\n", k, v.PrintDetails())
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
