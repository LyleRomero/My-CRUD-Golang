package application

import (
	"My-CRUD-Golang/internal/domain"
	"My-CRUD-Golang/internal/ports"
)

type ItemService struct {
	repository ports.ItemRepository
}

func NewItemService(repo ports.ItemRepository) *ItemService {
	return &ItemService{
		repository: repo,
	}
}

func (s *ItemService) GetAllItems() ([]domain.Item, error) {
	return s.repository.GetAll()
}

func (s *ItemService) GetItemByID(id string) (*domain.Item, error) {
	return s.repository.GetByID(id)
}

func (s *ItemService) CreateItem(item domain.Item) error {
	return s.repository.Create(item)
}

func (s *ItemService) UpdateItem(item domain.Item) error {
	return s.repository.Update(item)
}

func (s *ItemService) DeleteItem(id string) error {
	return s.repository.Delete(id)
}
