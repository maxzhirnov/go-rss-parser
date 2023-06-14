package main

import (
	"log"

	mongodb "github.com/maxzhirnov/go-rss-parser/db"
	parser "github.com/maxzhirnov/go-rss-parser/parser"
)

func DownloadNewNews(feeds []string, db mongodb.DB) error {
	log.Println("Downloading new news...")
	items := make([]mongodb.FeedItem, 0)

	for _, feed := range feeds {
		log.Printf("Parsing feed %s\n", feed)
		var feedItems []mongodb.FeedItem
		var err error
		if feedItems, err = parser.ParseFeed(feed); err != nil {
			return err
		}
		items = append(items, feedItems...)
	}

	parsedItems := 0
	insertedItems := 0

	for _, item := range items {
		parsedItems += 1
		if d, err := db.StoreItem(item); err != nil {
			log.Println(err)
			continue
		} else {
			if d != nil {
				insertedItems += 1
			}
		}
	}
	log.Printf("Parsed %d items\n", parsedItems)
	log.Printf("Inserted %d items\n", insertedItems)
	return nil
}
