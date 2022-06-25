package main

type ItemRepository interface {
	GetAll() ([]Item, error)
	GetItems(ids []ID) ([]Item, error)
}
