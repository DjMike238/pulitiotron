package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

const TOKEN = ""
const BOT_NAME = "pulitiotron"
const ADMIN = 14870908

var urls map[string][]string
var supported string

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(TOKEN),
	}
}

func (b *bot) handleInline(iq echotron.InlineQuery) {
	var url string
	var results []echotron.InlineQueryResult

	if iq.Query != "" {
		url = getCleanURL(iq.Query, urls)
	}

	switch url {
	case "":
		b.AnswerInlineQuery(
			iq.ID,
			results,
			&echotron.InlineQueryOptions{
				CacheTime:         1,
				SwitchPmText:      "Need help?",
				SwitchPmParameter: "start",
			},
		)

	case "unsupported":
		b.AnswerInlineQuery(
			iq.ID,
			results,
			&echotron.InlineQueryOptions{
				CacheTime:         1,
				SwitchPmText:      "Unsupported URL (click to learn more)",
				SwitchPmParameter: "unsupported",
			},
		)

	default:
		results = append(results, echotron.InlineQueryResultArticle{
			Type:        echotron.ARTICLE,
			ID:          fmt.Sprintf("%d", time.Now().Unix()),
			Title:       "Send clean URL",
			Description: url,
			InputMessageContent: echotron.InputTextMessageContent{
				MessageText: url,
			},
		})

		b.AnswerInlineQuery(
			iq.ID,
			results,
			&echotron.InlineQueryOptions{
				CacheTime: 1,
			},
		)
	}
}

func (b *bot) handleMessage(m echotron.Message) {
	msg := strings.TrimPrefix(m.Text, "/start ")

	switch msg {
	case "unsupported":
		b.SendMessage(
			"Hey\\! Looks like you tried to send me an unsupported URL\\.\n"+
				"\n"+
				"If you'd like to ask for support for a new URL type\\, you can open an issue or send a PR on the bot\\'s [GitHub page](https://github.com/DjMike238/pulitiotron)\\.",
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)

	case "start":
		fallthrough

	case "/start":
		b.SendMessage(
			"Welcome to *Pulitiotron*\\!\n"+
				"\n"+
				"Just write my nickname in the message box\\, put an URL after it and it will be magically cleaned\\!\n"+
				"\n"+
				"Bot made by @Dj\\_Mike238\\.\n"+
				"This bot is [open source](https://github.com/DjMike238/pulitiotron)\\!",
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)

	case "/reload":
		if b.chatID == ADMIN {
			b.SendMessage("Reloading URLs...", b.chatID, nil)
			urls = loadURLs()
			supported = createSupportedList(urls)
			b.SendMessage("URLs reloaded successfully.", b.chatID, nil)
		}

	case "/supported":
		b.SendMessage(
			fmt.Sprintf("This bot currently supports\\:\n\n%s", supported),
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)
	}
}

func (b *bot) Update(update *echotron.Update) {
	defer avertCrysis()

	if update.InlineQuery != nil {
		b.handleInline(*update.InlineQuery)
	} else if update.Message != nil {
		b.handleMessage(*update.Message)
	}
}

func avertCrysis() {
	if err := recover(); err != nil {
		log.Println(err)
		log.Println("Thread recovered. Crysis averted.")
	}
}

func main() {
	urls = loadURLs()
	supported = createSupportedList(urls)
	dsp := echotron.NewDispatcher(TOKEN, newBot)
	log.Println(dsp.Poll())
}
