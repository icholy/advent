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
	GuardID int64
}

func (r Record) String() string {
	return fmt.Sprintf("%s: %s %d", r.Time, r.Type, r.GuardID)
}

var (
	recordRe = regexp.MustCompile(`\[(.*)\] (.*)`)
	beginRe  = regexp.MustCompile(`Guard #(\d+) begins shift`)
	layout   = "2006-01-02 15:04"
)

func ParseRecordText(s string) (RecordType, int64) {
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
		id, err := strconv.ParseInt(m[1], 10, 64)
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
		return rr[i].Time.Before(rr[i].Time)
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
	ID         int64
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

type Tracker struct {
	guards  map[int64]*Guard
	current *Guard
	sleep   time.Time
	worst   *Guard
}

func NewTracker() *Tracker {
	return &Tracker{
		guards: make(map[int64]*Guard),
	}
}

func (t *Tracker) Guard(id int64) *Guard {
	g, ok := t.guards[id]
	if !ok {
		g = &Guard{ID: id}
		t.guards[id] = g
	}
	return g
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

	// update the worst guard (one that sleeps the most)
	if t.current != nil {
		if t.worst == nil || t.current.TotalSleep > t.worst.TotalSleep {
			t.worst = t.current
		}
	}
	return nil
}

func (t *Tracker) Worst() *Guard { return t.worst }

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
	fmt.Println(t.Worst())
}
