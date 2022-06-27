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

type RequestIDs struct {
	IDs []ID `json:"ids"`
}

type Items []Item

func (is *Items) Sort() {
	//IDの大きい順にソート
	sort.Slice((*is), func(i, j int) bool { return (*is)[i].ID > (*is)[j].ID })
}

type Price int

// Score is item  score struct with json.
type Score struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Scores []Score

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
