package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
)

type WeatherProvider struct {
	LocationCode string
}

type Weather struct {
	Status    string `json:"status"`
	Count     string `json:"count"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	ForeCasts []struct {
		City       string `json:"city"`
		AdCode     string `json:"adcode"`
		Province   string `json:"province"`
		ReportTime string `json:"reporttime"`
		Casts      []Cast `json:"casts"`
	} `json:"forecasts"`
}

type Cast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}

// NewWeatherProvider cityCode为高德天气api城市代码
func NewWeatherProvider(cityCode string) *WeatherProvider {
	return &WeatherProvider{
		LocationCode: cityCode,
	}
}

func (w *WeatherProvider) GetWeatherObj() (*Weather, error) {
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&extensions=all&key=%s", w.LocationCode, GlobalConf.ThirdParty.GaoDe.Key)
	content, err := httpGet(url)
	if err != nil {
		log.Println("http get failed", err)
		return nil, err
	}

	var weather Weather
	err = json.Unmarshal(content, &weather)
	if err != nil {
		log.Println("Unmarshal weather info failed", err)
		return nil, err
	}
	fmt.Printf("%+v\n", weather)
	if weather.InfoCode == "10000" {
		return &weather, nil
	} else {
		return nil, errors.New("get weather failed" + weather.Info)
	}

}

func (w *WeatherProvider) GetWeatherString() (string, error) {
	weather, err := w.GetWeatherObj()
	if err != nil {
		log.Println("get weather failed", err)
		return "", err
	}

	weatherStringTemplate := "%s天气预报：\n" +
		"今日（%s）：\n" +
		"白天天气：%s\n" +
		"平均气温：%s℃\n" +
		"晚上天气：%s\n" +
		"平均气温：%s℃\n" +
		"明日（%s）：\n" +
		"白天天气：%s\n" +
		"平均气温：%s℃\n" +
		"晚上天气：%s\n" +
		"平均气温：%s℃\n"

	weatherString := fmt.Sprintf(weatherStringTemplate, weather.ForeCasts[0].City, weather.ForeCasts[0].Casts[0].Date, weather.ForeCasts[0].Casts[0].DayWeather,
		weather.ForeCasts[0].Casts[0].DayTemp, weather.ForeCasts[0].Casts[0].NightWeather, weather.ForeCasts[0].Casts[0].NightTemp,
		weather.ForeCasts[0].Casts[1].Date, weather.ForeCasts[0].Casts[1].DayWeather, weather.ForeCasts[0].Casts[1].DayTemp,
		weather.ForeCasts[0].Casts[1].NightWeather, weather.ForeCasts[0].Casts[1].NightTemp)

	return weatherString, nil

}

func httpGet(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Println("get failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", err)
		return nil, err
	}
	return content, nil
}
