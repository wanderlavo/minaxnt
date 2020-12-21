package miner

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	TICK_SECONDS time.Duration = 10
)

type Stat struct {
	counter     uint64
	lastCounter uint64
	lastTime    time.Time
}

func (s *Stat) Start() {
	s.lastTime = time.Now().UTC()
	s.lastCounter = 0
	s.counter = 0

	var now time.Time
	var timeDiff time.Duration
	var currentCounter uint64
	var counterDiff uint64
	var rate float64
	t := time.NewTicker(TICK_SECONDS * time.Second)
	for {
		select {
		case <-t.C:
			now = time.Now().UTC()
			timeDiff = now.Sub(s.lastTime)
			currentCounter = s.counter
			counterDiff = currentCounter - s.lastCounter
			rate = float64(counterDiff) / timeDiff.Seconds()

			log.Infof("Total Speed: %s, Time: %.1fs", s.humanizeRate(rate), timeDiff.Seconds())

			s.lastCounter = currentCounter
			s.lastTime = now
		}
	}
}

func (s *Stat) Incr() {
	atomic.AddUint64(&s.counter, 1)
}

func (s *Stat) humanizeRate(rate float64) string {
	var hr string
	if rate/1000.0 <= 1.0 {
		hr = fmt.Sprintf("%.1f Work/s", rate)
	} else if rate/1000000.0 <= 1.0 {
		hr = fmt.Sprintf("%.1f KWork/s", rate/1000.0)
	} else if rate/1000000000.0 <= 1.0 {
		hr = fmt.Sprintf("%.1f MWork/s", rate/1000000.0)
	} else {
		hr = fmt.Sprintf("%.1f GWork/s", rate/1000000000.0)
	}
	return hr
}