// Package session contains pomodoro session management logic
package session

import (
	"context"
	"time"
)

type sessionState int

const (
	NotStarted sessionState = iota
	Running
	Paused
	Finished
)

const (
	PomodoroSession string = "pomodoro"
	BreakSession    string = "break"
)

type sessionService struct {
	repo        SessionRepository
	sessionType string
	duration    time.Duration
}

type Pomodoro struct {
	ID             int64
	Type           string
	State          int
	Duration       time.Duration
	ActualDuration time.Duration
	CreatedAt      time.Time
}

type SessionRepository interface {
	Last(context.Context) (Pomodoro, error)
	Add(context.Context, Pomodoro) error
	Update(context.Context, Pomodoro) error
}

func Newsession(repo SessionRepository, sessionType string, duration time.Duration) *sessionService {
	return &sessionService{
		repo:        repo,
		sessionType: sessionType,
		duration:    duration,
	}
}

func (s *sessionService) Start() error {
	return nil
}

func (s *sessionService) Pause() error {
	return nil
}

func (s *sessionService) Finish() error {
	return nil
}
