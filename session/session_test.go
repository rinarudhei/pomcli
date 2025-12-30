package session_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rinarudhei/pomcli/model"
	"github.com/rinarudhei/pomcli/session"
	"github.com/rinarudhei/pomcli/session/repository"
)

func TestNewSession(t *testing.T) {
	testCases := []struct {
		name           string
		pomDuration    time.Duration
		sbreakDuration time.Duration
		lbreakDuration time.Duration
	}{
		{
			name:           "initSession",
			pomDuration:    1 * time.Second,
			sbreakDuration: 1 * time.Second,
			lbreakDuration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.NewRepository()
			sqliteRepo, cleanup := generateSqliteSession(t)
			defer cleanup()
			s := session.NewSession(repo, sqliteRepo, tc.pomDuration, tc.sbreakDuration, tc.lbreakDuration)

			if s.SessionType != session.PomodoroSession {
				t.Errorf("expect session type %s, got %s", session.PomodoroSession, s.SessionType)
			}

			if s.PomodoroDuration != tc.pomDuration {
				t.Errorf("expect pomodoro duration %v, got %v", tc.pomDuration, s.PomodoroDuration)
			}
			if s.ShortBreakDuration != tc.sbreakDuration {
				t.Errorf("expect short break duration %v, got %v", tc.sbreakDuration, s.ShortBreakDuration)
			}
			if s.LongBreakDuration != tc.lbreakDuration {
				t.Errorf("expect long break duration %v, got %v", tc.lbreakDuration, s.LongBreakDuration)
			}
		})
	}
}

func generateSqliteSession(t *testing.T) (*repository.DbRepo, func()) {
	t.Helper()

	sqliteRepo, err := repository.NewSQLiteRepo(":memory:")
	if err != nil {
		t.Fatal("error generating sqlite repo")
	}
	return sqliteRepo, func() {
		sqliteRepo.Close()
	}
}

func generateRepoWithFinishedSession(t *testing.T) *repository.SessionRepository {
	t.Helper()

	r := repository.NewRepository()
	p := model.Pomodoro{
		Type:            session.PomodoroSession,
		State:           int(session.Finished),
		Starttime:       time.Now(),
		PlannedDuration: 2 * time.Second,
		ActualDuration:  2 * time.Second,
	}

	if err := r.Add(p); err != nil {
		t.Error("error generating repository with previously finished session")
	}

	return r
}

func TestStartSession(t *testing.T) {
	testCases := []struct {
		name       string
		mockedRepo session.InMemoryRepository
		expSession model.Pomodoro
		expErr     error
	}{
		// {
		// 	name:       "FirstTimeStartSession",
		// 	mockedRepo: repository.NewRepository(),
		// 	expSession: model.Pomodoro{
		// 		ID:              int64(1),
		// 		Type:            session.PomodoroSession,
		// 		State:           int(session.Finished),
		// 		PlannedDuration: 3 * time.Second,
		// 		ActualDuration:  3 * time.Second,
		// 	},
		// },
		{
			name:       "StartSessionWithPreviouslyFinishedSession",
			mockedRepo: generateRepoWithFinishedSession(t),
			expSession: model.Pomodoro{
				ID:              int64(2),
				Type:            session.PomodoroSession,
				State:           int(session.Finished),
				PlannedDuration: 2 * time.Second,
				ActualDuration:  2 * time.Second,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqliteRepo, cleanup := generateSqliteSession(t)
			defer cleanup()
			s := session.NewSession(tc.mockedRepo, sqliteRepo, tc.expSession.PlannedDuration, tc.expSession.PlannedDuration, tc.expSession.PlannedDuration)
			update := func(sessionState, timerString, history, activities, summary string) {}

			err := s.Start(context.Background(), update)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			last, err := tc.mockedRepo.Last()
			if err != nil {
				t.Errorf("error asserting last session state")
			}

			if last.ID != tc.expSession.ID {
				t.Errorf("expect ID %d, got %d", tc.expSession.ID, last.ID)
			}

			if last.PlannedDuration != tc.expSession.PlannedDuration {
				t.Errorf("expect duration %v, got %v", tc.expSession.PlannedDuration, last.PlannedDuration)
			}

			if last.Type != tc.expSession.Type {
				t.Errorf("expect session type %s, got %s", tc.expSession.Type, last.Type)
			}

			if last.State != tc.expSession.State {
				t.Errorf("expect last session state to be %s, got %s", fmt.Sprint(tc.expSession.State), fmt.Sprint(last.State))
			}
		})
	}
}
