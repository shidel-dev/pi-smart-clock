package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pi-clock/config"
	"strconv"
	"time"
)

type WeatherService struct {
	ApiKey  string
	City    string
	Country string
	Url     string
	Units   string
}

type Day struct {
	Date     int64 `json:"dt"`
	Temp     map[string]float64
	pressure float64
	Weather  []struct {
		id          int
		Main        string `json:"main"`
		Description string `json:"description"`
		icon        string
	} `json:"weather"`
	speed   interface{}
	Degrees int `json:"deg"`
	clouds  int
}

type DailyForcast struct {
	Day  string
	Temp struct {
		Min int
		Max int
	}

	Conditions string
}

func (d *DailyForcast) String() string {

	return string(d.Day[0:3]) + " " + strconv.Itoa(d.Temp.Min) + "-" + strconv.Itoa(d.Temp.Max)
}

func NewWeatherService() *WeatherService {
	service := &WeatherService{
		ApiKey:  config.Vars.OpenWeatherKey,
		City:    config.Vars.City,
		Country: config.Vars.Country,
		Url:     "http://api.openweathermap.org/data/2.5/forecast/daily",
		Units:   config.Vars.Units,
	}

	return service
}

func (w *WeatherService) GetCurrentForcast() string {
	values := &url.Values{}

	values.Add("q", w.City+","+w.Country)
	values.Add("APPID", w.ApiKey)
	values.Add("units", w.Units)

	weatherUrl, _ := url.Parse(w.Url)
	weatherUrl.RawQuery = values.Encode()

	res, err := http.Get(weatherUrl.String())

	fmt.Println(weatherUrl.String())

	if err != nil {
		fmt.Printf("%#v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("could not read opean weather responce")
	}

	var forcastJSON struct {
		cod     string
		message string
		city    interface{}
		cnt     interface{}
		coord   interface{}
		List    []Day `json:"list"`
	}

	fmt.Println(string(body))

	err = json.Unmarshal(body, &forcastJSON)

	if err != nil {
		fmt.Println(err.Error())
	}

	var forcast string

	for i, day := range forcastJSON.List {
		if i <= 3 {
			d := &DailyForcast{
				Day: time.Unix(day.Date, 0.00).Weekday().String(),
				Temp: struct {
					Min int
					Max int
				}{
					Min: int(day.Temp["min"]),
					Max: int(day.Temp["max"]),
				},
				Conditions: day.Weather[0].Description,
			}

			forcast = forcast + d.String() + "\n"
		}
	}

	return forcast
}
