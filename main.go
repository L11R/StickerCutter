package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
	"golang.org/x/net/proxy"
	"net/http"
	"github.com/spf13/viper"
	"fmt"
)

var log = logrus.New()

var (
	bot *tgbotapi.BotAPI
)

func main() {
	log.Formatter = new(logrus.TextFormatter)
	log.Info("Sticker Cutter started!")

	var err error

	// Load configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN env variable not specified!")
	}

	// Socks for Russia
	socksDialer, err := proxy.SOCKS5("tcp", "socks.druble.ru:1080", &proxy.Auth{
		User: "free",
		Password: "forfriends",
	}, proxy.Direct)
	if err != nil {
		log.Fatal("Proxy error!")
	}
	tr := http.Transport{
		Dial: socksDialer.Dial,
	}

	bot, err = tgbotapi.NewBotAPIWithClient(token, &http.Client{
		Transport: &tr,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("authorized on account @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		// Handle only messages
		if update.Message == nil {
			continue
		}

		if update.Message.Photo != nil {
			HandlePhoto(update)
		}

		if update.Message.VideoNote != nil {
			HandleVideoNote(update)
		}

		// commands
		switch update.Message.Command() {
		case "start":
			go StartCommand(update)
		}
	}
}
