package storage

import (
	"context"

	"github.com/go-faker/faker/v4"

	"github.com/mrksmt/test-task-7/internal/server/service"
)

type fakeStorage struct{}

var _ service.Storage = (*fakeStorage)(nil)

func NewFakeStorage() *fakeStorage {
	s := &fakeStorage{}
	return s
}

func (s *fakeStorage) GetSentence(
	ctx context.Context,
) (string, error) {
	return faker.Sentence(), nil
}
