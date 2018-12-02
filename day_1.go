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
	var sum int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		x, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		sum += x
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sum: %d\n", sum)
}
