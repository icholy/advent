package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	tunnel := NewTunnel(100000)
	tunnel.Init(input.State)

	for i := 0; i < 20; i++ {
		tunnel = tunnel.Apply(input.Rules...)
	}

	fmt.Println(tunnel.Min, tunnel.Max)
	fmt.Println(tunnel.PlantNumSum())
}

type Pot bool

func Format(pp []Pot) string {
	var pattern strings.Builder
	for _, p := range pp {
		pattern.WriteString(p.String())
	}
	return pattern.String()
}

type Tunnel struct {
	Size, Min, Max int
	Pots           []Pot
}

func (t Tunnel) Range(f func(int, Pot)) {
	for i := -t.Size; i < t.Size; i++ {
		f(i, t.At(i))
	}
}

func (t Tunnel) PlantNumSum() int {
	var sum int
	t.Range(func(i int, p Pot) {
		if p {
			sum += i
		}
	})
	return sum
}

func (t Tunnel) At(i int) Pot {
	return t.Pots[t.Size+i]
}

func (t *Tunnel) SetAt(i int, p Pot) {
	if i < t.Min {
		t.Min = i
	}
	if i > t.Max {
		t.Max = i
	}
	t.Pots[t.Size+i] = p
}

func (t Tunnel) String() string {
	return Format(t.Pots[t.Size+t.Min : t.Size+t.Max+1])
}

func (t *Tunnel) Init(pots []Pot) {
	for i, p := range pots {
		if p {
			t.SetAt(i, p)
		}
	}
}

func (t *Tunnel) Apply(rules ...Rule) Tunnel {
	next := NewTunnel(t.Size)
	// next.Init(t.Pots)
	for i := 0; i < len(t.Pots)-5; i++ {
		center := (i + 2) - t.Size
		for _, r := range rules {
			if r.Matches(t.Pots, i) {
				next.SetAt(center, r.To)
			}
		}
	}
	return next
}

func NewTunnel(size int) Tunnel {
	return Tunnel{
		Size: size,
		Pots: make([]Pot, 2*size),
	}
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

func (r Rule) Matches(pots []Pot, offset int) bool {
	for i, p := range r.Pattern {
		if pots[offset+i] != p {
			return false
		}
	}
	return true
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
