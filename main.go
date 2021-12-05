package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"crwlflt/m/v2/config"
	"crwlflt/m/v2/crwlrs"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Flat)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}

func worker(db orm.DB) {
	var changes []Change

	freeFlats := crwlrs.GetFlts()
	for _, flat := range freeFlats {

		flt := new(Flat)
		err := db.Model(flt).Where("number = ?", flat.Number).Select()

		values := GetFields(flat)
		query := db.Model(&values).TableExpr("flats")

		if err == nil {
			query.Where("id = ?", flt.Id).Update()
			itemChanges := GetChanges(flt.Number, values, ConvertFieldsToUnderscore(GetFields(flt)))
			log.Println(itemChanges)
			changes = append(changes, itemChanges...)
		}

		if err == pg.ErrNoRows {
			query.Insert()
		}
	}

	if len(changes) > 0 {
		formatMsg := func(changes []Change) string {
			var msg string
			for _, change := range changes {
				msg = fmt.Sprintf("%s\n%s", msg, fmt.Sprintf("%d %s %s %s", change.Number, change.Field, change.OldValue, change.NewValue))
			}
			return msg
		}

		text := fmt.Sprintf("<b>Список изменений:</b>\n<b>Номер Field Old New</b>%s", formatMsg(changes))

		msg := tgbotapi.NewMessage(config.AppConfig.TelegramChatId, text)
		msg.ParseMode = tgbotapi.ModeHTML

		bot, err := tgbotapi.NewBotAPI(config.AppConfig.TelegramApiToken)
		if err != nil {
			panic(err)
		}

		bot.Send(msg)
	}
}

func main() {
	log.Println("Start programm")
	log.Println(fmt.Sprintf("Run Every %d", config.AppConfig.RunTime))

	opt, err := pg.ParseURL(config.AppConfig.PgURL)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	defer db.Close()

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	worker(db)

	for range time.Tick(time.Duration(config.AppConfig.RunTime) * time.Second) {
		log.Println("run worker")
		worker(db)
		count, _ := db.Model((*Flat)(nil)).Count()
		log.Println(fmt.Sprintf("Count flats in db: %d", count))
	}

	log.Println("End programm")
}
