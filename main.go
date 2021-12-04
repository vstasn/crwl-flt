package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"crwlflt/m/v2/config"
	"crwlflt/m/v2/crwlrs"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Flat struct {
	Id           int64
	Number       int64 `pg:",unique"`
	Floor        int8
	Rooms        int8
	SquareTotal  float64
	Section      int8
	Type         string
	PropertyType string
	Price        string
	PriceM2      string
	Status       string
	StatusAlias  string
}

func (b *Flat) ChangedStatus(newFlat *Flat) bool {
	return b.StatusAlias != newFlat.StatusAlias
}

func (b *Flat) ChangedPrice(newFlat *Flat) bool {
	return b.Price != newFlat.Price
}

type Change struct {
	Number   int64
	OldValue string
	NewValue string
}

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
	var changesPrices, changesStatuses []Change

	freeFlats := crwlrs.GetFlts()
	for _, flat := range freeFlats {

		dbFlat := new(Flat)
		err := db.Model(dbFlat).Where("number = ?", flat.Number).Select()

		inrec, _ := json.Marshal(flat)
		var values map[string]interface{}
		json.Unmarshal(inrec, &values)

		if err == nil {
			query := db.Model(&values).TableExpr("flats").Where("id = ?", dbFlat.Id)
			_, err = query.Update()
			if err != nil {
				continue
			}
			newFlat := new(Flat)
			err = db.Model(newFlat).Where("id = ?", dbFlat.Id).Select()
			if err != nil {
				continue
			}

			if dbFlat.ChangedPrice(newFlat) {
				changesPrices = append(changesPrices, Change{Number: dbFlat.Number, OldValue: dbFlat.Price, NewValue: newFlat.Price})
			}

			if dbFlat.ChangedStatus(newFlat) {
				changesStatuses = append(changesStatuses, Change{Number: dbFlat.Number, OldValue: dbFlat.StatusAlias, NewValue: newFlat.StatusAlias})
			}
		}

		if err == pg.ErrNoRows {
			db.Model(&values).TableExpr("flats").Insert()
		}
	}

	if len(changesPrices) > 0 || len(changesStatuses) > 0 {

		formatMsg := func(msg string, changes []Change) string {
			for _, change := range changes {
				msg = fmt.Sprintf("%s\n%s", msg, fmt.Sprintf("%d %s %s", change.Number, change.OldValue, change.NewValue))
			}
			return msg
		}
		var text string
		if len(changesStatuses) > 0 {
			text = formatMsg("*Изменения в статусах:*\n_Номер Old New_", changesStatuses)
		}
		if len(changesPrices) > 0 {
			text1 := formatMsg("\n*Изменения в ценах:*\n_Номер Old New_", changesPrices)
			text = fmt.Sprintf("%s%s", text, text1)
		}

		msg := tgbotapi.NewMessage(config.AppConfig.TelegramChatId, text)
		msg.ParseMode = tgbotapi.ModeMarkdownV2

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
	}

	log.Println("End programm")
}
