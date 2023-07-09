package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"log"

	"errors"

	"github.com/cleanshaven/noaaweather"
	"github.com/godbus/dbus/v5"
)

const (
	UserAgentKey      = "User-Agent"
	UserAgent         = "Mozilla/5.0 (compatible; fpweb)"
	PointsUrlTemplate = "https://api.weather.gov/points/%s,%s"
	AlertUrlTemplate  = "https://api.weather.gov/alerts/active?point=%s,%s"
)

type serverObject struct{}

func (server serverObject) GetIcon(iconUrl string) (iconFile string, dbusError *dbus.Error) {
	iconFile = ""

	dbusError = nil

	newIconUrl, err := fixWeatherIconUrl(iconUrl)
	if err != nil {
		dbusError = dbus.MakeFailedError(err)
		return
	}

	icon, err := getNOAAInfo(newIconUrl)
	if err != nil {
		dbusError = dbus.MakeFailedError(err)
		return
	}

	file, err := ioutil.TempFile("", "i3weatherblock*.jpeg")
	if err != nil {
		dbusError = dbus.MakeFailedError(err)
		return
	}

	defer file.Close()

	_, err = file.Write(icon)
	if err != nil {
		dbusError = dbus.MakeFailedError(err)
	}

	iconFile = file.Name()
	return

}

// this is to fix what appears to be a bug in the icon sent by the NOAA Weather api
// There is a ,0 sitting at the end of the path in the url
// ex: https://api.weather.gov/icons/land/day/haze,0?size=small needs to be
// https://api.weather.gov/icons/land/day/haze?size=small
func fixWeatherIconUrl(url string) (string, error) {
	re, err := regexp.Compile(`,[^?]*`)
	if err != nil {
		return url, err
	}

	newUrl := re.ReplaceAllString(url, "")
	return newUrl, nil

}

func (server serverObject) GetWeather(longitude, latitude string) (forecast, alerts string, dbusError *dbus.Error) {
	log.Printf("got request %s longitude %s latitude", longitude, latitude)
	forecast = ""
	alerts = ""
	dbusError = nil

	forecastChannel := make(chan string)
	errorChannel := make(chan error)
	alertChannel := make(chan string)

	go concurrentGetForecast(latitude, longitude, forecastChannel, errorChannel)
	go concurrentGetAlert(latitude, longitude, alertChannel, errorChannel)

	for i := 0; i < 2; i++ {
		select {
		case forecast = <-forecastChannel:
		case alerts = <-alertChannel:
		case err := <-errorChannel:
			dbusError = dbus.MakeFailedError(err)
			return
		}
	}
	return
}

func concurrentGetAlert(latitude, longitude string, alertChannel chan string, errorChannel chan error) {
	alerts, err := getNOAAInfo(getAlertsUrl(latitude, longitude))
	if err != nil {
		errorChannel <- err
		return
	}
	alertChannel <- string(alerts)
}

func concurrentGetForecast(latitude, longitude string, forecastChannel chan string, errorChannel chan error) {
	forecastUrl, err := getForecastUrl(latitude, longitude)
	if err != nil {
		errorChannel <- err
		return
	}

	forecast, err := getForecast(forecastUrl)
	if err != nil {
		errorChannel <- err
		return
	}
	forecastChannel <- forecast
}

func getPointsUrl(latitude, longitude string) string {
	return fmt.Sprintf(PointsUrlTemplate, latitude, longitude)
}

func getAlertsUrl(latitude, longitude string) string {
	return fmt.Sprintf(AlertUrlTemplate, latitude, longitude)
}

func getNOAAInfo(requestUrl string) ([]byte, error) {
	request, err := getRequest(requestUrl)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return []byte{}, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func getForecastUrl(latitude, longitude string) (string, error) {
	location, err := getNOAAInfo(getPointsUrl(latitude, longitude))
	if err != nil {
		return "", err
	}

	var locationJson noaaweather.LocationJson
	err = json.Unmarshal(location, &locationJson)
	if err != nil {
		return "", err
	}
	return locationJson.Properties.ForecastHourly, nil
}

func getForecast(url string) (string, error) {
	result, err := getNOAAInfo(url)
	return string(result), err
}

func getRequest(requestUrl string) (*http.Request, error) {
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add(UserAgentKey, UserAgent)
	return request, nil
}
