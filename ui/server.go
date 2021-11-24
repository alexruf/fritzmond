package ui

import (
	bolt "go.etcd.io/bbolt"
)

type server struct {
	db *bolt.DB
}

func New(db *bolt.DB) server {
	return server{db: db}
}
