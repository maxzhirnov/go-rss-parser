package parser

import (
	"fmt"

	"github.com/maxzhirnov/go-rss-parser/db"
)

// type FeedItem struct {
// 	Title       string
// 	Description string
// 	Content     string
// 	URL         string
// 	PubDate     time.Time
// 	Author      string
// 	GUID        string
// 	Website     string
// 	Category    string
// }

type RSSParser interface {
	Parse() ([]db.FeedItem, error)
}

// func (f *db.FeedItem) String() string {
// 	return fmt.Sprintf("Title: %s\nDescription: %s\nContent: %s\nURL: %s\nPubDate: %s\nAuthor: %s\nGUID: %s\nWebsite: %s\nCategory: %s\n", f.Title, f.Description, f.Content, f.URL, f.PubDate, f.Author, f.GUID, f.Website, f.Category)
// }

func ParseFeed(feedURL string) ([]db.FeedItem, error) {
	var parser RSSParser

	switch feedURL {
	case "https://www.wordstream.com/feed":
		parser = &WordStreamParser{
			FeedURL:  "https://www.wordstream.com/feed",
			Category: "Advertising",
			SiteName: "www.wordstream.com",
		}
	case "https://phiture.com/feed/":
		parser = &WordStreamParser{
			FeedURL:  "https://phiture.com/feed/",
			Category: "Advertising",
			SiteName: "www.phiture.com",
		}
	case "https://www.revenuecat.com/blog/rss.xml":
		parser = &RevenueCatParser{
			FeedURL:  "https://www.revenuecat.com/blog/rss.xml",
			Category: "Advertising",
			SiteName: "www.revenuecat.com",
		}
	case "http://feeds.seroundtable.com/SearchEngineRoundtable1":
		parser = &SearchRountTableParser{
			FeedURL:  "http://feeds.seroundtable.com/SearchEngineRoundtable1",
			Category: "Advertising",
			SiteName: "www.seroundtable.com",
		}
	default:
		return nil, fmt.Errorf("Unknown feed URL: %s", feedURL)
	}

	return parser.Parse()
}
