package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Constraint struct {
	Before, After string
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s before %s", c.Before, c.After)
}

func ReadInput(file string) ([]Constraint, error) {
	var cc []Constraint
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var c Constraint
		_, err := fmt.Sscanf(
			sc.Text(),
			"Step %s must be finished before step %s can begin.",
			&c.Before,
			&c.After,
		)
		if err != nil {
			return nil, err
		}
		cc = append(cc, c)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cc, nil
}

func main() {
	constraints, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range constraints {
		fmt.Println(c)
	}
}
