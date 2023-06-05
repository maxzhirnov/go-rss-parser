package main

import (
	"log"

	mongodb "github.com/maxzhirnov/go-rss-parser/db"
	parser "github.com/maxzhirnov/go-rss-parser/parser"
)

func DownloadNewNews(feeds []string, db mongodb.DB) error {
	items := make([]mongodb.FeedItem, 0)

	for _, feed := range feeds {
		var feedItems []mongodb.FeedItem
		var err error
		if feedItems, err = parser.ParseFeed(feed); err != nil {
			return err
		}
		items = append(items, feedItems...)
	}

	for _, item := range items {
		if _, err := db.StoreItem(item); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
