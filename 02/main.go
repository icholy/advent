package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ToRunes(s string) []rune {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		rr = append(rr, r)
	}
	return rr
}

func NumDiff(a, b string) int {
	ar := ToRunes(a)
	br := ToRunes(b)
	if len(ar) != len(br) {
		return -1
	}
	var diff int
	for i := range br {
		if ar[i] != br[i] {
			diff++
		}
	}
	return diff
}

func Common(a, b string) string {
	ar := ToRunes(a)
	br := ToRunes(b)
	if len(ar) != len(br) {
		return ""
	}
	var common []rune
	for i := range ar {
		if ar[i] == br[i] {
			common = append(common, ar[i])
		}
	}
	return string(common)
}

func RuneCount(s string) map[rune]int {
	m := map[rune]int{}
	for _, r := range s {
		m[r]++
	}
	return m
}

type Checksum struct {
	n2, n3 int
}

func (c *Checksum) Update(s string) {
	var has2, has3 bool
	for _, n := range RuneCount(s) {
		switch n {
		case 2:
			has2 = true
		case 3:
			has3 = true
		}
	}
	if has2 {
		c.n2++
	}
	if has3 {
		c.n3++
	}
}

func (c Checksum) Sum() int {
	return c.n2 * c.n3
}

// Table indexes strings by their character counts it provides
// a not-so-great way of finding similar strings
type Table struct {
	m map[rune]map[int][]string
}

func (t *Table) runeMap(r rune) map[int][]string {
	if t.m == nil {
		t.m = make(map[rune]map[int][]string)
	}
	rm, ok := t.m[r]
	if !ok {
		rm = make(map[int][]string)
		t.m[r] = rm
	}
	return rm
}

func (t *Table) Insert(s string) {
	for r, n := range RuneCount(s) {
		rm := t.runeMap(r)
		rm[n] = append(rm[n], s)
	}
}

func (t *Table) Lookup(r rune, n int) []string {
	rm := t.runeMap(r)
	return rm[n]
}

func (t *Table) Similar(s string) (string, bool) {
	for r, n := range RuneCount(s) {
		for _, s0 := range t.Lookup(r, n) {
			if NumDiff(s, s0) == 1 {
				return s0, true
			}
		}
	}
	return "", false
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var chk Checksum
	var tbl Table
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		id := sc.Text()

		// update the checksum
		chk.Update(id)

		// check for a similar id
		if s, ok := tbl.Similar(id); ok {
			fmt.Printf("Common: %s\n", Common(id, s))
		}
		tbl.Insert(id)
	}
	fmt.Printf("Checksum: %d\n", chk.Sum())
}
