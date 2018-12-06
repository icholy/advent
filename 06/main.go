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

func Iterate(r image.Rectangle, f func(x, y int)) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			f(x, y)
		}
	}
}

func main() {
	coords, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	cv := draw.NewCanvas(50, 20)
	cv.Draw(cv.Bounds().Fill(), '.')

	Iterate(cv.Bounds().Image(), func(x, y int) {
		p := image.Pt(x, y)
		if i := Nearest(p, coords); i != -1 {
			cv.Draw(draw.FromImagePoint(p), LowerByte(i))
		}
	})

	for i, c := range coords {
		cv.Draw(draw.FromImagePoint(c), UpperByte(i))
	}

	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
