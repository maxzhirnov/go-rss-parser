package parser

import (
	"log"
	"time"

	"github.com/maxzhirnov/go-rss-parser/db"
	"github.com/mmcdole/gofeed"
)

type SearchRountTableParser struct {
	FeedURL  string
	Category string
	SiteName string
}

func (p SearchRountTableParser) Parse() ([]db.FeedItem, error) {
	var result []db.FeedItem

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(p.FeedURL)
	if err != nil {
		return nil, err
	}

	for _, i := range feed.Items {

		//PROCESSING FIELDS
		// title := i.Title
		// description := i.Description
		// content := i.Content
		// link := i.Link
		// published := i.Published
		// authorName := i.Author.Name
		// guid := i.GUID
		title := i.Title

		description, err := parseHTML(i.Description)
		if err != nil {
			log.Println(err)
			description = i.Content
		}

		content, err := parseHTML(i.Content)
		if err != nil {
			log.Println(err)
			content = i.Content
		}

		link := i.Link

		const layout = "Mon, 02 Jan 2006 15:04:05 -0700"
		published, err := time.Parse(layout, i.Published)
		if err != nil {
			log.Fatal(err)
		}

		authorName := ""

		if i.Author != nil {
			authorName = i.Author.Name
		}

		guid := i.GUID
		//PROCESSING FIELDS ENDS

		item := db.FeedItem{
			Title:       title,
			Description: description,
			Content:     content,
			URL:         link,
			PubDate:     published,
			Author:      authorName,
			GUID:        guid,
			Website:     p.SiteName,
			Category:    p.Category,
		}
		result = append(result, item)
	}

	return result, nil
}
