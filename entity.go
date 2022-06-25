package main

// Item is item struct with json.
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ID int

type IDs struct {
	IDs []ID `json:"ids"`
}
