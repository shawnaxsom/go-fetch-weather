package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

func fetchWeather(key string) weather {
	currentTime := time.Now()
	yesterday := currentTime.Format("2006-01-02")
	today := currentTime.Format("2006-01-02")
  fmt.Println("Current date: ", today)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://api.meteostat.net/v2/stations/hourly"+
		"?station=10637&start=%v&end=%v", yesterday, today)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Error requesting Meteostat data")
	}

	req.Header.Add("x-api-key", key)
	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
 	fmt.Println(body)

	if err != nil {
		panic(err.Error())
	}

	var data weather
	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data)
	b, err := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))

	return data
}

func main() {
	godotenv.Load(".env")

	if key := environmentVariable("METEOSTAT_API_KEY"); key == "" {
		log.Fatalf("Please enter a METEOSTAT_API_KEY in a .env file, or use an environment variable")
	} else {
		data := fetchWeather(key)

		lastEntry := data.Data[len(data.Data) - 1]

		fmt.Println()
		fmt.Printf("Latest Time: %v\n", lastEntry.Time)
		fmt.Printf("Latest Temp: %v\n", lastEntry.Temp)
		fmt.Printf("Latest Prcp: %v\n", lastEntry.Prcp)
	}
}
