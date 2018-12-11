package main

import (
	"fmt"
	"image"
)

func main() {
	serial := 7672
	fmt.Printf("Answer (Part 1): %s\n", PartOne(serial))
	fmt.Printf("Answer (Part 2): %s\n", PartTwo(serial))
}

func PartOne(serial int) image.Point {
	var (
		bestpower int
		bestcell  image.Point
	)
	ForEachKernel(3, image.Rect(1, 1, 300, 300), func(kernel image.Rectangle) {
		var power int
		ForEachCell(kernel, func(cell image.Point) {
			power += Power(cell, serial)
		})
		if power > bestpower {
			bestpower = power
			bestcell = kernel.Min
		}
	})
	return bestcell
}

func PartTwo(serial int) PartTwoResult {
	var (
		result = make(chan PartTwoResult)
		best   PartTwoResult
	)

	go PartTwoRange(serial, 1, 50, result)
	go PartTwoRange(serial, 51, 100, result)
	go PartTwoRange(serial, 101, 150, result)
	go PartTwoRange(serial, 151, 200, result)
	go PartTwoRange(serial, 201, 250, result)
	go PartTwoRange(serial, 251, 300, result)

	for i := 0; i < 6; i++ {
		if res := <-result; res.Power > best.Power {
			best = res
		}
	}

	return best
}

type PartTwoResult struct {
	Point image.Point
	Power int
	Size  int
}

func (r PartTwoResult) String() string {
	return fmt.Sprintf("%s, %d", r.Point, r.Size)
}

func PartTwoRange(serial, min, max int, out chan PartTwoResult) {
	var (
		bestpower int
		bestsize  int
		bestcell  image.Point
	)
	for size := min; size <= max; size++ {
		ForEachKernel(size, image.Rect(1, 1, 300, 300), func(kernel image.Rectangle) {
			var power int
			ForEachCell(kernel, func(cell image.Point) {
				power += Power(cell, serial)
			})
			if power > bestpower {
				bestpower = power
				bestcell = kernel.Min
				bestsize = size
			}
		})
	}
	out <- PartTwoResult{bestcell, bestpower, bestsize}
}

func ForEachKernel(size int, bounds image.Rectangle, f func(r image.Rectangle)) {
	kernel := image.Rect(0, 0, size-1, size-1).Add(bounds.Min)
	for kernel.Max.Y <= bounds.Max.Y {
		f(kernel)
		kernel = kernel.Add(image.Pt(1, 0))
		if kernel.Max.X > bounds.Max.X {
			dx := kernel.Min.X - bounds.Min.X
			kernel = kernel.Add(image.Pt(-dx, 1))
		}
	}
}

func ForEachCell(r image.Rectangle, f func(p image.Point)) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			f(image.Pt(x, y))
		}
	}
}

func Power(cell image.Point, serial int) int {
	// Find the fuel cell's rack ID, which is its X coordinate plus 10.
	rackID := cell.X + 10

	// Begin with a power level of the rack ID times the Y coordinate.
	power := rackID * cell.Y

	// Increase the power level by the value of the grid serial number (your puzzle input)
	power += serial

	// Set the power level to itself multiplied by the rack ID.
	power *= rackID

	// Keep only the hundreds digit of the power level (so 12345 becomes 3; numbers with no hundreds digit become 0).
	power = (power / 100) % 10

	// Subtract 5 from the power level.
	power -= 5

	return power
}
