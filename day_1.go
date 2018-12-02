package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("day_1.input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var freq int64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		x, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		freq += x
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Frequency: %d\n", freq)
}
