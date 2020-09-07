package repositories

import (
	"errors"
	"sync"
)

type (
	// Repository struct
	repository struct {
		lock     sync.RWMutex
		Database []*Atheletes
	}
	// IRepository interface
	IRepository interface {
		SaveAthelete(b *Atheletes) error
		GetAtheletes() []*Atheletes
	}
	Atheletes struct {
		Code    string `json:"code"`
		Number  string `json:"number"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}
)

var (
	repo IRepository
	once sync.Once
)

// Singleton instance of Athlete accounts repository
func GetInstance() IRepository {
	once.Do(func() {

		repo = &repository{
			Database: []*Atheletes{},
		}
	})
	return repo
}

func (r *repository) SaveAthelete(b *Atheletes) error {
	if b == nil {
		return errors.New("athlete is nil")
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Database = append(r.Database, b)
	return nil
}

func (r *repository) GetAtheletes() []*Atheletes {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.Database
}
