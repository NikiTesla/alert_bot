package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	storage Storage
}

func NewBot() *Bot {
	botToken, ok := os.LookupEnv("TELEBOT_API")
	if !ok {
		log.Fatal("bot token is not presented")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.WithError(err).Fatal("can't create Bot API")
	}
	bot.Debug = true

	return &Bot{
		storage: NewMemoryStorage(),
		bot:     bot,
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

func (b *Bot) NotifySubscribers(data []byte) {
	if err := b.notifySubscribers(data); err != nil {
		log.WithError(err).Error("cannot notify subscribers")
	}
}

func (b *Bot) Close() {
	log.Info("bot was closed")
}
