package main

import (
	"fmt"
	"github.com/flasherup/telegrum_upwork_bot/upwork"
	"github.com/flasherup/telegrum_upwork_bot/utils"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)



func main() {
	config, err := utils.LoadConfig("config.yml")
	if err !=  nil {
		log.Panic(err)
	}

	processors := map[string]*upwork.RSSProcessor{}

	uw := upwork.NewUpwork(config.Upwork)
	ch := uw.Run(time.Second *15)

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	chatIds := []int64{388017091}

	updates, err := bot.GetUpdatesChan(u)
	for {
		select {
		case resp := <-ch:
			if len(chatIds) == 0 {
				continue
			}
			if resp.Error == nil {
				p, exist := processors[resp.Feed.Id]
				if !exist {
					p = upwork.NewRSSProcessor()
					processors[resp.Feed.Id] = p
				}

				entries := p.Check(resp.Feed.Entries)
				for _,v := range entries {
					msg := upworkEntryToBotMessage(&v, 0)
					for _,chatId := range chatIds {
						msg.ChatID =chatId
						bot.Send(msg)
					}
				}
			}

		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}
			chatId := update.Message.Chat.ID
			if !isIdExist(chatIds, chatId) {
				chatIds = append(chatIds,chatId)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			fmt.Println("chatId", chatId)
			bot.Send(msg)
		}
	}
}

func upworkEntryToBotMessage(entry *upwork.UWEntry, chatId int64) tgbotapi.MessageConfig {
	summary := upwork.ParseSummary(entry.Summary)
	text := fmt.Sprintf("*%s*\n", entry.Title)
	update, err := time.Parse("2006-01-02T15:04:05+00:00", entry.Updated)
	if err != nil {
		fmt.Println("Time parse Error")
	} else {
		text += fmt.Sprintf("\n_%s_\n\n", update.Format("Mon Jan 2 15:04:05"))
	}
	text += summary.Text + "\n"
	text += "\n"
	text += fmt.Sprintf("*Posteg on*: %s\n", summary.PostedOn)
	text += fmt.Sprintf("*Category*: %s\n", summary.Category)
	text += fmt.Sprintf("*Country*: %s\n", summary.Country)
	if summary.Budget > 0 {
		text += fmt.Sprintf("*Budget*: $%d\n", summary.Budget)
	}
	if len(summary.HourlyRange) > 0 {
		text += fmt.Sprintf("*Hourly Range*: $%g-$%g\n", summary.HourlyRange[0], summary.HourlyRange[1])
	}
	text += fmt.Sprintf("[More info...](%s)\n", summary.Link)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	return msg
}

func isIdExist(ids []int64, id int64) bool {
	for _,v := range ids {
		if v == id {
			return true
		}
	}

	return false
}