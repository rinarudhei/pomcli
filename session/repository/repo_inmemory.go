// Package repository contains methods to interact with database
package repository

import (
	"sync"

	"github.com/rinarudhei/pomcli/model"
)

type SessionRepository struct {
	sync.RWMutex
	sessions []model.Pomodoro
}

func NewRepository() *SessionRepository {
	return &SessionRepository{}
}

func (s *SessionRepository) Last() (model.Pomodoro, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.sessions) == 0 {
		return model.Pomodoro{}, nil
	}

	last := s.sessions[len(s.sessions)-1]
	return last, nil
}

func (s *SessionRepository) Add(p model.Pomodoro) error {
	s.Lock()
	defer s.Unlock()

	p.ID = int64(len(s.sessions) + 1)
	s.sessions = append(s.sessions, p)

	return nil
}

func (s *SessionRepository) Update(p model.Pomodoro) error {
	s.Lock()
	defer s.Unlock()

	for i, session := range s.sessions {
		if session.ID == p.ID {
			s.sessions[i] = p
			break
		}
	}

	return nil
}
