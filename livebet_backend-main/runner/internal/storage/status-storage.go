package storage

import (
	"livebets/runner/internal/entity"
	"sync"
)

type StatusStorage struct {
	sync.RWMutex
	statuses map[string]entity.StatusBookmaker
}

func NewStatusStorage() *StatusStorage {
	return &StatusStorage{
		statuses: make(map[string]entity.StatusBookmaker),
		RWMutex:  sync.RWMutex{},
	}
}

func (s *StatusStorage) ReadAll() map[string]entity.StatusBookmaker {
	s.RLock()
	defer s.RUnlock()

	result := make(map[string]entity.StatusBookmaker)
	for i, val := range s.statuses {
		result[i] = val
	}

	return result
}

func (s *StatusStorage) SetStatus(status entity.StatusBookmaker) {
	s.Lock()
	defer s.Unlock()

	s.statuses[status.Name] = status
}
