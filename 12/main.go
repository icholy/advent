package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Pot bool

func (p Pot) String() string {
	if p {
		return "#"
	}
	return "."
}

type Rule struct {
	From []Pot
	To   Pot
}

func (r Rule) String() string {
	var from strings.Builder
	for _, p := range r.From {
		from.WriteString(p.String())
	}
	return fmt.Sprintf("%s => %s", &from, r.To)
}

func (r Rule) Matches(state []Pot, center int) bool {
	return state[center-2] == r.From[0] &&
		state[center-1] == r.From[1] &&
		state[center] == r.From[2] &&
		state[center+1] == r.From[3] &&
		state[center+2] == r.From[4]
}

type Input struct {
	State []Pot
	Rules []Rule
}

func ParseState(s string) ([]Pot, error) {
	var state []Pot
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

	// read the rules
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
	for _, r := range input.Rules {
		fmt.Println(r)
	}
}
