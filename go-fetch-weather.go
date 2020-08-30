package main

import (
	"os"
	"log"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
)

type weather struct {
	Meta struct {
		Source    string  `json:"source"`
		ExecTime  float64 `json:"exec_time"`
		Generated string  `json:"generated"`
	} `json:"meta"`
	Data []struct {
		Time string  `json:"time"`
		Temp int     `json:"temp"`
		Dwpt float64 `json:"dwpt"`
		Rhum int     `json:"rhum"`
		Prcp int     `json:"prcp"`
		Snow int     `json:"snow"`
		Wdir int     `json:"wdir"`
		Wspd float64 `json:"wspd"`
		Wpgt int     `json:"wpgt"`
		Pres float64 `json:"pres"`
		Tsun int     `json:"tsun"`
		Coco int     `json:"coco"`
	} `json:"data"`
}

// use godot package to load/read the .env file and
// return the value of the key
func environmentVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func fetchWeather(key string) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://api.meteostat.net/v2/stations/hourly" +
		"?station=10637&start=2020-02-01&end=2020-02-04", nil)

  if err != nil {
    log.Fatalf("Error requesting Meteostat data")
  }

	req.Header.Add("x-api-key", key)
	res, err := client.Do(req)

	fmt.Println("Response")
	fmt.Println(res)

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
			panic(err.Error())
	}

	var data weather
	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data)
	os.Exit(0)
}

func main() {
	fmt.Println("Hello world!")

	godotenv.Load(".env")

	if key := environmentVariable("METEOSTAT_API_KEY"); key == "" {
		fmt.Println("Please enter a METEOSTAT_API_KEY in a .env file, or use an environment variable")
	} else {
		fmt.Printf("Found key: %v\n", key)

		fetchWeather(key)
	}
}
