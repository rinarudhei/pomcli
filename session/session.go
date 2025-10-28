// Package session contains pomodoro session management logic
package session

import (
	"context"
	"fmt"
	"time"

	"github.com/rinarudhei/pomcli/model"
)

type sessionState int

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
	Last() (*model.Pomodoro, error)
	Add(p *model.Pomodoro) error
	Update(p *model.Pomodoro) error
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

func (s *SessionService) Start(ctx context.Context, update CallbackFunc) error {
	last, err := s.repo.Last()
	if err != nil {
		return err
	}
	if last.State == int(Running) {
		return nil
	}

	var dur time.Duration
	switch s.sessionType {
	case PomodoroSession:
		dur = s.pomodoroDuration
	case ShortBreakSession:
		dur = s.shortBreakDuration
	case LongBreakSession:
		dur = s.longBreakDuration
	}

	if last.ID == 0 || last.State == int(Finished) {
		session := &model.Pomodoro{
			Type:            s.sessionType,
			State:           int(Running),
			Starttime:       time.Now(),
			PlannedDuration: dur,
		}
		last.Starttime = time.Now()
		s.repo.Add(session)
		return s.tick(ctx, update)
	}

	if last.State == int(NotStarted) {
		last.Starttime = time.Now()
	}

	last.State = int(Running)
	if err := s.repo.Update(last); err != nil {
		return err
	}

	return s.tick(ctx, update)
}

func (s *SessionService) Pause() error {
	p, err := s.repo.Last()
	if err != nil {
		return err
	}
	p.State = int(Paused)
	if err := s.repo.Update(p); err != nil {
		return err
	}
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

type CallbackFunc func(sessionState string, timerString string)

func (s *SessionService) tick(ctx context.Context, update CallbackFunc) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	p, err := s.repo.Last()
	if err != nil {
		return err
	}
	expire := time.After(p.PlannedDuration - p.ActualDuration)
	for {
		select {
		case <-t.C:
			p, err := s.repo.Last()
			if err != nil {
				return err
			}
			if p.State == int(Paused) {
				return nil
			}
			if p.State == int(Running) {
				p.ActualDuration += time.Second
			}
			currentTimer := p.PlannedDuration - p.ActualDuration
			update(s.sessionType, durationToDisplayString(currentTimer))
			if err := s.repo.Update(p); err != nil {
				return err
			}
		case <-ctx.Done():
			p, err := s.repo.Last()
			if err != nil {
				return err
			}
			p.State = int(Finished)
			return s.repo.Update(p)
		case <-expire:
			p, err := s.repo.Last()
			if err != nil {
				return err
			}
			p.State = int(Finished)
			return s.repo.Update(p)
		}
	}
}
