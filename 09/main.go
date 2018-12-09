package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

type Marble struct {
	Num  int
	Next *Marble // clockwise
}

type Circle struct {
	Size    int
	First   *Marble
	Current *Marble
}

func (c *Circle) Clockwise(n int) *Marble {
	m := c.Current
	for i := 0; i < n; i++ {
		m = m.Next
	}
	return m
}

func (c *Circle) Each(f func(*Marble)) {
	m := c.First
	for i := 0; i < c.Size; i++ {
		if i > 0 {
			m = m.Next
		}
		f(m)
	}
}

func (c *Circle) String() string {
	var b strings.Builder
	c.Each(func(m *Marble) {
		if m == c.Current {
			fmt.Fprintf(&b, "(%d) ", m.Num)
		} else {
			fmt.Fprintf(&b, "%d ", m.Num)
		}
	})
	return b.String()
}

func (c *Circle) Insert(marble int) {
	m := &Marble{Num: marble}
	if c.Current == nil {
		m.Next = m
		c.Current = m
		c.First = m
	} else {
		prev := c.Clockwise(1)
		after := c.Clockwise(2)
		prev.Next = m
		m.Next = after
		c.Current = m
	}
	c.Size++
}

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", input)
}
