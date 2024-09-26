package db

import (
	"My-CRUD-Golang/internal/domain"
	"My-CRUD-Golang/internal/ports"
	"errors"
)

type MemoryRepository struct {
	items []domain.Item
}

func NewMemoryRepository() ports.ItemRepository {
	return &MemoryRepository{
		items: []domain.Item{},
	}
}

func (r *MemoryRepository) GetAll() ([]domain.Item, error) {
	return r.items, nil
}

func (r *MemoryRepository) GetByID(id string) (*domain.Item, error) {
	for _, item := range r.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, errors.New("Item not found")
}

func (r *MemoryRepository) Create(item domain.Item) error {
	r.items = append(r.items, item)
	return nil
}

func (r *MemoryRepository) Update(updatedItem domain.Item) error {
	for i, item := range r.items {
		if item.ID == updatedItem.ID {
			r.items[i] = updatedItem
			return nil
		}
	}
	return errors.New("Item not found")
}

func (r *MemoryRepository) Delete(id string) error {
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}
