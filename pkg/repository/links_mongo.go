package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/wellWINeo/ShortLink"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type LinksMongo struct {
	collection *mongo.Collection
}

func NewLinksMongo(db *mongo.Database) *LinksMongo {
	return &LinksMongo{collection: db.Collection("links")}
}

func (l *LinksMongo) CreateLink(shortLink, originLink string) (string, error) {
	link := ShortLink.Link{
		Id:         shortLink,
		OriginLink: originLink,
		CreatedAt:  time.Now(),
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	insertResult, err := l.collection.InsertOne(ctx, link)
	if err != nil {
		return "", err
	}
	return insertResult.InsertedID.(string), nil
}

func (l *LinksMongo) UpdateTTL(shortLink string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"_id": shortLink}
	update := bson.M{"$set": bson.M{"createdAt": time.Now()}}
	_, err := l.collection.UpdateOne(ctx, filter, update, nil)
	return err
}

func (l *LinksMongo) GetLink(shortLink string) (string, error) {
	var link ShortLink.Link
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"_id": shortLink}
	err := l.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return link.OriginLink, nil
}
