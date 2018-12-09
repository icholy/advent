package main

import (
	"fmt"
	"log"
	"os"
)

type Input struct {
	NumPlayers       int
	LastMarblePoints int
}

func ReadInput(file string) (Input, error) {
	var input Input
	f, err := os.Open(file)
	if err != nil {
		return input, err
	}
	defer f.Close()
	if _, err := fmt.Fscanf(f,
		"%d players; last marble is worth %d points",
		&input.NumPlayers,
		&input.LastMarblePoints,
	); err != nil {
		return input, err
	}
	return input, nil
}

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", input)
}
