package telegram

import log "github.com/sirupsen/logrus"

func (b *Bot) notifySubscribers(data []byte) error {
	log.Infof("subscribers was notified about: %s", string(data))

	return nil
}
