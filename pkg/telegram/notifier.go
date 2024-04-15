package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
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
	log.Infof("subscribers was notified about: %s", string(data))

	return nil
}
