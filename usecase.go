package main

import (
	"fmt"
)

type SearchUsecase interface {
	GetAll() (Items, error)
	Search(SearchCondition) (Items, error)
}

type searchUsecase struct {
	itemUsecase    ItemUsecase
	rankingUsecase RankingUsecase
	cache          CacheRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ SearchUsecase = (*searchUsecase)(nil)

func NewSearchUsecase(itemUsecase ItemUsecase, rankingUsecase RankingUsecase, cache CacheRepository) SearchUsecase {
	return &searchUsecase{itemUsecase: itemUsecase, rankingUsecase: rankingUsecase, cache: cache}
}

func (u *searchUsecase) GetAll() (Items, error) {
	is, err := u.itemUsecase.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %v", err)
	}

	return is, nil
}

func (u *searchUsecase) Search(c SearchCondition) (Items, error) {
	is, err := u.itemUsecase.GetItems(c)
	if err != nil {
		return nil, fmt.Errorf("failed to GetItems: %v", err)
	}
	is, err = u.rankingUsecase.Ranking(is)
	if err != nil {
		return nil, fmt.Errorf("failed to GetScores: %v", err)
	}
	return is, err
}
