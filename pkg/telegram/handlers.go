package telegram

import (
	"alert_bot/pkg/domain"
	"alert_bot/pkg/errs"
	"alert_bot/pkg/model"
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const (
	startCommand             = "start"
	subscribeCommand         = "subscribe"
	unsubscribeCommand       = "unsubscribe"
	sendDataCommand          = "send_data"
	subscriptionCheckCommand = "check_subscription"
	healthCommand            = "healthz"

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
			go b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	logger := b.logger.WithField("chatId", chatId).WithField("command", message.Command())

	switch message.Command() {
	case startCommand:
		b.handleStart(message)

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

	case subscriptionCheckCommand:
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

	case sendDataCommand:
		response := "Enter your data to save in format..."
		if err := b.storage.SetStatus(message.Chat.ID, model.SendingData); err != nil {
			logger.WithError(err).Error("cannot set status")
			response = "Sending data failed. Status model failed"
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

func (b *Bot) handleStart(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	if _, err := b.bot.Send(tgbotapi.NewMessage(chatId, welcomeMessage)); err != nil {
		b.logger.WithField("chatId", chatId).WithError(err).Error("unable to send start message")
	}
	if err := b.setCommands(context.TODO(), message.Chat.ID, domain.UnknownRole); err != nil {
		b.logger.WithField("chatId", chatId).WithError(err).Error("unable to set commands")
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

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	status, err := b.storage.GetStatus(message.Chat.ID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		b.logger.WithError(err).Error("failed to get status")
		return
	}

	switch status {
	case model.SendingData:
		b.logger.Debug("processing sending data status")
		b.sendData(message)
	default:
		b.handleStart(message)
	}
}

func (b *Bot) sendData(message *tgbotapi.Message) {
	response := "Data was successfully sent"
	defer func() {
		if _, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, response)); err != nil {
			b.logger.WithField("chatId", message.Chat.ID).WithError(err).Error("unable to send start message")
		}
	}()
	// TODO fix naming for different use cases
	x, y, ok := strings.Cut(message.Text, "=")
	if !ok {
		response = "Invalid data format"
		return
	}

	if err := b.storage.SetStatus(message.Chat.ID, model.DefaultStatus); err != nil {
		log.WithError(err).WithField("chatId", message.Chat.ID).Error("failed to remove chat's status")
		response = "Error occured. Youre next message should be considered as data sending. You should provide it once again"
		return
	}

	response = fmt.Sprintf("You've entered x = %s, y = %s", x, y)
}
