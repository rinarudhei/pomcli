// Package session contains pomodoro session management logic
package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rinarudhei/pomcli/model"
	"github.com/rinarudhei/pomcli/notify"
	"github.com/rinarudhei/pomcli/utils"
)

type SessionState int

func (s SessionState) String() string {
	switch s {
	case NotStarted:
		return "NotStarted"
	case Running:
		return "Running"
	case Paused:
		return "Paused"
	case Finished:
		return "Finished"
	default:
		return ""
	}
}

const (
	NotStarted SessionState = iota
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
	Repo               InMemoryRepository
	SqliteRepo         SqliteRepository
	SessionType        string
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

type InMemoryRepository interface {
	Last() (model.Pomodoro, error)
	Add(p model.Pomodoro) error
	Update(p model.Pomodoro) error
}

type SqliteRepository interface {
	Add(p model.Pomodoro) error
	GetTodayPomodoro() ([]model.Pomodoro, error)
}

func NewSession(repo InMemoryRepository, sqliteRepo SqliteRepository, pomodoroDuration, shortBreakDuration, longBreakDuration time.Duration) *SessionService {
	return &SessionService{
		Repo:               repo,
		SqliteRepo:         sqliteRepo,
		SessionType:        PomodoroSession,
		PomodoroDuration:   pomodoroDuration,
		ShortBreakDuration: shortBreakDuration,
		LongBreakDuration:  longBreakDuration,
	}
}

func (s *SessionService) GetInfos() (string, error) {
	pomodoros, err := s.SqliteRepo.GetTodayPomodoro()
	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	var history string
	for _, p := range pomodoros {
		history += fmt.Sprintf("     🍅 %.0f minutes (%v)\n", (p.PlannedDuration - p.ActualDuration).Minutes(), p.Starttime.Format(time.DateTime))
	}
	return history, nil
}

func (s *SessionService) Decrement(update CallbackFunc) error {
	last, err := s.Repo.Last()
	if err != nil {
		return err
	}

	if last.State == int(Running) || last.State == int(Paused) {
		return nil
	}
	var currentPlannedDuration time.Duration
	switch s.SessionType {
	case PomodoroSession:
		if s.PomodoroDuration <= 5*time.Minute {
			return nil
		}
		s.PomodoroDuration -= 5 * time.Minute
		currentPlannedDuration = s.PomodoroDuration
	case ShortBreakSession:
		if s.ShortBreakDuration <= 5*time.Minute {
			return nil
		}
		s.ShortBreakDuration -= 5 * time.Minute
		currentPlannedDuration = s.ShortBreakDuration
	case LongBreakSession:
		if s.LongBreakDuration <= 5*time.Minute {
			return nil
		}
		s.LongBreakDuration -= 5 * time.Minute
		currentPlannedDuration = s.LongBreakDuration
	}

	update(s.SessionType, durationToDisplayString(currentPlannedDuration), "")
	return nil
}

func (s *SessionService) Increment(update CallbackFunc) error {
	last, err := s.Repo.Last()
	if err != nil {
		return err
	}
	if last.State == int(Running) || last.State == int(Paused) {
		return nil
	}
	var currentPlannedDuration time.Duration
	switch s.SessionType {
	case PomodoroSession:
		s.PomodoroDuration += 5 * time.Minute
		currentPlannedDuration = s.PomodoroDuration
	case ShortBreakSession:
		s.ShortBreakDuration += 5 * time.Minute
		currentPlannedDuration = s.ShortBreakDuration
	case LongBreakSession:
		s.LongBreakDuration += 5 * time.Minute
		currentPlannedDuration = s.LongBreakDuration
	}

	update(s.SessionType, durationToDisplayString(currentPlannedDuration), "")
	return nil
}

func (s *SessionService) Start(ctx context.Context, update CallbackFunc) error {
	last, err := s.Repo.Last()
	if err != nil {
		return err
	}
	if last.State == int(Running) {
		return nil
	}

	var dur time.Duration
	switch s.SessionType {
	case PomodoroSession:
		dur = s.PomodoroDuration
	case ShortBreakSession:
		dur = s.ShortBreakDuration
	case LongBreakSession:
		dur = s.LongBreakDuration
	}

	if last.ID == 0 || last.State == int(Finished) {
		session := model.Pomodoro{
			Type:            s.SessionType,
			State:           int(Running),
			Starttime:       time.Now(),
			PlannedDuration: dur,
		}
		if err := s.Repo.Add(session); err != nil {
			return err
		}
		return s.tick(ctx, update)
	}

	if last.State == int(NotStarted) {
		last.Starttime = time.Now()
	}

	last.State = int(Running)
	if err := s.Repo.Update(last); err != nil {
		return err
	}

	return s.tick(ctx, update)
}

func (s *SessionService) Pause() error {
	p, err := s.Repo.Last()
	if err != nil {
		return err
	}
	p.State = int(Paused)
	if err := s.Repo.Update(p); err != nil {
		return err
	}
	return nil
}

// SwitchState update session type, and return updated session type, duration in string, and title

func (s *SessionService) SwitchState() (model.SwitchStateResponse, error) {
	var durationString string
	var title string

	last, err := s.Repo.Last()
	if err != nil {
		return model.SwitchStateResponse{}, err
	}

	if last.State == int(Running) || last.State == int(Paused) {
		return model.SwitchStateResponse{}, utils.ErrSwitchInActiveState
	}

	switch s.SessionType {
	case PomodoroSession:
		s.SessionType = ShortBreakSession
		durationString = durationToDisplayString(s.ShortBreakDuration)
		title = "Short Break"
	case ShortBreakSession:
		s.SessionType = LongBreakSession
		durationString = durationToDisplayString(s.LongBreakDuration)
		title = "Long Break"
	case LongBreakSession:
		s.SessionType = PomodoroSession
		durationString = durationToDisplayString(s.PomodoroDuration)
		title = "Pomodoro Focus"
	}

	return model.SwitchStateResponse{DurationString: durationString, NextSessionType: s.SessionType, Title: title}, nil
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

type CallbackFunc func(sessionState, timerString, history string)

func (s *SessionService) tick(ctx context.Context, update CallbackFunc) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	p, err := s.Repo.Last()
	if err != nil {
		return err
	}
	expire := time.After(p.PlannedDuration - p.ActualDuration)

	for {
		select {
		case <-t.C:
			p, err := s.Repo.Last()
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
			update(s.SessionType, durationToDisplayString(currentTimer), "")
			if err := s.Repo.Update(p); err != nil {
				return err
			}
		case <-ctx.Done():
			p, err := s.Repo.Last()
			if err != nil {
				return err
			}
			p.State = int(Finished)

			if err := s.SqliteRepo.Add(p); err != nil {
				return err
			}
			return s.Repo.Update(p)
		case <-expire:
			p, err := s.Repo.Last()
			if err != nil {
				return err
			}
			p.State = int(Finished)
			n := notify.New(s.SessionType, "Done", notify.SeverityNormal)
			if err := n.Send(); err != nil {
				fmt.Println(err)
			}
			if err := s.SqliteRepo.Add(p); err != nil {
				return err
			}

			infos, err := s.GetInfos()
			if err != nil {
				return err
			}
			update(s.SessionType, "", infos)

			return s.Repo.Update(p)
		}
	}
}
