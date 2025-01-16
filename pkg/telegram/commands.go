package telegram

import (
	"alert_bot/pkg/domain"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var role2commands = map[domain.Role][]tgbotapi.BotCommand{
	domain.UnknownRole: {
		{Command: startCommand, Description: "start"},
		{Command: healthCommand, Description: "health check"},
		{Command: subscribeCommand, Description: "subscribe"},
		{Command: unsubscribeCommand, Description: "unsubscribe"},
		{Command: subscriptionCheckCommand, Description: "check subscription"},
		{Command: sendDataCommand, Description: "send data"},
	},
}

func (b *Bot) setCommands(_ context.Context, chatID int64, role domain.Role) error {
	commands, ok := role2commands[role]
	if !ok {
		return fmt.Errorf("no commands found for role %s", role)
	}

	commandCfg := tgbotapi.NewSetMyCommandsWithScopeAndLanguage(
		tgbotapi.BotCommandScope{
			Type:   "chat",
			ChatID: chatID,
		},
		"",
		commands...,
	)

	if _, err := b.bot.Request(commandCfg); err != nil {
		return fmt.Errorf("b.bot.Request: %w", err)
	}
	return nil
}
