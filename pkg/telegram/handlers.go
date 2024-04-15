package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := update.Message
		if msg.IsCommand() {
			go b.handleCommand(update.Message)
			continue
		}
		go b.handleStart(update.Message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	log.WithField("type", "command").Info(message.Text)
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	log.WithField("type", "start").Info(message.Text)
}
