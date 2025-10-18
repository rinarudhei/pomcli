// Package repository contains methods to interact with database
package repository

import (
	"github.com/rinarudhei/pomcli/model"
)

type SessionRepository struct{}

func NewRepository() *SessionRepository {
	return &SessionRepository{}
}

func (s *SessionRepository) Last() (model.Pomodoro, error) {
	return model.Pomodoro{}, nil
}

func (s *SessionRepository) Add() error {
	return nil
}

func (s *SessionRepository) Update() error {
	return nil
}
