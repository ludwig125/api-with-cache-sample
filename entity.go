package main

import (
	"fmt"
	"sort"
	"strconv"
)

// Item is item struct with json.
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ID int

// type RequestIDs struct {
// 	IDs []ID `json:"ids"`
// }

type Items []Item

func (is Items) RemoveItems(removeIDs []ID) Items {
	var items Items
	for _, i := range is {
		if !ID(i.ID).contains(removeIDs) {
			items = append(items, i)
		}
	}
	return items
}

func (id ID) contains(removeIDs []ID) bool {
	for _, r := range removeIDs {
		if id == r {
			return true
		}
	}
	return false
}

func (is *Items) Sort() {
	//IDの大きい順にソート
	sort.Slice((*is), func(i, j int) bool { return (*is)[i].ID > (*is)[j].ID })
}

func (is Items) Limit(limit int) Items {
	newIs := make(Items, len(is))
	copy(newIs, is)
	if limit >= len(newIs) {
		return newIs
	}
	return newIs[:limit]
}

// SortByScore sort Items by score
func (is Items) SortByScore(scores Scores) Items {
	newIs := make(Items, len(scores)) // scoreのかずだけ確保
	copy(newIs, is)

	itemMap := make(map[ID]Item, len(is))
	for _, item := range is {
		itemMap[ID(item.ID)] = item
	}

	for i, s := range scores {
		if item, ok := itemMap[ID(s.ID)]; ok {
			newIs[i] = item
		}
	}

	return newIs
}

type Price int

// Score is item  score struct with json.
type Score struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Scores []Score

func (ss *Scores) Sort() {
	//IDの大きい順にソート
	sort.Slice((*ss), func(i, j int) bool { return (*ss)[i].Score > (*ss)[j].Score })
}

const (
	Invalid = iota
	PriceEqualTo
	PriceLessThanAndEqualTo
	PriceGreaterThanAndEqualTo
)

type SearchCondition struct {
	// ID         ID
	// Name       string
	Price      Price
	Expression string
}

func NewSearchCondition(price, expr string) (*SearchCondition, error) {
	p, err := strconv.Atoi(price)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to int: %s", price)
	}
	return &SearchCondition{
		Price:      Price(p),
		Expression: expr,
	}, nil
}

func (c SearchCondition) CheckCond() int {
	switch {
	case c.Price != 0:
		if c.Expression == "lessthan" {
			return PriceLessThanAndEqualTo
		}
		if c.Expression == "greaterthan" {
			return PriceGreaterThanAndEqualTo
		}
		return PriceEqualTo
	default:
		return Invalid
	}
}
