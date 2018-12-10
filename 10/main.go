package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/icholy/draw"
	"github.com/spakin/disjoint"
)

type Light struct {
	Pos, Vel image.Point
	Reaching *disjoint.Element
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (l Light) Touches(other Light) bool {
	delta := l.Pos.Sub(other.Pos)
	return Abs(delta.X) <= 1 && Abs(delta.Y) <= 1
}

func (l Light) Position(seconds int) image.Point {
	return l.Pos.Add(l.Vel.Mul(seconds))
}

func (l Light) String() string {
	return fmt.Sprintf("pos=%s vel=%s", l.Pos, l.Vel)
}

func ReadInput(file string) ([]Light, error) {
	var lights []Light
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var l Light
		_, err := fmt.Sscanf(sc.Text(),
			"position=<%d, %d> velocity=<%d, %d>",
			&l.Pos.X, &l.Pos.Y, &l.Vel.X, &l.Vel.Y,
		)
		if err != nil {
			return nil, err
		}
		l.Reaching = disjoint.NewElement()
		lights = append(lights, l)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return lights, nil
}

func Bounds(lights []Light) image.Rectangle {
	var b image.Rectangle
	for i, l := range lights {
		r := image.Rectangle{l.Pos, l.Pos}
		if i == 0 {
			b = r
		}
		b = b.Union(r)
	}
	return b
}

func Search(lights []Light) int {
	var seconds int
	smallest := len(lights)
	for i := 0; i < 1000000; i++ {
		groups := Groups(Simulate(lights, i))
		if n := len(groups); n < smallest {
			smallest = n
			seconds = i
		}
	}
	return seconds
}

func Groups(lights []Light) [][]Light {
	for _, l := range lights {
		*l.Reaching = *disjoint.NewElement()
	}
	for i, l1 := range lights {
		for j, l2 := range lights {
			if i != j && l1.Touches(l2) {
				disjoint.Union(l1.Reaching, l2.Reaching)
			}
		}
	}
	groups := map[*disjoint.Element][]Light{}
	for _, l := range lights {
		root := l.Reaching.Find()
		groups[root] = append(groups[root], l)
	}
	var ll [][]Light
	for _, g := range groups {
		ll = append(ll, g)
	}
	return ll
}

func Simulate(lights []Light, seconds int) []Light {
	simulated := make([]Light, len(lights))
	for i, l := range lights {
		l.Pos = l.Position(seconds)
		simulated[i] = l
	}
	return simulated
}

func Draw(lights []Light) error {
	cv := draw.NewCanvas(300, 300)
	cv.Draw(cv.Bounds().Fill(), '.')
	for i, g := range Groups(lights) {
		for _, l := range g {
			cv.Draw(draw.FromImagePoint(l.Pos), 'A'+byte(i))
		}
	}
	return cv.WriteTo(os.Stdout)
}

func main() {
	lights, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	Draw(Simulate(lights, Search(lights)))
}
