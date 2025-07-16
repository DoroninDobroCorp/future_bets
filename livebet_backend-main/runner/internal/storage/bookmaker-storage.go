package storage

import (
	"livebets/runner/cmd/config"
	"livebets/runner/internal/entity"
	"sync"
)

type BookmakerStorage struct {
	sync.RWMutex
	bookmakers map[string]entity.Bookmaker
}

func NewBookmakerStorage(cfg []config.BookmakerConfig) *BookmakerStorage {
	bookmaker := make(map[string]entity.Bookmaker)
	for _, val := range cfg {
		bookmaker[val.Name] = entity.Bookmaker{
			Replicas:     val.Replicas,
			ReplicasName: val.ReplicasName,
			Name:         val.Name,
			Path:         val.Path,
			API:          val.API,
		}
	}

	return &BookmakerStorage{
		bookmakers: bookmaker,
	}
}

func (b *BookmakerStorage) ReadAll() map[string]entity.Bookmaker {
	b.RLock()
	defer b.RUnlock()

	result := make(map[string]entity.Bookmaker)
	for i, val := range b.bookmakers {
		result[i] = val
	}

	return result
}

func (b *BookmakerStorage) SetReplicas(command entity.Command) {
	b.Lock()
	defer b.Unlock()

	bookmaker, ok := b.bookmakers[command.Name]
	if !ok {
		return
	}

	bookmaker.Replicas = 0
	if command.Run {
		bookmaker.Replicas = 1
	}

	b.bookmakers[command.Name] = bookmaker

	return
}
