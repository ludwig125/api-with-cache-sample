package main

import (
	"context"
	"fmt"
)

type RankingUsecase interface {
	Ranking(Items) (Items, error)
}

type rankingUsecase struct {
	repository ScoreRepository
	cache      CacheRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ RankingUsecase = (*rankingUsecase)(nil)

func NewRankingUsecase(repository ScoreRepository, cache CacheRepository) RankingUsecase {
	return &rankingUsecase{repository: repository, cache: cache}
}

func (u *rankingUsecase) Ranking(is Items) (Items, error) {
	ss, err := u.getScores(extractIDsFromItems(is))
	if err != nil {
		return nil, fmt.Errorf("failed to getScores: %v", err)
	}

	// 事前にScoreを大きい順にsort
	ss.Sort()
	fmt.Println("ss", ss)

	is = is.SortByScore(ss)
	fmt.Println("is", is)

	for _, i := range ss {
		if i.Score == 0 {
			fmt.Println("put Item", i)
			_ = u.cache.Put(context.Background(), i.ID, true)
		}
	}

	return is, err
}

func extractIDsFromItems(is Items) []ID {
	ids := make([]ID, len(is))
	for i, ite := range is {
		ids[i] = ID(ite.ID)
	}
	return ids
}

func (u *rankingUsecase) getScores(ids []ID) (Scores, error) {
	ss, err := u.repository.GetScores(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to GetScores: %v", err)
	}
	return ss, nil
}
