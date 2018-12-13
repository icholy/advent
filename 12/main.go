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
	tunnel := NewTunnel()
	tunnel.Init(input.State)

	PrintTunnel(1, tunnel)

	for i := 0; i < 3000; i++ {
		tunnel = tunnel.Apply(input.Rules...)
		PrintTunnel(i+1, tunnel)
	}

	fmt.Println(tunnel.PlantNumSum())
}

func PrintTunnel(gen int, t *Tunnel) {
	fmt.Printf("%2d: %s\n", gen, t.Shape())
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
	Initialized bool
	Min, Max    int
	Pots        map[int]Pot
}

func (t *Tunnel) Range(min, max int, f func(int, Pot)) {
	for i := min; i <= max; i++ {
		f(i, t.At(i))
	}
}

func (t *Tunnel) PlantNumSum() int {
	var sum int
	t.Range(t.Min, t.Max, func(i int, p Pot) {
		if p {
			sum += i
		}
	})
	return sum
}

func (t *Tunnel) At(i int) Pot {
	return t.Pots[i]
}

func (t *Tunnel) SetAt(i int, p Pot) {
	if !t.Initialized {
		t.Min = i
		t.Max = i
		t.Initialized = true
	}
	if i < t.Min {
		t.Min = i
	}
	if i > t.Max {
		t.Max = i
	}
	if p {
		t.Pots[i] = p
	} else {
		delete(t.Pots, i)
	}
}

func (t *Tunnel) RangeString(min, max int) string {
	var s strings.Builder
	t.Range(min, max, func(_ int, p Pot) {
		s.WriteString(p.String())
	})
	return s.String()
}

func (t *Tunnel) Extents() (min, max int) {
	var initialized bool
	t.Range(t.Min, t.Max, func(i int, p Pot) {
		if !p {
			return
		}
		if !initialized {
			max = i
			min = i
			initialized = true
		}
		if i > max {
			max = i
		}
		if i < min {
			min = i
		}
	})
	return min, max
}

func (t *Tunnel) Shape() string {
	min, max := t.Extents()
	return t.RangeString(min, max)
}

func (t *Tunnel) String() string {
	return t.RangeString(t.Min, t.Max)
}

func (t *Tunnel) Init(pots []Pot) {
	for i, p := range pots {
		if p {
			t.SetAt(i, p)
		}
	}
}

func (t *Tunnel) Apply(rules ...Rule) *Tunnel {
	next := NewTunnel()
	t.Range(t.Min-2, t.Max+2, func(center int, p Pot) {
		for _, r := range rules {
			if r.Matches(t, center) {
				next.SetAt(center, r.To)
			}
		}
	})
	return next
}

func NewTunnel() *Tunnel {
	return &Tunnel{
		Pots: map[int]Pot{},
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

func (r Rule) Matches(t *Tunnel, center int) bool {
	return t.At(center-2) == r.Pattern[0] &&
		t.At(center-1) == r.Pattern[1] &&
		t.At(center) == r.Pattern[2] &&
		t.At(center+1) == r.Pattern[3] &&
		t.At(center+2) == r.Pattern[4]
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
		if len(rule.Pattern) != 5 {
			return nil, fmt.Errorf("expected rule pattern to be 5 pots long, got %d", len(rule.Pattern))
		}
		rule.To = to == "#"
		rules = append(rules, rule)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &Input{state, rules}, nil
}
