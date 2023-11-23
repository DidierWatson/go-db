package invoice

import (
	"github.com/DidierWatson/go-db/pkg/invoiceheader"
	"github.com/DidierWatson/go-db/pkg/invoiceitem"
)

// Invoice Model
type Model struct {
	Header *invoiceheader.Model
	Items  invoiceitem.Models
}

type Storage interface {
	Create(*Model) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
