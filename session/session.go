// Package session contains pomodoro session management logic
package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/rinarudhei/pomcli/model"
)

type sessionState int

var ErrInvalidInput = errors.New("invalid input")

const (
	NotStarted sessionState = iota
	Running
	Paused
	Finished
)

const (
	PomodoroSession   string = "pomodoro"
	ShortBreakSession string = "shortbreak"
	LongBreakSession  string = "longbreak"
)

type SessionService struct {
	repo               SessionRepository
	sessionType        string
	pomodoroDuration   time.Duration
	shortBreakDuration time.Duration
	longBreakDuration  time.Duration
}

type SessionRepository interface {
	Last() (model.Pomodoro, error)
	Add() error
	Update() error
}

func NewSession(repo SessionRepository, pomodoroDuration, shortBreakDuration, longBreakDuration time.Duration) *SessionService {
	return &SessionService{
		repo:               repo,
		sessionType:        PomodoroSession,
		pomodoroDuration:   pomodoroDuration,
		shortBreakDuration: shortBreakDuration,
		longBreakDuration:  longBreakDuration,
	}
}

func (s *SessionService) Start() error {
	return nil
}

func (s *SessionService) Pause() error {
	return nil
}

func (s *SessionService) Finish() error {
	return nil
}

// SwitchState update session type, and return updated session type, duration in string, and title
func (s *SessionService) SwitchState() (string, string, string) {
	var durationString string
	var title string
	switch s.sessionType {
	case PomodoroSession:
		s.sessionType = ShortBreakSession
		durationString = durationToDisplayString(s.shortBreakDuration)
		title = "Short Break"
	case ShortBreakSession:
		s.sessionType = LongBreakSession
		durationString = durationToDisplayString(s.longBreakDuration)
		title = "Long Break"
	case LongBreakSession:
		s.sessionType = PomodoroSession
		durationString = durationToDisplayString(s.pomodoroDuration)
		title = "Pomodoro Focus"
	}

	return durationString, s.sessionType, title
}

func durationToDisplayString(d time.Duration) string {
	minutes := int(d.Truncate(time.Minute).Minutes())
	seconds := int(d.Seconds()) % 60
	minutesString := fmt.Sprintf("%d", minutes)
	secondsString := fmt.Sprintf("%d", seconds)
	if minutes < 10 {
		minutesString = "0" + minutesString
	}
	if seconds < 10 {
		secondsString = "0" + secondsString
	}

	return fmt.Sprintf("%s:%s", minutesString, secondsString)
}
