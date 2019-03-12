package reader

import (
	root ".."
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFileByID(fid string) (*root.FileMeta, error) {
	var fm root.FileMeta

	id, herr := primitive.ObjectIDFromHex(fid)
	if herr != nil {
		return nil, herr
	}
	filter := bson.D{{"_id", id}}
	err := metaService.Collection.FindOne(metaService.Context, filter).Decode(&fm)
	if err != nil {
		return nil, err
	}
	return &fm, err
}
