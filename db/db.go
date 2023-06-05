package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

func New(uri string, dbName string, collectionName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &DB{
		Client:     client,
		Database:   database,
		Collection: collection,
	}, nil
}

func (db *DB) StoreItem(item FeedItem) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := db.Exists(item.GUID)
	if err != nil {
		log.Printf("Could not check if item exists: %v", err)
		return nil, err
	}
	if exists {
		log.Printf("Item with GUID %s already exists", item.GUID)
		return nil, nil
	}

	result, err := db.Collection.InsertOne(ctx, bson.M{
		"title":              item.Title,
		"description":        item.Description,
		"content":            item.Content,
		"url":                item.URL,
		"pubDate":            item.PubDate,
		"author":             item.Author,
		"guid":               item.GUID,
		"website":            item.Website,
		"category":           item.Category,
		"publishedToChannel": false,
	})

	if err != nil {
		log.Printf("Could not store item: %v", err)
		return nil, err
	}

	return result, nil
}

func (db *DB) GetUnpublishedItems() ([]FeedItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.Collection.Find(ctx, bson.M{"publishedToChannel": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []FeedItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (db *DB) GetItemsForLastNHours(hours int) ([]FeedItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate the time 'n' hours ago
	nHoursAgo := time.Now().Add(-time.Duration(hours) * time.Hour)

	filter := bson.M{
		"pubDate": bson.M{
			"$gte": nHoursAgo,
		},
		"publishedToChannel": false,
	}

	cursor, err := db.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []FeedItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (db *DB) UpdatePublishedStatusToTrue(guid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"guid": guid}
	update := bson.M{
		"$set": bson.M{
			"publishedToChannel": true,
		},
	}

	_, err := db.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (db *DB) GetMostRecentItem(hours int) (*FeedItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate the time 'n' hours ago
	nHoursAgo := time.Now().Add(-time.Duration(hours) * time.Hour)

	filter := bson.M{
		"pubDate": bson.M{
			"$gte": nHoursAgo,
		},
		"publishedToChannel": false,
	}

	// Define the sorting option: in this case, descending by "pubDate".
	opts := options.FindOne().SetSort(bson.D{{Key: "pubDate", Value: -1}})

	var item FeedItem
	err := db.Collection.FindOne(ctx, filter, opts).Decode(&item)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no document found, return nil
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}
