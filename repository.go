package main

type ItemRepository interface {
	GetAll() (Items, error)
	SearchByPriceEqualTo(price Price) (Items, error)
	SearchByPriceLessThanAndEqualTo(price Price) (Items, error)
	SearchByPriceGreaterThanAndEqualTo(price Price) (Items, error)
}

type ScoreRepository interface {
	GetScores(ids []ID) (Scores, error)
}
