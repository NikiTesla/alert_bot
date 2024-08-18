package telegram

import (
	"alert_bot/pkg/storage"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	storage storage.Storage
	logger  *log.Entry
}

func NewBot(logger *log.Entry) *Bot {
	botToken, ok := os.LookupEnv("TELEBOT_API")
	if !ok {
		log.Fatal("bot token is not presented")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.WithError(err).Fatal("can't create Bot API")
	}
	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
	}

	return &Bot{
		storage: storage.New(),
		bot:     bot,
		logger:  log.WithField("type", "telegram-bot"),
	}
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChannel, err := b.bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return fmt.Errorf("cannot create update channel, err: %w", err)
	}

	log.Info("Bot is handling updates")
	b.handleUpdates(updateChannel)

	return nil
}

func (b *Bot) NotifySubscribers(data []byte) error {
	if err := b.notifySubscribers(data); err != nil {
		return fmt.Errorf("cannot notify subscribers, err: %w", err)
	}
	return nil
}

func (b *Bot) NotifySubscribersWithImage(imageFilename string) error {
	if err := b.notifySubscribersWithImage(imageFilename); err != nil {
		return fmt.Errorf("cannot notify subscribers, err: %w", err)
	}
	return nil
}

func (b *Bot) Close() {
	log.Info("bot was closed")
}
