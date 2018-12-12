package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Rule struct {
	From []bool
	To   bool
}

type Input struct {
	State []bool
	Rules []Rule
}

func ParseState(s string) ([]bool, error) {
	var state []bool
	for _, c := range s {
		if c != '#' && c != '.' {
			return nil, fmt.Errorf("invalid state char: %q", c)
		}
		state = append(state, c == '#')
	}
	return state, nil
}

func ReadInput(file string) (*Input, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	// read the initial state
	if !sc.Scan() {
		return nil, fmt.Errorf("failed to read initial state: %v", sc.Err())
	}
	var initial string
	if _, err := fmt.Sscanf(sc.Text(), "initial state: %s", &initial); err != nil {
		return nil, fmt.Errorf("failed to scan initial state: %v", err)
	}
	state, err := ParseState(initial)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial state: %v", err)
	}

	// scan blank line
	if !sc.Scan() || sc.Text() != "" {
		return nil, fmt.Errorf("missing blank line")
	}

	var rules []Rule
	for sc.Scan() {
		var from, to string
		if _, err := fmt.Sscanf(sc.Text(), "%s => %s", &from, &to); err != nil {
			return nil, fmt.Errorf("failed to scan rule: %v", err)
		}
		var rule Rule
		rule.From, err = ParseState(from)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule from: %v", err)
		}
		rule.To = to == "#"
		rules = append(rules, rule)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &Input{state, rules}, nil
}

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Input", input)
}
