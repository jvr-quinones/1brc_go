package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// Welcome to the One Billion Row Challenge in GO

	file := openFile("sample.txt")
	counter := countLinesV1(file)
	fmt.Println("Counter V1:", counter)

	counter = countLinesV2(file)
	fmt.Println("Counter V2:", counter)

}

func countLinesV1(file *os.File) (count int) {
	file.Seek(0, 0)
	buffer := bufio.NewReader(file)
	for line, err := buffer.ReadString('\n'); (len(line) > 0) || err == nil; line, err = buffer.ReadString('\n') {
		count++
	}
	return count
}

func countLinesV2(file *os.File) (count int) {
	file.Seek(0, 0)
	buffer := bufio.NewScanner(file)
	for buffer.Scan() {
		count++
	}
	return count
}

func openFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("Error opening %q", name)
	}
	return file
}
