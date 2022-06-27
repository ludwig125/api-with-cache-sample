package main

import (
	"fmt"
)

type ItemUsecase interface {
	GetAll() (Items, error)
	Search(SearchCondition) (Items, error)
	// GetScores(Items) (Items, error)
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

func (s *itemUsecase) GetAll() (Items, error) {
	is, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %v", err)
	}

	return is, nil
}

func (s *itemUsecase) Search(c SearchCondition) (Items, error) {
	is, err := s.searchItems(c)
	if err != nil {
		return nil, fmt.Errorf("failed to searchItems: %v", err)
	}

	is.Sort()

	// ここでlimit処理入れる

	ss, err := s.getScores(extractIDsFromItems(is))
	if err != nil {
		return nil, fmt.Errorf("failed to getScores: %v", err)
	}
	fmt.Println("ss", ss)

	return is, err
}

func (s *itemUsecase) searchItems(c SearchCondition) (Items, error) {
	var is []Item
	var err error

	switch c.CheckCond() {
	case PriceEqualTo:
		is, err = s.repository.SearchByPriceEqualTo(c.Price)
	case PriceLessThanAndEqualTo:
		is, err = s.repository.SearchByPriceLessThanAndEqualTo(c.Price)
	case PriceGreaterThanAndEqualTo:
		is, err = s.repository.SearchByPriceGreaterThanAndEqualTo(c.Price)
	default:
		return nil, fmt.Errorf("invalid condition: %#v", c)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search: %v", err)
	}
	if len(is) == 0 {
		return nil, fmt.Errorf("did not meet the conditions %#v", c)
	}
	return is, nil
}

func extractIDsFromItems(is Items) []ID {
	ids := make([]ID, len(is))
	for i, ite := range is {
		ids[i] = ID(ite.ID)
	}
	return ids
}

func (s *itemUsecase) getScores(ids []ID) (Scores, error) {
	ss, err := s.repository.GetScores(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to GetScores: %v", err)
	}
	return ss, nil
}
