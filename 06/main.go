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

func Closest(index int, points []image.Point) int {
	var (
		point   = points[index]
		closest int
		dist    = -1
	)
	for i, p := range points {
		if i == index {
			continue
		}
		if d := Distance(point, p); d < dist || dist == -1 {
			closest = i
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

func IndexByte(i int) byte {
	if i > 26 {
		return '?'
	}
	return 'A' + byte(i)
}

func main() {
	coords, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	cv := draw.NewCanvas(50, 20)
	cv.Draw(cv.Bounds().Fill(), '.')
	for i, c := range coords {
		cv.Draw(draw.FromImagePoint(c), IndexByte(i))
	}
	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
