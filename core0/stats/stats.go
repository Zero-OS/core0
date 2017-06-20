package stats

import (
	"encoding/json"
	"time"
)

const (
	Average      Operation = "A"
	Differential Operation = "D"
)

type Operation string

type Sample struct {
	Avg   float64 `json:"avg"`
	Total float64 `json:"total"`
	Max   float64 `json:"max"`
	Count uint    `json:"count"`
	Start int64   `json:"start"`
}

/*
Feed values on now, for the specific aggregate duration

*/
func (m *Sample) Feed(value float64, now int64, duration int64) *Sample {
	period := (now / duration) * duration

	if period != 0 && m.Start < period {
		//start a new period
		update := *m

		m.Total = value
		m.Avg = value
		m.Max = value
		m.Count = 1
		m.Start = period

		return &update
	}

	if m.Start == 0 {
		m.Start = period
	}

	m.Total += value
	m.Count += 1
	m.Avg = m.Total / float64(m.Count)
	if value > m.Max {
		m.Max = value
	}

	return nil
}

type Samples map[int64]*Sample

type State struct {
	Operation Operation `json:"op"`
	LastValue float64   `json:"last_value"`
	LastTime  int64     `json:"last_time"`
	Tags      string    `json:"tags,omitempty"`
	Samples   Samples   `json:"samples"`
}

func NewState(op Operation, durations ...int64) *State {
	s := State{
		Operation: op,
		Samples:   Samples{},
		LastTime:  -1,
	}

	for _, d := range durations {
		s.Samples[d] = &Sample{}
	}

	return &s
}

func LoadState(data []byte) (*State, error) {
	var state State
	return &state, json.Unmarshal(data, &state)
}

func (s *State) avg(now int64, value float64) {
	for d, sample := range s.Samples {
		sample.Feed(value, now, d)
	}
}

func (s *State) init(now int64, value float64) {
	for d, sample := range s.Samples {
		if s.Operation == Average {
			sample.Feed(value, now, d)
		}
	}
}

func (s *State) FeedOn(now int64, value float64) Samples {
	orig := value
	defer func() {
		s.LastValue = orig
		s.LastTime = now
	}()

	if s.LastTime == -1 {
		s.init(now, value)
		return nil
	}

	if s.Operation == Differential {
		value = (value - s.LastValue) / float64(now-s.LastTime)
	}

	updates := Samples{}
	for d, sample := range s.Samples {
		if update := sample.Feed(value, now, d); update != nil {
			updates[d] = update
		}
	}

	return updates
}

func (s *State) Feed(value float64) Samples {
	return s.FeedOn(time.Now().Unix(), value)
}
