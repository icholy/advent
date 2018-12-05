package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"unicode"
)

type RunePredicate func(rune) bool

func IsValid(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func Without(bad rune) RunePredicate {
	return func(r rune) bool {
		return r != bad
	}
}

func Filter(rr []rune, preds ...RunePredicate) []rune {
	var out []rune
	for _, r := range rr {
		allow := true
		for _, p := range preds {
			if !p(r) {
				allow = false
				break
			}
		}
		if allow {
			out = append(out, r)
		}
	}
	return out
}

func ToRunes(s string) []rune {
	var rr []rune
	for _, r := range s {
		rr = append(rr, r)
	}
	return rr
}

func Cancels(a, b rune) bool {
	return a != b && unicode.ToUpper(a) == unicode.ToUpper(b)
}

func ReduceRunes(in []rune) ([]rune, bool) {
	var (
		out  []rune
		size = len(in)
		done = true
	)
	for i := 0; i < size; i++ {
		ch := in[i]
		if j := i + 1; j < size && Cancels(ch, in[j]) {
			i++
			done = false
		} else {
			out = append(out, ch)
		}
	}
	return out, done
}

func Reduce(s string, preds ...RunePredicate) string {
	rr := ToRunes(s)
	rr = Filter(rr, preds...)
	var done bool
	for !done {
		rr, done = ReduceRunes(rr)
	}
	return string(rr)
}

func PartOne(s string) int {
	return len(Reduce(s, IsValid))
}

func PartTwo(s string) int {
	min := len(Reduce(s, IsValid))
	for ch := 'a'; ch <= 'z'; ch++ {
		reduced := Reduce(s,
			IsValid,
			Without(ch),
			Without(unicode.ToUpper(ch)),
		)
		if size := len(reduced); size < min {
			min = size
		}
	}
	return min
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	polymers := string(data)
	fmt.Printf("Answer (Part 1): %d\n", PartOne(polymers))
	fmt.Printf("Answer (Part 2): %d\n", PartTwo(polymers))
}
