package db

import (
	"fmt"
	"time"
)

type FeedItem struct {
	Title              string    `json:"title" bson:"title"`
	Description        string    `json:"description" bson:"description"`
	Content            string    `json:"content" bson:"content"`
	URL                string    `json:"url" bson:"url"`
	PubDate            time.Time `json:"pubDate" bson:"pubDate"`
	Author             string    `json:"author" bson:"author"`
	GUID               string    `json:"guid" bson:"guid"`
	Website            string    `json:"website" bson:"website"`
	Category           string    `json:"category" bson:"category"`
	PublishedToChannel bool      `json:"publishedToChannel" bson:"publishedToChannel"`
}

func (f *FeedItem) String() string {
	return fmt.Sprintf("Title: %s\nDescription: %s\nContent: %s\nURL: %s\nPubDate: %s\nAuthor: %s\nGUID: %s\nWebsite: %s\nCategory: %s\n", f.Title, f.Description, f.Content, f.URL, f.PubDate, f.Author, f.GUID, f.Website, f.Category)
}
