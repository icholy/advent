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
	return Abs(delta.X)+Abs(delta.Y) <= 2
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

func Groups(lights []Light) [][]Light {
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

func Draw(lights []Light) error {
	cv := draw.NewCanvas(40, 40)
	center := image.Pt(cv.Width()/2, cv.Height()/2)
	cv.Draw(cv.Bounds().Fill(), '.')
	for i, g := range Groups(lights) {
		for _, l := range g {
			p := l.Pos.Add(center)
			cv.Draw(draw.FromImagePoint(p), 'A'+byte(i))
		}
	}
	return cv.WriteTo(os.Stdout)
}

func Simulate(lights []Light, seconds int) []Light {
	simulated := make([]Light, len(lights))
	for i, l := range lights {
		l.Pos = l.Position(seconds)
		simulated[i] = l
	}
	return simulated
}

func main() {
	lights, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lights = Simulate(lights, 3)
	Draw(lights)
}
