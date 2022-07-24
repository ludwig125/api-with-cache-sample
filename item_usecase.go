package main

import (
	"context"
	"fmt"
)

type ItemUsecase interface {
	GetAll() (Items, error)
	GetItems(SearchCondition) (Items, error)
}

type itemUsecase struct {
	repository ItemRepository
	cache      CacheRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ItemUsecase = (*itemUsecase)(nil)

func NewItemUsecase(repository ItemRepository, cache CacheRepository) ItemUsecase {
	return &itemUsecase{repository: repository, cache: cache}
}

func (u *itemUsecase) GetAll() (Items, error) {
	is, err := u.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %v", err)
	}

	return is, nil
}

func (u *itemUsecase) GetItems(c SearchCondition) (Items, error) {
	is, err := u.searchItems(c)
	if err != nil {
		return nil, fmt.Errorf("failed to searchItems: %v", err)
	}

	is.Sort()

	var cachedItemIDs []ID
	for _, i := range is {
		if ok, _ := u.cache.Get(context.Background(), i.ID); ok {
			cachedItemIDs = append(cachedItemIDs, ID(i.ID))
		}
	}
	if len(cachedItemIDs) != 0 {
		is = is.RemoveItems(cachedItemIDs)
	}

	limit := 10
	is = is.Limit(limit)

	return is, err
}

func (u *itemUsecase) searchItems(c SearchCondition) (Items, error) {
	var is []Item
	var err error

	switch c.CheckCond() {
	case PriceEqualTo:
		is, err = u.repository.SearchByPriceEqualTo(c.Price)
	case PriceLessThanAndEqualTo:
		is, err = u.repository.SearchByPriceLessThanAndEqualTo(c.Price)
	case PriceGreaterThanAndEqualTo:
		is, err = u.repository.SearchByPriceGreaterThanAndEqualTo(c.Price)
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
