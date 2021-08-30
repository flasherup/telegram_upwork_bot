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

	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	users := config.Telegram.Users

	updates, err := bot.GetUpdatesChan(u)
	for {
		select {
		case resp := <-ch:
			if len(users) == 0 {
				continue
			}
			if resp.Error == nil {
				p, exist := processors[resp.Feed.Id]
				if !exist {
					p = upwork.NewRSSProcessor()
					processors[resp.Feed.Id] = p
				}

				entries := p.Check(resp.Feed.Entries, config.Filters.SkipDuration)
				for _,v := range entries {
					for _,user := range users {
						msg := upworkEntryToBotMessage(&v, user)
						msg.ChatID = user.Id
						bot.Send(msg)
					}
				}
			}

		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}
			chatId := update.Message.Chat.ID
			if !isIdExist(users, chatId) {
				t := time.Unix(int64(update.Message.Date), 0)
				zone, offset := t.Zone()
				user := utils.User{
					Id: chatId,
					ZoneName:zone,
					ZoneOffset:offset,
				}
				users = append(users,user)
			}


			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func upworkEntryToBotMessage(entry *upwork.Entry, user utils.User) tgbotapi.MessageConfig {
	text := fmt.Sprintf("*%s*\n", entry.Title)
	loc := time.FixedZone(user.ZoneName, user.ZoneOffset)
	local := entry.Updated.In(loc)
	text += fmt.Sprintf("\n_%s_\n\n", local.Format("Mon Jan 2 15:04:05"))
	text += entry.Summary.Text + "\n"
	text += "\n"
	text += fmt.Sprintf("*Posteg on*: %s\n", entry.Summary.PostedOn)
	text += fmt.Sprintf("*Category*: %s\n", entry.Summary.Category)
	text += fmt.Sprintf("*Country*: %s\n", entry.Summary.Country)
	if entry.Summary.Budget > 0 {
		text += fmt.Sprintf("*Budget*: $%d\n", entry.Summary.Budget)
	}
	if len(entry.Summary.HourlyRange) > 0 {
		text += fmt.Sprintf("*Hourly Range*: $%g-$%g\n", entry.Summary.HourlyRange[0], entry.Summary.HourlyRange[1])
	}
	text += fmt.Sprintf("[More info...](%s)\n", entry.Summary.Link)
	msg := tgbotapi.NewMessage(user.Id, text)
	msg.ParseMode = "Markdown"
	return msg
}

func isIdExist(users []utils.User, id int64) bool {
	for _,user := range users {
		if user.Id == id {
			return true
		}
	}

	return false
}