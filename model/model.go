// Package model define the structs of data model, request, or response object
package model

import "time"

type Pomodoro struct {
	ID              int64
	Type            string
	State           int
	Starttime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
}

type SwitchStateResponse struct {
	DurationString  string
	NextSessionType string
	Title           string
}

type Activity struct {
	ID          int64
	Message     string
	CompletedAt time.Time
}

type Summary struct {
	CurrentDate   time.Time
	ActivityCount int
	SessionCount  int
	FocusDuration time.Duration
}
