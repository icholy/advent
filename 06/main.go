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

func Nearest(point image.Point, points []image.Point) int {
	var (
		closest  int
		distance = -1
	)
	for i, p := range points {
		d := Distance(point, p)
		switch {
		case d == distance:
			closest = -1
		case d < distance || distance == -1:
			closest = i
			distance = d
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

func IndexByte(i int) byte {
	return '0' + byte(i)
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

func PartOne(coords []image.Point) int {
	areas := map[image.Point]int{}
	Iterate(Bounds(coords), func(p image.Point) {
		if i := Nearest(p, coords); i > 0 {
			areas[coords[i]]++
		}
	})
	var max int
	for c, area := range areas {
		if area > max && IsFinite(c, coords) {
			max = area
		}
	}
	return max
}

func Draw(coords []image.Point) error {
	bounds := image.Rectangle{
		Min: image.ZP,
		Max: Bounds(coords).Max,
	}
	cv := draw.NewCanvas(bounds.Dx()+1, bounds.Dy()+1)
	cv.Draw(cv.Bounds().Fill(), '.')

	finite := map[int]bool{}
	for i, p := range coords {
		finite[i] = IsFinite(p, coords)
	}

	Iterate(cv.Bounds().Image(), func(p image.Point) {
		if i := Nearest(p, coords); i != -1 {
			if finite[i] {
				cv.Draw(draw.FromImagePoint(p), IndexByte(i))
			} else {
				cv.Draw(draw.FromImagePoint(p), '*')
			}
		}
	})

	for _, c := range coords {
		cv.Draw(draw.FromImagePoint(c), '*')
	}

	return cv.WriteTo(os.Stdout)
}

func main() {
	coords, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Answer (Part 1): %d\n", PartOne(coords))

	if err := Draw(coords); err != nil {
		log.Fatal(err)
	}
}
