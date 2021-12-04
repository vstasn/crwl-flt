package crwlrs

import (
	"crwlflt/m/v2/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Flt struct {
	Status       string `json:"status"`
	StatusAlias  string `json:"status_alias"`
	Number       string `json:"number"`
	PropertyType string `json:"property_type"`
	Price        string `json:"price"`
	PriceM2      string `json:"price_m2"`
	Floor        string `json:"floor"`
	Type         string `json:"type"`
	SquareTotal  string `json:"square_total"`
	Rooms        string `json:"rooms"`
	Section      string `json:"section"`
}

type Products struct {
	Flts []Flt `json:"products"`
}

type Data struct {
	Data Products `json:"data"`
}

func GetFlts() []Flt {
	url := fmt.Sprintf("%s&v=%d", config.AppConfig.FltURL, time.Now().UnixMilli())

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var data Data

	json.Unmarshal(byteValue, &data)

	var flats []Flt
	for _, flt := range data.Data.Flts {
		if flt.PropertyType == "flat" && flt.Section != "1" {
			flats = append(flats, flt)
		}
	}

	return flats
}
