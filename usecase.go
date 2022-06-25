package main

import (
	"fmt"
)

type ItemUsecase interface {
	GetAll() ([]Item, error)
	GetItems(ids []ID) ([]Item, error)
}

type itemUsecase struct {
	// config       Config
	repository ItemRepository
	// exRepository ExcludeRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ItemUsecase = (*itemUsecase)(nil)

func NewItemUsecase(repository ItemRepository) ItemUsecase {
	return &itemUsecase{repository: repository}
}

func (s *itemUsecase) GetAll() ([]Item, error) {
	is, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %v", err)
	}

	return is, nil
}

func (s *itemUsecase) GetItems(ids []ID) ([]Item, error) {
	is, err := s.repository.GetItems(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to GetItems: %v", err)
	}

	return is, nil
}
