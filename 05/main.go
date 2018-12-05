package main

import "fmt"
import "io/ioutil"
import "log"
import "unicode"

func IsValid(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func Clean(rr []rune) []rune {
	var out []rune
	for _, r := range rr {
		if IsValid(r) {
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
		out     []rune
		size    = len(in)
		changed bool
	)
	for i := 0; i < size; i++ {
		ch := in[i]
		if j := i + 1; j < size && Cancels(ch, in[j]) {
			i++
			changed = true
		} else {
			out = append(out, ch)
		}
	}
	return out, changed
}

func Reduce(s string) string {
	var (
		rr      = Clean(ToRunes(s))
		changed = true
	)
	for changed {
		rr, changed = ReduceRunes(rr)
	}
	return string(rr)
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	polymers := string(data)
	reduced := Reduce(polymers)
	fmt.Println(len(reduced))
}
