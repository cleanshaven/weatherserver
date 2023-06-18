package main

import (
	"log"
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var server serverObject = serverObject{}

func TestIconLoad(t *testing.T) {
	fileName, err := server.GetIcon(`https://api.weather.gov/icons/land/day/haze,0?size=small`)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf(fileName)
}

func TestUrl(t *testing.T) {
	myUrl, err := url.Parse(`https://api.weather.gov/icons/land/day/haze,0?size=small`)
	if err != nil {
		log.Printf(err.Error())
	}

	spew.Dump(myUrl)
	log.Println(myUrl.Path)
	log.Println(myUrl.RawPath)
	log.Println(myUrl.RawQuery)
}

func TestIconFix(t *testing.T) {
	newString, err := fixWeatherIconUrl(`https://api.weather.gov/icons/land/day/haze,0?size=small`)
	if err != err {
		t.Fatalf(err.Error())
	}

	if newString != `https://api.weather.gov/icons/land/day/haze?size=small` {
		t.Fatalf("Wanted %s got %s",
			`https://api.weather.gov/icons/land/day/haze?size=small`,
			newString)
	}
}
