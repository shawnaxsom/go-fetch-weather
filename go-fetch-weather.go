package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type weatherData []struct {
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
}

type weather struct {
	Meta struct {
		Source    string  `json:"source"`
		ExecTime  float64 `json:"exec_time"`
		Generated string  `json:"generated"`
	} `json:"meta"`
	Data weatherData `json:"data"`
}

type stations struct {
	Meta struct {
		ExecTime  float64 `json:"exec_time"`
		Generated string  `json:"generated"`
	} `json:"meta"`
	Data []struct {
		ID   string `json:"id"`
		Name struct {
			En string `json:"en"`
		} `json:"name"`
		Country   string      `json:"country"`
		Region    string      `json:"region"`
		National  interface{} `json:"national"`
		Wmo       string      `json:"wmo"`
		Icao      string      `json:"icao"`
		Iata      interface{} `json:"iata"`
		Latitude  float64     `json:"latitude"`
		Longitude float64     `json:"longitude"`
		Elevation int         `json:"elevation"`
		Timezone  string      `json:"timezone"`
		Active    bool        `json:"active"`
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

var tr *http.Transport = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
var client *http.Client = &http.Client{Transport: tr}

func findStation(key string, query string) stations {
	url := fmt.Sprintf("https://api.meteostat.net/v2/stations/search?query=%v", strings.Replace(query, " ", "%20", -1))
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Error requesting Meteostat data")
	}

	req.Header.Add("x-api-key", key)
	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(body)

	if err != nil {
		panic(err.Error())
	}

	var data stations
	json.Unmarshal(body, &data)
	// fmt.Printf("Results: %v\n", data)
	// b, err := json.MarshalIndent(data, "", "  ")
	// fmt.Println(string(b))

	return data
}

func fetchWeather(key string, stationId string) weather {
	currentTime := time.Now()
	yesterday := currentTime.Format("2006-01-02")
	today := currentTime.Format("2006-01-02")
	fmt.Println("Current date: ", today)

	url := fmt.Sprintf("https://api.meteostat.net/v2/stations/hourly"+
		"?station=%v&start=%v&end=%v", stationId, yesterday, today)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Error requesting Meteostat data")
	}

	req.Header.Add("x-api-key", key)
	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(body)

	if err != nil {
		panic(err.Error())
	}

	var data weather
	json.Unmarshal(body, &data)
	return data
}

func filterResultsWithInfo(data weatherData) weatherData {
	// Meteostat includes all times of day, even those that haven't occurred yet.
	// Later times have all zeros for data besides time.
	var filteredData weatherData
	for i := range data {
		if data[i].Temp != 0 && data[i].Wdir != 0 && data[i].Wspd != 0 {
			filteredData = append(filteredData, data[i])
		}
	}

	if len(filteredData) > 0 {
		return filteredData
	}

	return data
}

func main() {
	stationName := flag.String("station", "Indianapolis", "Name of weather station to search for and retrieve weather from")

	flag.Parse()

	godotenv.Load(".env")

	if key := environmentVariable("METEOSTAT_API_KEY"); key == "" {
		log.Fatalf("Please enter a METEOSTAT_API_KEY in a .env file, or use an environment variable")
	} else {
		stationData := findStation(key, *stationName)
		// s, _ := json.MarshalIndent(stationData, "", "\t");
		// fmt.Print(string(s))
		fmt.Println("Station:", stationData.Data[0].Name.En)
		data := fetchWeather(key, stationData.Data[0].ID)
		filteredData := filterResultsWithInfo(data.Data)

		// b, err := json.MarshalIndent(filteredData, "", "  ")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// fmt.Println(string(b))

		lastEntry := filteredData[len(filteredData)-1]

		fmt.Printf("Latest Time: %v\n", lastEntry.Time)
		fmt.Printf("Latest Temp (C): %v\n", lastEntry.Temp)
		fahrenheit := float64(lastEntry.Temp)*(9.0/5.0) + 32
		fmt.Printf("Latest Temp (F): %v\n", fahrenheit)
		fmt.Printf("Latest Prcp: %v\n", lastEntry.Prcp)
	}
}
