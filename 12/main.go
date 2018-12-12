package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Pot bool

func Format(pp []Pot) string {
	var pattern strings.Builder
	for _, p := range pp {
		pattern.WriteString(p.String())
	}
	return pattern.String()
}

func (p Pot) String() string {
	if p {
		return "#"
	}
	return "."
}

type Rule struct {
	Pattern []Pot
	To      Pot
}

func (r Rule) String() string {
	return fmt.Sprintf("%s => %s", Format(r.Pattern), r.To)
}

func (r Rule) Matches(state []Pot, offset int) bool {
	for i, p := range r.Pattern {
		if state[offset+i] != p {
			return false
		}
	}
	return true
}

func (r Rule) Update(dst, src []Pot) {
	for i := 0; i < len(src)-5; i++ {
		if r.Matches(src, i) {
			dst[i+2] = true
		}
	}
}

type Input struct {
	State []Pot
	Rules []Rule
}

func ParsePots(s string) ([]Pot, error) {
	var pots []Pot
	for _, c := range s {
		if c != '#' && c != '.' {
			return nil, fmt.Errorf("invalid pot char: %q", c)
		}
		pots = append(pots, c == '#')
	}
	return pots, nil
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
	state, err := ParsePots(initial)
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
		var pattern, to string
		if _, err := fmt.Sscanf(sc.Text(), "%s => %s", &pattern, &to); err != nil {
			return nil, fmt.Errorf("failed to scan rule: %v", err)
		}
		var rule Rule
		rule.Pattern, err = ParsePots(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule pattern: %v", err)
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

	prev := input.State
	next := make([]Pot, len(prev))

	for _, r := range input.Rules {
		fmt.Println(r)
		r.Update(next, prev)
	}

	fmt.Println(Format(input.State))
	fmt.Println(Format(next))
}
