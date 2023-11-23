package invoiceheader

import (
	"database/sql"
	"time"
)

type Model struct {
	ID        uint
	Client    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, *Model) error
}

// invoiceheader's service
type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

// used to migrate product
func (s *Service) Migrate() error {
	return s.storage.Migrate()
}
