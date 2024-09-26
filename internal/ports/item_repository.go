package ports

import "My-CRUD-Golang/internal/domain"

type ItemRepository interface {
	GetAll() ([]domain.Item, error)
	GetByID(id string) (*domain.Item, error)
	Create(item domain.Item) error
	Update(item domain.Item) error
	Delete(id string) error
}
