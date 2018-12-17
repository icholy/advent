package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"sort"

	"github.com/icholy/draw"
)

func main() {
	cv, err := ReadCanvas("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	tt := ParseTracks(cv)
	sim := NewSimulation(tt)
	if err := sim.Validate(); err != nil {
		log.Fatal(err)
	}
	for sim.NumCarts() > 1 {
		sim.Tick()
	}
	fmt.Printf("Answer (Part One): %s\n", sim.Collisions[0])
	for _, c := range sim.Carts {
		if !c.Crashed {
			fmt.Printf("Answer (Part Two): %s\n", c.Position)
		}
	}
}

type Direction byte

func (d Direction) Valid() bool {
	switch d {
	case North, South, East, West:
		return true
	default:
		return false
	}
}

func (d Direction) Turn(t Turn) Direction {
	if t != Left && t != Right {
		return d
	}
	switch d {
	case North:
		if t == Left {
			return West
		} else {
			return East
		}
	case South:
		if t == Left {
			return East
		} else {
			return West
		}
	case East:
		if t == Left {
			return North
		} else {
			return South
		}
	case West:
		if t == Left {
			return South
		} else {
			return North
		}
	}
	return d
}

const (
	North = Direction('^')
	South = Direction('v')
	East  = Direction('>')
	West  = Direction('<')
)

type Turn int

const (
	Left = Turn(iota)
	Strait
	Right
)

func (t Turn) Next() Turn {
	switch t {
	case Left:
		return Strait
	case Strait:
		return Right
	case Right:
		return Left
	default:
		panic("invalid turn")
	}
}

type Cart struct {
	Direction Direction
	Position  image.Point
	Track     *Track
	NextTurn  Turn
	Crashed   bool
}

func (c *Cart) LessThan(other *Cart) bool {
	if c.Position.Y < other.Position.Y {
		return true
	}
	return c.Position.X < other.Position.X
}

func (c *Cart) Intersect(other *Track) {
	if c.NextTurn != Strait {
		c.Track = other
		c.Direction = c.Direction.Turn(c.NextTurn)
	}
	c.NextTurn = c.NextTurn.Next()
}

func (c *Cart) Step() {

	switch c.Position {
	case c.Track.TopLeft():
		if c.Direction == North {
			c.Direction = East
		} else {
			c.Direction = South
		}
	case c.Track.TopRight():
		if c.Direction == East {
			c.Direction = South
		} else {
			c.Direction = West
		}
	case c.Track.BottomLeft():
		if c.Direction == West {
			c.Direction = North
		} else {
			c.Direction = East
		}
	case c.Track.BottomRight():
		if c.Direction == East {
			c.Direction = North
		} else {
			c.Direction = West
		}
	}

	switch c.Direction {
	case North:
		c.Position.Y--
	case South:
		c.Position.Y++
	case East:
		c.Position.X++
	case West:
		c.Position.X--
	default:
		panic("invalid position")
	}
}

func (c Cart) Draw(cv draw.Canvas, _ byte) {
	p := draw.FromImagePoint(c.Position)
	cv.Draw(p, byte(c.Direction))
}

type Track struct {
	Carts      []*Cart
	Intersects []image.Point
	Rect       image.Rectangle
}

func (t Track) TopLeft() image.Point     { return t.Rect.Min }
func (t Track) TopRight() image.Point    { return image.Pt(t.Rect.Max.X, t.Rect.Min.Y) }
func (t Track) BottomLeft() image.Point  { return image.Pt(t.Rect.Min.X, t.Rect.Max.Y) }
func (t Track) BottomRight() image.Point { return t.Rect.Max }

func (t Track) Draw(cv draw.Canvas, b byte) {
	r := draw.FromImageRect(t.Rect)
	cv.Draw(r.Left(), '|')
	cv.Draw(r.Right(), '|')
	cv.Draw(r.Top(), '-')
	cv.Draw(r.Bottom(), '-')
	cv.Draw(r.TopLeft(), '/')
	cv.Draw(r.TopRight(), '\\')
	cv.Draw(r.BottomLeft(), '\\')
	cv.Draw(r.BottomRight(), '/')
}

func ParseTracks(cv draw.Canvas) []*Track {
	var tt []*Track
	for x := 0; x < cv.Width(); x++ {
		for y := 0; y < cv.Height(); y++ {
			// top left of the rect
			if cv.At(x, y) == '/' && cv.Contains(x, y+1) && (cv.At(x, y+1) == '|' || cv.At(x, y+1) == '+') {
				tt = append(tt, ParseTrack(cv, x, y))
			}
		}
	}
	return tt
}

// ParseRect takes a canvas and a top-left corner of the rect
func ParseTrack(cv draw.Canvas, x, y int) *Track {
	min := image.Pt(x, y)
	max := min

	// move right until we find '\\'
	for cv.At(max.X, max.Y) != '\\' {
		max.X++
	}
	// move down until we find '/'
	for cv.At(max.X, max.Y) != '/' {
		max.Y++
	}
	rect := image.Rectangle{min, max}

	// find the intersect points and carts
	track := &Track{Rect: rect}

	// top row
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		b := cv.At(x, rect.Min.Y)
		p := image.Pt(x, rect.Min.Y)
		if b == '+' {
			track.Intersects = append(track.Intersects, p)
		}
		if d := Direction(b); d.Valid() {
			track.Carts = append(track.Carts, &Cart{d, p, track, Left, false})
		}
	}
	// bottom row
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		b := cv.At(x, rect.Max.Y)
		p := image.Pt(x, rect.Max.Y)
		if b == '+' {
			track.Intersects = append(track.Intersects, p)
		}
		if d := Direction(b); d.Valid() {
			track.Carts = append(track.Carts, &Cart{d, p, track, Left, false})
		}
	}
	// left side
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		b := cv.At(rect.Min.X, y)
		p := image.Pt(rect.Min.X, y)
		if b == '+' {
			track.Intersects = append(track.Intersects, p)
		}
		if d := Direction(b); d.Valid() {
			track.Carts = append(track.Carts, &Cart{d, p, track, Left, false})
		}
	}
	// right side
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		b := cv.At(rect.Max.X, y)
		p := image.Pt(rect.Max.X, y)
		if b == '+' {
			track.Intersects = append(track.Intersects, p)
		}
		if d := Direction(b); d.Valid() {
			track.Carts = append(track.Carts, &Cart{d, p, track, Left, false})
		}
	}

	return track
}

type Simulation struct {
	Intersects map[image.Point][]*Track
	Carts      []*Cart
	Tracks     []*Track
	Collisions []image.Point
	Occupied   map[image.Point]*Cart
}

func NewSimulation(tracks []*Track) *Simulation {
	sim := &Simulation{
		Intersects: make(map[image.Point][]*Track),
		Tracks:     tracks,
		Occupied:   make(map[image.Point]*Cart),
	}
	for _, t := range tracks {
		for _, p := range t.Intersects {
			sim.Intersects[p] = append(sim.Intersects[p], t)
		}
		sim.Carts = append(sim.Carts, t.Carts...)
	}
	for _, c := range sim.Carts {
		sim.Occupied[c.Position] = c
	}
	return sim
}

func (sim *Simulation) NumCarts() int {
	return len(sim.Occupied)
}

func (sim *Simulation) OtherTrack(c *Cart) (*Track, bool) {
	if tt, ok := sim.Intersects[c.Position]; ok {
		for _, t := range tt {
			if t != c.Track {
				return t, true
			}
		}
	}
	return nil, false
}

func (sim *Simulation) Tick() {
	sort.Slice(sim.Carts, func(i, j int) bool {
		return sim.Carts[i].LessThan(sim.Carts[j])
	})

	for _, c := range sim.Carts {

		if c.Crashed {
			continue
		}

		delete(sim.Occupied, c.Position)
		if t, ok := sim.OtherTrack(c); ok {
			c.Intersect(t)
		}
		c.Step()

		if other, ok := sim.Occupied[c.Position]; ok {
			c.Crashed = true
			other.Crashed = true
			sim.Collisions = append(sim.Collisions, c.Position)
			delete(sim.Occupied, c.Position)
		} else {
			sim.Occupied[c.Position] = c
		}
	}
}

func (sim *Simulation) TickN(n int) {
	for i := 0; i < n; i++ {
		sim.Tick()
	}
}

func (sim *Simulation) Validate() error {
	for _, tt := range sim.Intersects {
		if len(tt) != 2 {
			return fmt.Errorf("invalid intersect: %#v", len(tt))
		}
	}
	return nil
}

func (sim *Simulation) Draw(cv draw.Canvas, _ byte) {
	for _, t := range sim.Tracks {
		cv.Draw(t, 0)
	}
	for p, _ := range sim.Intersects {
		cv.Draw(draw.FromImagePoint(p), '+')
	}
	for _, p := range sim.Collisions {
		cv.Draw(draw.FromImagePoint(p), 'X')
	}
	for _, c := range sim.Carts {
		if !c.Crashed {
			cv.Draw(c, 0)
		}
	}
}

func ReadCanvas(file string) (draw.Canvas, error) {
	cv := draw.NewCanvas(150, 150)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := cv.ReadFrom(f); err != nil {
		return nil, err
	}
	return cv, nil
}
