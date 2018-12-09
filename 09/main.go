package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Input struct {
	NumPlayers int
	NumMarbles int
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
		&input.NumMarbles,
	); err != nil {
		return input, err
	}
	return input, nil
}

type Marble struct {
	Num  int
	Next *Marble // clockwise
	Prev *Marble
}

func (m *Marble) InsertAfter(next *Marble) {
	next.Prev = m
	next.Next = m.Next
	m.Next.Prev = next
	m.Next = next
}

func (m *Marble) Delete() {
	m.Next.Prev = m.Prev
	m.Prev.Next = m.Next
}

type Circle struct {
	Size    int
	First   *Marble
	Current *Marble
}

func NewCircle() *Circle {
	m := &Marble{Num: 0}
	m.Next = m
	m.Prev = m
	return &Circle{
		Size:    1,
		First:   m,
		Current: m,
	}
}

func (c *Circle) Clockwise(n int) *Marble {
	m := c.Current
	for i := 0; i < n; i++ {
		m = m.Next
	}
	return m
}

func (c *Circle) CounterClockwise(n int) *Marble {
	m := c.Current
	for i := 0; i < n; i++ {
		m = m.Prev
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

func (c Circle) String() string {
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

func (c *Circle) Score(marble int) int {
	m := c.CounterClockwise(7)
	score := m.Num + marble
	m.Delete()
	c.Current = m.Next
	c.Size--
	return score
}

func (c *Circle) Place(marble int) {
	m := &Marble{Num: marble}
	c.Clockwise(1).InsertAfter(m)
	c.Current = m
	c.Size++
}

func PartOne(input Input) int {
	var (
		circle  = NewCircle()
		players = make([]int, input.NumPlayers)
	)
	for i := 1; i <= input.NumMarbles; i++ {
		player := i % input.NumPlayers
		if i != 0 && i%23 == 0 {
			players[player] += circle.Score(i)
		} else {
			circle.Place(i)
		}
	}
	var max int
	for _, score := range players {
		if score > max {
			max = score
		}
	}
	return max
}

func PartTwo(input Input) int {
	var (
		circle  = NewCircle()
		players = make([]int, input.NumPlayers)
		marbles = 100 * input.NumMarbles
	)
	for i := 1; i <= marbles; i++ {
		player := i % input.NumPlayers
		if i != 0 && i%23 == 0 {
			players[player] += circle.Score(i)
		} else {
			circle.Place(i)
		}
	}
	var max int
	for _, score := range players {
		if score > max {
			max = score
		}
	}
	return max
}

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer (Part 1): %d\n", PartOne(input))
	fmt.Printf("Answer (Part 2): %d\n", PartTwo(input))
}
