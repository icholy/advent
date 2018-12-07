package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Edge struct {
	From, To string
}

func (e Edge) String() string {
	return fmt.Sprintf("%s -> %s", e.From, e.To)
}

func ReadInput(file string) ([]Edge, error) {
	var ee []Edge
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var e Edge
		_, err := fmt.Sscanf(
			sc.Text(),
			"Step %s must be finished before step %s can begin.",
			&e.From,
			&e.To,
		)
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return ee, nil
}

func main() {
	edges, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range edges {
		fmt.Println(e)
	}
}
