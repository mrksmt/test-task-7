package service

import (
	"context"
)

type Service interface {
	GetSentence(ctx context.Context) (string, error)
}

type Storage interface {
	GetSentence(ctx context.Context) (string, error)
}

type service struct {
	storage Storage
}

func NewService(
	storage Storage,
) Service {
	s := &service{storage: storage}
	return s
}

func (s *service) GetSentence(
	ctx context.Context,
) (string, error) {
	return s.storage.GetSentence(ctx)
}
