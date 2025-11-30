package atylabclock

import "time"

type ClockInterface interface {
	Now() time.Time
}

type ClockStruct struct{}

func NewClock() *ClockStruct {
	return &ClockStruct{}
}

func (t *ClockStruct) Now() time.Time {
	return time.Now()
}
