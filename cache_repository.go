package main

import "context"

type CacheRepository interface {
	Get(context.Context, ID) (bool, error)
	Put(context.Context, ID, bool) error
}
