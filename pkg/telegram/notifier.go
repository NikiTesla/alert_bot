package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) notifySubscribers(data []byte) error {
	subscribersUids, err := b.storage.GetSubscribersUids()
	if err != nil {
		return fmt.Errorf("cannot get subscribers' uids, err: %w", err)
	}

	var msg tgbotapi.MessageConfig
	for _, subscriberUid := range subscribersUids {
		msg = tgbotapi.NewMessage(subscriberUid, string(data))

		if _, err := b.bot.Send(msg); err != nil {
			return fmt.Errorf("cannot send message to chatId %d, err: %w", subscriberUid, err)
		}
	}

	return nil
}

func (b *Bot) notifySubscribersWithImage(imageFilename string) error {
	subscribersUids, err := b.storage.GetSubscribersUids()
	if err != nil {
		return fmt.Errorf("cannot get subscribers' uids, err: %w", err)
	}

	var msg tgbotapi.PhotoConfig
	for _, subscriberUid := range subscribersUids {
		msg = tgbotapi.NewPhoto(subscriberUid, tgbotapi.FilePath(imageFilename))

		if _, err := b.bot.Send(msg); err != nil {
			return fmt.Errorf("cannot send message to chatId %d, err: %w", subscriberUid, err)
		}
	}

	return nil
}
