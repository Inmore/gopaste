package storage

import "github.com/inmore/gopaste/internal/model"

type Storage interface {
	Save(p *model.Paste) error
	Load(id string) (*model.Paste, error)
	DeleteExpired() (n int, err error)
	Close() error
}
