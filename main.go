package main

import (
	_ "embed"
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

const (
	BOT_NAME = "pulitiotron"
	ADMIN    = 14870908
)

var (
	urls      map[string][]string
	supported string

	//go:embed token
	token string
)

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) handleInline(iq echotron.InlineQuery) {
	var (
		url     string
		urlType URLType

		results []echotron.InlineQueryResult
	)

	if iq.Query != "" {
		urlType, url = getCleanURL(iq.Query, urls)
	}

	switch urlType {
	case Empty:
		b.AnswerInlineQuery(
			iq.ID,
			results,
			&echotron.InlineQueryOptions{
				CacheTime:         1,
				SwitchPmText:      "Need help?",
				SwitchPmParameter: "start",
			},
		)
		return

	case Unsupported:
		b.AnswerInlineQuery(
			iq.ID,
			results,
			&echotron.InlineQueryOptions{
				CacheTime:         1,
				SwitchPmText:      "Unsupported URL (click to learn more)",
				SwitchPmParameter: "unsupported",
			},
		)
		return

	// When trying to clean a Twitter URL, an option to send a
	// vxTwitter (vxtwitter.com) version of the URL will be added.
	case Twitter:
		vxURL := strings.Replace(url, "twitter", "vxtwitter", 1)

		results = append(results, echotron.InlineQueryResultArticle{
			Type:        echotron.InlineArticle,
			ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:       "Send vxTwitter URL",
			Description: vxURL,
			InputMessageContent: echotron.InputTextMessageContent{
				MessageText: vxURL,
			},
		})
	}

	results = append(results, echotron.InlineQueryResultArticle{
		Type:        echotron.InlineArticle,
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
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

	case "start", "/start":
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
	dsp := echotron.NewDispatcher(token, newBot)
	log.Println(dsp.Poll())
}
