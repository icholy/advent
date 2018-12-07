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

type Step struct {
	Name string
	Done bool
	Deps []*Step
}

func (s Step) String() string { return s.Name }

func (s Step) Ready() bool {
	for _, d := range s.Deps {
		if !d.Done {
			return false
		}
	}
	return true
}

type Graph map[string]*Step

func (g Graph) Step(name string) *Step {
	if _, ok := g[name]; !ok {
		g[name] = &Step{Name: name}
	}
	return g[name]
}

func (g Graph) Next() *Step {
	var next *Step
	for _, s := range g {
		if !s.Done && s.Ready() {
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

func main() {
	constraints, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	g := make(Graph)
	for _, c := range constraints {
		g.Add(c)
	}

	for s := g.Next(); s != nil; s = g.Next() {
		fmt.Print(s)
		s.Done = true
	}
}
