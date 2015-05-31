package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type configJSON struct {
	OpenWeatherKey string
	Country        string
	City           string
	Units          string
}

func load() configJSON {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Could not open config file!")
	}

	var config configJSON
	json.Unmarshal(b, &config)
	return config
}

var Vars configJSON = load()
