package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Welcome to the One Billion Row Challenge in GO
	fileName := readFlags()
	file, fileErr := os.Open(fileName)
	if fileErr != nil {
		log.Fatalf("Error opening %q", fileName)
	}
	defer file.Close()

	log.Println("Started reading file")
	t0 := time.Now()
	samples, names := processBuffer(file)
	log.Println("Finished reading file in", time.Since(t0))
	for _, name := range names {
		details, err := samples[name].PrintDetails()
		if err != nil {
			log.Println(samples[name])
			log.Fatalf("Error on station %q: %v", name, err)
		}
		fmt.Printf("%s: %v\n", name, details)
	}
}

func readFlags() string {
	mid := flag.Bool("mid", false, "Program will use the 10% of the big file")
	small := flag.Bool("small", false, "Program will use the small sample file")
	flag.Parse()

	if *small {
		return "samples_100K.txt"
	} else if *mid {
		return "samples_100M.txt"
	} else {
		return "samples_1B.txt"
	}
}
