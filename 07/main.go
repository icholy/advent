package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
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

type State int

const (
	Todo State = iota
	Working
	Done
)

type Step struct {
	Name  string
	State State
	Deps  []*Step
}

type ByName []*Step

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s Step) String() string { return s.Name }

func (s Step) Ready() bool {
	for _, d := range s.Deps {
		if d.State != Done {
			return false
		}
	}
	return true
}

func (s Step) Duration() time.Duration {
	for _, r := range s.Name {
		return (time.Duration(r) - 'A' + 1 + 60) * time.Second
	}
	return 0
}

type Graph map[string]*Step

func (g Graph) Done() bool {
	for _, s := range g {
		if s.State != Done {
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

func (g Graph) Todo() []*Step {
	var todo []*Step
	for _, s := range g {
		if s.State == Todo && s.Ready() {
			todo = append(todo, s)
		}
	}
	sort.Sort(ByName(todo))
	return todo
}

func (g Graph) Next() *Step {
	var next *Step
	for _, s := range g {
		if s.State == Todo && s.Ready() {
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

type Worker struct {
	Step *Step
	End  time.Time
	Idle bool
}

func (w *Worker) Update(now time.Time) {
	if !w.Idle && now.After(w.End) {
		w.Idle = true
		w.Step.State = Done
	}
}

func (w *Worker) Do(now time.Time, s *Step) {
	s.State = Working
	w.Idle = false
	w.Step = s
	w.End = now.Add(s.Duration())
}

func (w Worker) String() string {
	if w.Idle {
		return "idle"
	}
	return w.Step.Name
}

type Workers []*Worker

func NewWorkers(n int) Workers {
	var ww Workers
	for i := 0; i < n; i++ {
		ww = append(ww, &Worker{Idle: true})
	}
	return ww
}

func (ww Workers) Idle() (*Worker, bool) {
	for _, w := range ww {
		if w.Idle {
			return w, true
		}
	}
	return nil, false
}

func (ww Workers) String() string {
	var b strings.Builder
	for i, w := range ww {
		fmt.Fprintf(&b, "worker_%d=%s ", i+1, w)
	}
	return b.String()
}

func PartOne(constraints []Constraint) string {
	g := make(Graph)
	for _, c := range constraints {
		g.Add(c)
	}
	var seq strings.Builder
	for s := g.Next(); s != nil; s = g.Next() {
		s.State = Done
		seq.WriteString(s.Name)
	}
	return seq.String()
}

func PartTwo(constraints []Constraint) time.Duration {
	g := make(Graph)
	for _, c := range constraints {
		g.Add(c)
	}
	started := time.Unix(0, 0)
	workers := NewWorkers(5)
	for now := started; true; now = now.Add(time.Second) {
		for _, w := range workers {
			w.Update(now)
		}
		for _, s := range g.Todo() {
			if w, ok := workers.Idle(); ok {
				w.Do(now, s)
			}
		}
		if g.Done() {
			return now.Sub(started)
		}
	}
	panic("unreachable")
}

func main() {
	constraints, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer (Part 1): %s\n", PartOne(constraints))
	fmt.Printf("Answer (Part 2): %s\n", PartTwo(constraints))
}
