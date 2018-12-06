package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/icholy/draw"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Distance(a, b image.Point) int {
	delta := a.Sub(b)
	return Abs(delta.X) + Abs(delta.Y)
}

func IsNearest(center, other image.Point, points []image.Point) bool {
	dist := Distance(center, other)
	for _, p := range points {
		if Distance(p, center) <= dist {
			return false
		}
	}
	return true
}

func Nearest(point image.Point, points []image.Point) int {
	var (
		closest int
		dist    = -1
	)
	for i, p := range points {
		if p == point {
			continue
		}
		d := Distance(point, p)
		switch {
		case d == dist:
			closest = -1
		case d < dist || dist == -1:
			closest = i
			dist = d
		}
	}
	return closest
}

func IsFinite(point image.Point, points []image.Point) bool {
	var above, below, left, right bool
	for _, p := range points {
		if p.X < point.X {
			left = true
		}
		if p.X > point.X {
			right = true
		}
		if p.Y < point.Y {
			above = true
		}
		if p.Y > point.Y {
			below = true
		}
		if above && below && left && right {
			return true
		}
	}
	return false
}

func ReadInput(file string) ([]image.Point, error) {
	var pp []image.Point
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var p image.Point
		_, err := fmt.Sscanf(sc.Text(), "%d, %d", &p.X, &p.Y)
		if err != nil {
			return nil, err
		}
		pp = append(pp, p)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return pp, nil
}

func UpperByte(i int) byte {
	if i > 26 {
		return '?'
	}
	return 'A' + byte(i)
}

func LowerByte(i int) byte {
	if i > 26 {
		return '?'
	}
	return 'a' + byte(i)
}

func Max(a, b image.Point) image.Point {
	if b.X > a.X || a.X == -1 {
		a.X = b.X
	}
	if b.Y > a.Y || a.Y == -1 {
		a.Y = b.Y
	}
	return a
}

func Min(a, b image.Point) image.Point {
	if b.X < a.X || a.X == -1 {
		a.X = b.X
	}
	if b.Y < a.Y || a.Y == -1 {
		a.Y = b.Y
	}
	return a
}

var Unset = image.Point{-1, -1}

func Bounds(points []image.Point) image.Rectangle {
	var min, max = Unset, Unset
	for _, p := range points {
		min = Min(min, p)
		max = Max(max, p)
	}
	return image.Rectangle{min, max}
}

func Iterate(r image.Rectangle, f func(image.Point)) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			f(image.Pt(x, y))
		}
	}
}

func Area(center, other image.Point, coords []image.Point, bounds image.Rectangle) int {
	if !IsFinite(center, coords) {
		return -1
	}
	var area int
	Iterate(bounds, func(p image.Point) {
		if IsNearest(center, p, coords) {
			area++
		}
	})
	return area
}

func main() {
	coords, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	cv := draw.NewCanvas(50, 20)
	cv.Draw(cv.Bounds().Fill(), '.')

	Iterate(cv.Bounds().Image(), func(p image.Point) {
		if i := Nearest(p, coords); i != -1 {
			cv.Draw(draw.FromImagePoint(p), LowerByte(i))
		}
	})

	for i, c := range coords {
		if IsFinite(c, coords) {
			cv.Draw(draw.FromImagePoint(c), '*')
		} else {
			cv.Draw(draw.FromImagePoint(c), UpperByte(i))
		}
	}

	cv.Draw(draw.FromImageRect(Bounds(coords)), '*')

	areas := map[int]int{}
	Iterate(Bounds(coords), func(p image.Point) {
		i := Nearest(p, coords)
		areas[i]++
	})

	for i, area := range areas {
		fmt.Println(string([]byte{UpperByte(i)}), area)
	}

	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
