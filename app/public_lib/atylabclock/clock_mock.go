package atylabclock

import "time"

type ClockStructMock struct {
	now time.Time
}

func NewClockMock(now time.Time) *ClockStructMock {
	return &ClockStructMock{
		now: now,
	}
}

func (t *ClockStructMock) Now() time.Time {
	return t.now
}
