package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	Coord struct {
		Lon float64
		Lat float64
	}
	// Weather json.RawMessage
	// Base    json.RawMessage
	// Main    json.RawMessage
}

func main() {
	apiURL := "https://samples.openweathermap.org/data/2.5/weather?q=London,uk&appid=b6907d289e10d714a6e88b30761fae22"
	fmt.Println(apiURL)
	resp, err := http.Get(apiURL)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonResult := Response{}

	decoder := json.NewDecoder(resp.Body)
	jsonerr := decoder.Decode(&jsonResult)

	// json.Unmarshal(resp.Body)

	fmt.Println(jsonerr)
	fmt.Println(jsonResult)
	fmt.Println(jsonResult.Coord)
}
