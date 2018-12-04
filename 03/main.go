package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/icholy/draw"
)

// SplitX splits r into two rectangles along x.
// If x doesn't intersect r, r is returned unchanged
func SplitX(x int, r image.Rectangle) []image.Rectangle {
	if x <= r.Min.X || r.Max.X <= x {
		return []image.Rectangle{r}
	}
	return []image.Rectangle{
		image.Rect(r.Min.X, r.Min.Y, x, r.Max.Y),
		image.Rect(x, r.Min.Y, r.Max.X, r.Max.Y),
	}
}

// SplitY splits r into two rectangles along y.
// If y doesn't intersect r, r is returned unchanged
func SplitY(y int, r image.Rectangle) []image.Rectangle {
	if y <= r.Min.Y || r.Max.Y <= y {
		return []image.Rectangle{r}
	}
	return []image.Rectangle{
		image.Rect(r.Min.X, r.Min.Y, r.Max.X, y),
		image.Rect(r.Min.X, y, r.Max.X, r.Max.Y),
	}
}

// SplitAll splits every rectangle in rr by the point's x and y values.
// The ignore rectangle is omitted from the output
func SplitAll(p image.Point, rr []image.Rectangle, ignore image.Rectangle) []image.Rectangle {
	var ss []image.Rectangle
	for _, r := range rr {
		for _, r := range SplitX(p.X, r) {
			for _, r := range SplitY(p.Y, r) {
				if !r.Eq(ignore) {
					ss = append(ss, r)
				}
			}
		}
	}
	return ss
}

// Subtract a from b. This returns the parts of a that
// are not covered by a.
func Subtract(existing, b image.Rectangle) []image.Rectangle {
	rr := []image.Rectangle{b}
	if b.Overlaps(existing) {
		i := b.Intersect(existing)
		rr = SplitAll(i.Min, rr, i)
		rr = SplitAll(i.Max, rr, i)
		rr = SplitAll(image.Pt(i.Max.X, i.Min.Y), rr, i)
		rr = SplitAll(image.Pt(i.Min.X, i.Max.Y), rr, i)
	}
	return rr
}

// SubtractAll substracts all rects in a from b. This returns the parts of b
// that are not covered by any rects in a.
func SubtractAll(existing []image.Rectangle, b image.Rectangle) []image.Rectangle {
	for _, e := range existing {
		if b.Overlaps(e) {
			var unique []image.Rectangle
			for _, r := range Subtract(e, b) {
				unique = append(unique, SubtractAll(existing, r)...)
			}
			return unique
		}
	}
	return []image.Rectangle{b}
}

// Union appends the parts of r which are not already covered by the
// rectangles in rr
func Union(rr []image.Rectangle, r image.Rectangle) []image.Rectangle {
	return append(rr, SubtractAll(rr, r)...)
}

// Area returns the summed areas of rr
func Area(rr []image.Rectangle) int {
	var area int
	for _, r := range rr {
		area += r.Dx() * r.Dy()
	}
	return area
}

var inputRe = regexp.MustCompile(`#\d+ @ (\d+),(\d+): (\d+)x(\d+)`)

// ParseClaim parses an input claim line to a rectangle
func ParseClaim(s string) (image.Rectangle, error) {
	m := inputRe.FindStringSubmatch(s)
	if len(m) != 5 {
		return image.ZR, fmt.Errorf("no match")
	}
	x, err := strconv.Atoi(m[1])
	if err != nil {
		return image.ZR, err
	}
	y, err := strconv.Atoi(m[2])
	if err != nil {
		return image.ZR, err
	}
	w, err := strconv.Atoi(m[3])
	if err != nil {
		return image.ZR, err
	}
	h, err := strconv.Atoi(m[4])
	if err != nil {
		return image.ZR, err
	}
	min := image.Pt(x, y)
	return image.Rectangle{
		Min: min,
		Max: min.Add(image.Pt(w, h)),
	}, nil
}

// Draw the rectangles to stdout. This can get very big!
func Draw(rr []image.Rectangle) error {
	var maxX, maxY int
	for _, r := range rr {
		if r.Max.X > maxX {
			maxX = r.Max.X
		}
		if r.Max.Y > maxY {
			maxY = r.Max.Y
		}
	}
	cv := draw.NewCanvas(maxX+1, maxY+1)
	cv.Draw(cv.Bounds().Box(), 0)
	for i, r := range rr {
		cv.Draw(draw.FromImageRect(r), '0'+byte(i))
	}
	return cv.WriteTo(os.Stdout)
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// find the overlapped regions
	var claims, overlaps []image.Rectangle
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		claim, err := ParseClaim(sc.Text())
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range claims {
			if c.Overlaps(claim) {
				overlaps = Union(overlaps, c.Intersect(claim))
			}
		}
		claims = append(claims, claim)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	// find intact ID
	for i, c := range claims {
		intact := true
		for _, r := range overlaps {
			if c.Overlaps(r) {
				intact = false
				break
			}
		}
		if intact {
			fmt.Printf("Intact Claim: #%d\n", i+1)
		}
	}

	fmt.Printf("Overlap Area: %d\"\n", Area(overlaps))
}
