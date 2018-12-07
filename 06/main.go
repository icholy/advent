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

func Nearest(point image.Point, points []image.Point) (image.Point, bool) {
	var (
		closest  image.Point
		found    bool
		distance = -1
	)
	for _, p := range points {
		d := Distance(point, p)
		switch {
		case d == distance:
			closest = p
			found = false
		case d < distance || distance == -1:
			closest = p
			distance = d
			found = true
		}
	}
	return closest, found
}

func IsFinite(point image.Point, points []image.Point) bool {
	var above, below, left, right bool
	for _, p := range points {
		delta := point.Sub(p)
		dx, dy := Abs(delta.X), Abs(delta.Y)
		if p.X < point.X && dy <= dx {
			left = true
		}
		if p.X > point.X && dy <= dx {
			right = true
		}
		if p.Y < point.Y && dx <= dy {
			above = true
		}
		if p.Y > point.Y && dx <= dy {
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

func Bounds(points []image.Point) image.Rectangle {
	var b image.Rectangle
	for i, p := range points {
		if i == 0 {
			b.Min = p
			b.Max = p
		}
		if p.X < b.Min.X {
			b.Min.X = p.X
		}
		if p.Y < b.Min.Y {
			b.Min.Y = p.Y
		}
		if p.X > b.Max.X {
			b.Max.X = p.X
		}
		if p.Y > b.Max.Y {
			b.Max.Y = p.Y
		}
	}
	return b
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
		if c, ok := Nearest(p, coords); ok {
			areas[c]++
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

func PartTwo(coords []image.Point) int {
	var area int
	Iterate(Bounds(coords), func(p image.Point) {
		var sum int
		for _, c := range coords {
			if d := Distance(p, c); d < 1000 {
				sum += d
			}
		}
		if sum < 10000 {
			area++
		}
	})
	return area
}

func Draw(coords []image.Point) error {
	bounds := image.Rectangle{
		Min: image.ZP,
		Max: Bounds(coords).Max,
	}
	cv := draw.NewCanvas(bounds.Dx()+1, bounds.Dy()+1)
	cv.Draw(cv.Bounds().Fill(), '.')

	finite := map[image.Point]bool{}
	for _, p := range coords {
		finite[p] = IsFinite(p, coords)
	}

	Iterate(cv.Bounds().Image(), func(p image.Point) {
		if c, ok := Nearest(p, coords); ok {
			if finite[c] {
				cv.Draw(draw.FromImagePoint(p), '$')
			} else {
				cv.Draw(draw.FromImagePoint(p), '%')
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
	if err := Draw(coords); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer (Part 1): %d\n", PartOne(coords))
	fmt.Printf("Answer (Part 2): %d\n", PartTwo(coords))
}
