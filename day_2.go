package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("day_2.input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var (
		freq int64
		seen = map[int64]bool{}
		sc   = bufio.NewScanner(f)
	)
	for sc.Scan() {
		x, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		freq += x
		if seen[freq] {
			fmt.Printf("Duplicate: %d", freq)
			return
		}
		seen[freq] = true
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Duplicate not found in %d samples\n", len(seen))
}
