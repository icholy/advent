package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type RecordType string

const (
	Begin   RecordType = "BEGIN"
	Sleep   RecordType = "SLEEP"
	Wake    RecordType = "WAKE"
	Invalid RecordType = "INVALID"
)

type Record struct {
	Time    time.Time
	Text    string
	Type    RecordType
	GuardID int
}

func (r Record) String() string {
	return fmt.Sprintf("%s: %s %d", r.Time, r.Type, r.GuardID)
}

var (
	recordRe = regexp.MustCompile(`\[(.*)\] (.*)`)
	beginRe  = regexp.MustCompile(`Guard #(\d+) begins shift`)
	layout   = "2006-01-02 15:04"
)

func ParseRecordText(s string) (RecordType, int) {
	switch s {
	case "wakes up":
		return Wake, 0
	case "falls asleep":
		return Sleep, 0
	default:
		m := beginRe.FindStringSubmatch(s)
		if len(m) != 2 {
			return Invalid, 0
		}
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return Invalid, 0
		}
		return Begin, id
	}
}

func ParseRecord(s string) (Record, error) {
	m := recordRe.FindStringSubmatch(s)
	if len(m) != 3 {
		return Record{}, fmt.Errorf("invalid record: %s", s)
	}
	date := m[1]
	t, err := time.Parse(layout, date)
	if err != nil {
		return Record{}, err
	}
	text := m[2]
	typ, id := ParseRecordText(text)
	return Record{
		Time:    t,
		Text:    text,
		Type:    typ,
		GuardID: id,
	}, nil
}

func ParseRecords(r io.Reader) ([]Record, error) {
	var rr []Record
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		rec, err := ParseRecord(sc.Text())
		if err != nil {
			return nil, err
		}
		rr = append(rr, rec)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(rr, func(i, j int) bool {
		return rr[i].Time.Before(rr[j].Time)
	})
	return rr, nil
}

type TimeRange struct {
	Start, End time.Time
}

func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

type Guard struct {
	ID         int
	Shifts     int
	Sleeps     []TimeRange
	TotalSleep time.Duration
}

func (g *Guard) String() string {
	return fmt.Sprintf("%d: shifts=%d slept=%s ", g.ID, g.Shifts, g.TotalSleep)
}

func (g *Guard) Sleep(start, end time.Time) {
	tr := TimeRange{start, end}
	g.Sleeps = append(g.Sleeps, tr)
	g.TotalSleep += tr.Duration()
}

// minute -> count
type MinuteCount map[int]int

func (mc MinuteCount) Max() (min, count int) {
	var minMax, nMax int
	for min, n := range mc {
		if n > nMax {
			nMax = n
			minMax = min
		}
	}
	return minMax, nMax
}

// hours -> minute count
type Histogram map[int]MinuteCount

func (h Histogram) Hour(hour int) MinuteCount {
	if _, ok := h[hour]; !ok {
		h[hour] = make(MinuteCount)
	}
	return h[hour]
}

func (h Histogram) Update(tr TimeRange) {
	for t := tr.Start; t.Before(tr.End); t = t.Add(time.Minute) {
		hour, minute := t.Hour(), t.Minute()
		m := h.Hour(hour)
		m[minute]++
	}
}

type Tracker struct {
	guards  map[int]*Guard
	current *Guard
	sleep   time.Time
}

func NewTracker() *Tracker {
	return &Tracker{
		guards: make(map[int]*Guard),
	}
}

func (t *Tracker) Guard(id int) *Guard {
	g, ok := t.guards[id]
	if !ok {
		g = &Guard{ID: id}
		t.guards[id] = g
	}
	return g
}

func (t *Tracker) Guards() []*Guard {
	var gg []*Guard
	for _, g := range t.guards {
		gg = append(gg, g)
	}
	sort.Slice(gg, func(i, j int) bool {
		return gg[i].TotalSleep < gg[j].TotalSleep
	})
	return gg
}

func (t *Tracker) Update(r Record) error {
	switch r.Type {
	case Begin:
		t.current = t.Guard(r.GuardID)
		t.current.Shifts++
	case Sleep:
		t.sleep = r.Time
	case Wake:
		t.current.Sleep(t.sleep, r.Time)
	default:
		return fmt.Errorf("invalid event")
	}
	return nil
}

func PartTwo(gg []*Guard) int {
	var maxCount, maxMinute, maxID int
	for _, g := range gg {
		hist := make(Histogram)
		for _, s := range g.Sleeps {
			hist.Update(s)
		}
		min, count := hist.Hour(0).Max()
		if count > maxCount {
			maxMinute = min
			maxCount = count
			maxID = g.ID
		}
	}
	return maxID * maxMinute
}

func PartOne(gg []*Guard) int {
	var worst *Guard
	for _, g := range gg {
		if worst == nil || g.TotalSleep > worst.TotalSleep {
			worst = g
		}
	}
	hist := make(Histogram)
	for _, s := range worst.Sleeps {
		hist.Update(s)
	}
	min, _ := hist.Hour(0).Max()
	return worst.ID * min
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rr, err := ParseRecords(f)
	if err != nil {
		log.Fatal(err)
	}
	t := NewTracker()
	for _, r := range rr {
		if err := t.Update(r); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Answer (Part 1): %d\n", PartOne(t.Guards()))
	fmt.Printf("Answer (Part 2): %d\n", PartTwo(t.Guards()))
}
