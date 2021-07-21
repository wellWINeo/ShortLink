package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Links interface {
	CreateLink(shortLink, originLink string) (string, error)
	GetLink(shortLink string) (string, error)
	UpdateTTL(shortLink string) error
}

type Repository struct {
	Links
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{Links: NewLinksMongo(db)}
}
