package main

import (
	"fmt"
	"log"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhirnov/go-rss-parser/db"
	"github.com/maxzhirnov/go-rss-parser/translate"
)

func publishNewsToChannel(db *db.DB, channelName string, bot *tgbot.BotAPI, translator *translate.Translator, translate bool) error {

	recentNewsItem, err := db.GetMostRecentItem(240)
	if err != nil {
		return err
	}

	if recentNewsItem == nil {
		log.Println("recentNewsItem is nil")
		return nil
	}

	var postTitle string = recentNewsItem.Title
	var postText string = recentNewsItem.Description
	var postSourceLink string = recentNewsItem.URL

	if translate {
		if postTitleReq, err := translator.Translate(recentNewsItem.Title, "ru"); err != nil {
			postTitle = recentNewsItem.Title
		} else {
			postTitle = postTitleReq.Translations[0].Text
		}

		if postTextReq, err := translator.Translate(recentNewsItem.Description, "ru"); err != nil {
			postText = recentNewsItem.Description
		} else {
			postText = postTextReq.Translations[0].Text
		}
	}

	msgText := formatPostMessage(postTitle, postText, postSourceLink)
	message := tgbot.NewMessageToChannel(channelName, msgText)
	message.ParseMode = "Markdown"

	if _, err = bot.Send(message); err != nil {
		return err
	}

	if err := db.UpdatePublishedStatusToTrue(recentNewsItem.GUID); err != nil {
		return err
	}
	return nil
}

func formatPostMessage(postTitle, postText, postSourceLink string) string {
	message := fmt.Sprintf("*%s*\n\n%s\n\n%s", postTitle, postText, postSourceLink)
	if os.Getenv("ENV") == "development" {
		return fmt.Sprintf("[dev]\n\n%s", message)
	}
	return message
}
