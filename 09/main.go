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
}

func (m *Marble) InsertAfter(next *Marble) {
	next.Next = m.Next
	m.Next = next
}

func (m *Marble) Delete() {
	*m = *m.Next
}

type Circle struct {
	Size    int
	First   *Marble
	Current *Marble
}

func NewCircle() *Circle {
	m := &Marble{Num: 0}
	m.Next = m
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
	return c.Clockwise(c.Size - (n % c.Size))
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
	score := m.Num
	m.Delete()
	c.Current = m
	c.Size--
	return score
}

func (c *Circle) Place(marble int) {
	m := &Marble{Num: marble}
	c.Clockwise(1).InsertAfter(m)
	c.Current = m
	c.Size++
}

func main() {
	input, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", input)

	var (
		circle  = NewCircle()
		players = make([]int, input.NumPlayers)
	)

	fmt.Printf("[-] %s\n", circle)
	for i := 1; i <= input.NumMarbles; i++ {
		player := i % input.NumPlayers
		if i != 0 && i%23 == 0 {
			players[player] += circle.Score(i)
		} else {
			circle.Place(i)
		}
		fmt.Printf("[player=%d, marble=%d] %s\n", player, i, circle)
	}

	for i, score := range players {
		fmt.Printf("Player %d: %d\n", i, score)
	}
}
