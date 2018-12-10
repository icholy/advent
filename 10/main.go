package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/icholy/draw"
)

type Light struct {
	Pos, Vel image.Point
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
			fmt.Println(sc.Text())
			return nil, err
		}
		lights = append(lights, l)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return lights, nil
}

func Draw(lights []Light, seconds int) error {
	cv := draw.NewCanvas(40, 40)
	center := image.Pt(cv.Width()/2, cv.Height()/2)
	cv.Draw(cv.Bounds().Fill(), '.')
	for _, l := range lights {
		p := l.Position(seconds).Add(center)
		cv.Draw(draw.FromImagePoint(p), '#')
	}
	return cv.WriteTo(os.Stdout)
}

func main() {
	lights, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	if err := Draw(lights, 3); err != nil {
		log.Fatal(err)
	}
}
