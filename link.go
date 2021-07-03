package ShortLink

type Link struct {
	//Id bson.ObjectId `bson:"_id"`
	Id         string `bson:"_id"`
	OriginLink string `bson:"origin_link"`
}
