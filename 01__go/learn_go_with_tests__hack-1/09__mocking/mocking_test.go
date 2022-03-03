package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

const (
	write = "write"
	sleep = "sleep"
)

type SpySleepr struct {
	Calls int
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

func (ss *SpySleepr) Sleep() {
	ss.Calls++
}

type SleepCountdownOperations struct {
	Calls []string
}

func (sco *SleepCountdownOperations) Sleep() {
	sco.Calls = append(sco.Calls, sleep)
}

func (sco *SleepCountdownOperations) Write(p []byte) (n int, err error) {
	sco.Calls = append(sco.Calls, write)
	return
}

func TestCounddown(t *testing.T) {
	buffer := bytes.Buffer{}
	spySleeper := SpySleepr{}
	Countdown(&buffer, &spySleeper)
	want := `3
2
1
Go!`
	got := buffer.String()
	if want != got {
		t.Errorf("got %q, want %q", got, want)
	}

	if spySleeper.Calls != 4 {
		t.Errorf("not enough calls to a sleeper, want %d, got %d", 4, spySleeper.Calls)
	}

	t.Run("sleep before every print", func(t *testing.T) {
		SpySleepPrinter := &SleepCountdownOperations{}
		Countdown(SpySleepPrinter, SpySleepPrinter)
		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}
		if !reflect.DeepEqual(want, SpySleepPrinter.Calls) {
			t.Errorf("wanted calls %v, got %v", want, SpySleepPrinter.Calls)
		}
	})
}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second
	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
	sleeper.Sleep()
	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v, but slept for %v", sleepTime, spyTime.durationSlept)
	}
}
