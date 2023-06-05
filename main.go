package main

import (
	"log"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	mongodb "github.com/maxzhirnov/go-rss-parser/db"
	"github.com/maxzhirnov/go-rss-parser/translate"
	cron "github.com/robfig/cron/v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	var (
		env                string = os.Getenv("ENV")
		mongodbConnString  string = os.Getenv("MONGO_CONN")
		telegramBotToken   string = os.Getenv("TELEGRAM_TOKEN")
		yTranslateFolderId string = os.Getenv("YTRANSLATE_FOLDER_ID")
		yTranslateToken    string = os.Getenv("YTRANSLATE_TOKEN")
		tgChannelName      string = os.Getenv("TG_CHANNEL_NAME")
		err                error
		db                 *mongodb.DB
		bot                *tgbot.BotAPI
		translator         *translate.Translator
		feeds              []string = []string{
			"https://www.wordstream.com/feed",
			"https://phiture.com/feed/",
			"https://www.revenuecat.com/blog/rss.xml",
			"http://feeds.seroundtable.com/SearchEngineRoundtable1",
		}
	)

	// Initilazing dependencies
	if db, err = mongodb.New(mongodbConnString, "rss-feed", "items"); err != nil {
		log.Fatalf("Error iniyilizing mondodb: %s\n", err)
	}
	if bot, err = tgbot.NewBotAPI(telegramBotToken); err != nil {
		log.Fatalf("Error iniyilizing telegram bot: %s\n", err)
	}
	if translator, err = translate.NewTranslator(yTranslateFolderId, yTranslateToken); err != nil {
		log.Fatalf("Error iniyilizing yTranslate: %s\n", err)
	}

	if env == "production" {
		// Setting up cron jobs
		c := cron.New()
		// Every 30 minutes download new news from list of RSS feeds and store in Mongo if not already exist
		c.AddFunc("*/10 * * * *", func() {
			if err := DownloadNewNews(feeds, *db); err != nil {
				log.Fatal(err)
			}
		})

		//  Every 10 minutes check if there is new news and send it to the channel
		c.AddFunc("* * * * *", func() {
			if err := publishNewsToChannel(db, tgChannelName, bot, translator); err != nil {
				log.Fatal(err)
			}
		})
		c.Start()
	} else {
		if err := publishNewsToChannel(db, tgChannelName, bot, translator); err != nil {
			log.Fatal(err)
		}
	}

	select {}
}
