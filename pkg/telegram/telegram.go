package telegram

import (
	"alert_bot/pkg/storage"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		logger.Fatal("bot token is not presented")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.WithError(err).Fatal("can't create Bot API")
	}
	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
	}

	storage, err := storage.New()
	if err != nil {
		logger.WithError(err).Fatal("failed to create storage")
	}

	return &Bot{
		storage: storage,
		bot:     bot,
		logger:  logger.WithField("type", "telegram-bot"),
	}
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	log.Info("Bot is handling updates")
	b.handleUpdates(b.bot.GetUpdatesChan(updateConfig))

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
