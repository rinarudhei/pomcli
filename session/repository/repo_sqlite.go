// Package repository contains methods to interact with database
package repository

import (
	"sync"

	"github.com/rinarudhei/pomcli/model"
)

type sessionRepository struct {
	sync.RWMutex
	sessions []model.Pomodoro
}

func NewRepository() *sessionRepository {
	return &sessionRepository{}
}

func (s *sessionRepository) Last() (*model.Pomodoro, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.sessions) == 0 {
		return &model.Pomodoro{}, nil
	}

	last := s.sessions[len(s.sessions)-1]
	return &last, nil
}

func (s *sessionRepository) Add(p *model.Pomodoro) error {
	s.Lock()
	defer s.Unlock()

	if len(s.sessions) == 0 {
		p.ID = int64(1)
	} else {
		p.ID = int64(len(s.sessions))
	}
	s.sessions = append(s.sessions, *p)
	return nil
}

func (s *sessionRepository) Update(p *model.Pomodoro) error {
	s.Lock()
	defer s.Unlock()

	for i, ses := range s.sessions {
		if ses.ID == p.ID {
			s.sessions[i] = *p
			break
		}
	}
	return nil
}
