package ShortLink

import "time"

type Link struct {
	Id         string    `bson:"_id"`
	OriginLink string    `bson:"origin_link"`
	CreatedAt  time.Time `bson:"createdAt"`
}
