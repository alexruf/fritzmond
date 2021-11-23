package ui

import (
	"github.com/nakabonne/tstorage"
)

type server struct {
	storage tstorage.Storage
}

func New(storage tstorage.Storage) server {
	return server{storage: storage}
}
