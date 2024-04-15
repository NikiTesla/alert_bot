package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

const (
	subscribeCommand   = "subscribe"
	unsubscribeCommand = "unsubscribe"
	subscriptionCheck  = "check_subscription"
	healthCommand      = "healthz"

	welcomeMessage = "Hello! You can subscribe and get updates here for something interesting 🙃"
	errorReponse   = "Error occured. Try one more time or write @krechetov_n and ask him to fix it 🤕"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := update.Message
		if msg.IsCommand() {
			go b.handleCommand(update.Message)
		} else {
			go b.handleStart(update.Message)
		}
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	logger := log.WithField("chatId", chatId).WithField("command", message.Command())

	switch message.Command() {
	case subscribeCommand:
		response := "You were successfully subscribed to updates"
		if err := b.subscribe(chatId); err != nil {
			logger.WithError(err).Error("cannot subscribe")
			response = errorReponse
		}

		if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, response)); err != nil {
			logger.WithError(err).Error("unable to send response")
		}
	case unsubscribeCommand:
		response := "You were successfully unsubscribed from updates"
		if err := b.unsubscribe(chatId); err != nil {
			logger.WithError(err).Error("cannot unsubscribe")
			response = errorReponse
		}

		if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, response)); err != nil {
			logger.WithError(err).Error("unable to send response")
		}
	case subscriptionCheck:
		response := "You are not subscribed"
		ok, err := b.checkSubscription(chatId)
		if err != nil {
			logger.WithError(err).Error("cannot check subscription")
			response = errorReponse
		}
		if ok {
			response = "You are subscribed"
		}

		if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, response)); err != nil {
			logger.WithError(err).Error("unable to send response")
		}
	case healthCommand:
		healthzResponse := "Everything is great!"
		if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, healthzResponse)); err != nil {
			logger.WithError(err).Error("unable to send response")
		}
	default:
		if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("%s command is unknown", message.Text))); err != nil {
			logger.WithError(err).Error("unable to send response")
		}
	}
}

func (b *Bot) subscribe(chatId int64) error {
	return b.storage.Subscribe(chatId)
}

func (b *Bot) unsubscribe(chatId int64) error {
	return b.storage.Unsubscribe(chatId)
}

func (b *Bot) checkSubscription(chatId int64) (bool, error) {
	subscribersUids, err := b.storage.GetSubscribersUids()
	if err != nil {
		return false, fmt.Errorf("cannot check subscribers, err: %w", err)
	}

	for _, subscriberUid := range subscribersUids {
		if subscriberUid == chatId {
			return true, nil
		}
	}
	return false, nil
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, welcomeMessage)); err != nil {
		log.WithField("chatId", chatId).WithError(err).Error("unable to send start message")
	}
}
