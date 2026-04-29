package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type RatingResult struct {
	Id          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username    string        `json:"username,omitempty" bson:"username"`
	OldRating   int16         `json:"old_rating,omitempty" bson:"old_rating"`
	NewRating   int16         `json:"new_rating,omitempty" bson:"new_rating"`
	ContestId   bson.ObjectID `json:"contest_id,omitempty" bson:"contest_id"`
	Delta       int16         `json:"delta,omitempty" bson:"delta"`
	Rank        uint16        `json:"rank,omitempty" bson:"rank"`
	Performance int16         `json:"performance,omitempty" bson:"performance"`
}
