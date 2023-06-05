package parser

import (
	"log"
	"strings"
	"time"

	"github.com/maxzhirnov/go-rss-parser/db"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

type WordStreamParser struct {
	FeedURL  string
	Category string
	SiteName string
}

func (p WordStreamParser) Parse() ([]db.FeedItem, error) {
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

		description, err := extractFirstPContent(i.Description)
		if err != nil {
			log.Println(err)
			continue
		}

		content, err := parseHTML(i.Content)
		if err != nil {
			log.Println(err)
			content = i.Content
		}

		link := i.Link

		published, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", i.Published)
		if err != nil {
			log.Println(err)
			published = time.Now()
		}

		authorName := i.Author.Name
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

func extractFirstPContent(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "p" {
			return n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result := f(c)
			if result != "" {
				return result
			}
		}
		return ""
	}
	return f(doc), nil
}
