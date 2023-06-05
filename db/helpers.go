package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) Exists(guid string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := db.Collection.CountDocuments(ctx, bson.M{"guid": guid})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
