package main

import (
	"log"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

const TOKEN = "1878453925:AAGPYaw4QkTKKScYd-xRPb_-Lrfu1yDWFsg"
const BOT_NAME = "pulitiotron"

var urls map[string][]string

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

	if url != "" {
		results = append(results, echotron.InlineQueryResultArticle{
			Type: echotron.ARTICLE,
			ID: url,
			Title: "Send clean URL",
			Description: url,
			InputMessageContent: echotron.InputTextMessageContent{
				MessageText: url,
			},
		})

		b.AnswerInlineQuery(iq.ID, results, nil)
		return
	}

	// FIXME: Show "Unsupported URL" button when an unsupported URL is sent.
	b.AnswerInlineQuery(
		iq.ID,
		results,
		&echotron.InlineQueryOptions{
			SwitchPmText: "Unsupported URL",
			SwitchPmParameter: "start",
		},
	)
}

func (b *bot) Update(update *echotron.Update) {
	defer avertCrysis()

	if update.InlineQuery != nil {
		b.handleInline(*update.InlineQuery)
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

	dsp := echotron.NewDispatcher(TOKEN, newBot)
	log.Println(dsp.Poll())
}
