// Package model define the structs of data model, request, or response object
package model

import "time"

type Pomodoro struct {
	ID             int64
	Type           string
	State          int
	Duration       time.Duration
	ActualDuration time.Duration
	CreatedAt      time.Time
}
