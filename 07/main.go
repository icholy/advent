package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Constraint struct {
	Before, After string
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s before %s", c.Before, c.After)
}

type Step struct {
	Name string
	Done bool
	Deps []*Step
}

type ByName []*Step

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s Step) String() string { return s.Name }

func (s Step) IsReady() bool {
	for _, d := range s.Deps {
		if !d.Done {
			return false
		}
	}
	return true
}

type Graph map[string]*Step

func (g Graph) Done() bool {
	for _, s := range g {
		if !s.Done {
			return false
		}
	}
	return true
}

func (g Graph) Step(name string) *Step {
	if _, ok := g[name]; !ok {
		g[name] = &Step{Name: name}
	}
	return g[name]
}

func (g Graph) Ready() []*Step {
	var ready []*Step
	for _, s := range g {
		if !s.Done && s.IsReady() {
			ready = append(ready, s)
		}
	}
	sort.Sort(ByName(ready))
	return ready
}

func (g Graph) Next() *Step {
	var next *Step
	for _, s := range g {
		if !s.Done && s.IsReady() {
			if next == nil || s.Name < next.Name {
				next = s
			}
		}
	}
	return next
}

func (g Graph) Add(c Constraint) {
	n := g.Step(c.After)
	n.Deps = append(n.Deps, g.Step(c.Before))
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

func PartOne(constraints []Constraint) string {
	g := make(Graph)
	for _, c := range constraints {
		g.Add(c)
	}
	var seq strings.Builder
	for !g.Done() {
		s := g.Next()
		s.Done = true
		seq.WriteString(s.Name)
	}
	return seq.String()
}

func PartTwo(constraints []Constraint) string {
	g := make(Graph)
	for _, c := range constraints {
		g.Add(c)
	}
	var seq strings.Builder
	for !g.Done() {
		for _, s := range g.Ready() {
			s.Done = true
			seq.WriteString(s.Name)
		}
		seq.WriteString("-")
	}
	return seq.String()
}

func main() {
	constraints, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer (Part 1): %s\n", PartOne(constraints))
	fmt.Printf("Answer (Part 2): %s\n", PartTwo(constraints))
}
