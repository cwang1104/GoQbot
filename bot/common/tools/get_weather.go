package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"qbot/pkg/utils"
	"strings"
	"sync"
	"time"
)

var (
	weatherCache = map[string]*Weather{}
	mu           sync.Mutex
)

type WeatherProvider struct {
	LocationCode string
	WeatherKey   string
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
		WeatherKey:   GetWeatherKey(cityCode),
	}
}

func (w *WeatherProvider) GetWeatherObj() (*Weather, error) {

	//如果当天当城市的天气已经被缓存 则直接返回
	if _, exist := GetWeatherCache(w.WeatherKey); exist {
		weather, _ := GetWeatherCache(w.WeatherKey)
		return weather, nil
	}
	//从api接口请求
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&extensions=all&key=%s", w.LocationCode, utils.GlobalConf.ThirdParty.GaoDe.Key)
	content, err := HttpGet(url)
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

	if weather.InfoCode == "10000" {
		//将今天的天气信息存入map中
		AddWeatherCache(w.WeatherKey, &weather)
		//将昨天的删除
		DeleteExpireWeather()
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

func AddWeatherCache(key string, weather *Weather) {
	mu.Lock()
	defer mu.Unlock()
	weatherCache[key] = weather
}

func GetWeatherCache(key string) (weather *Weather, exist bool) {
	mu.Lock()
	defer mu.Unlock()
	weather, exist = weatherCache[key]
	return
}

func DelWeatherCache(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(weatherCache, key)
}

//GetWeatherKey 根据当天时间和城市代码生成map-key
func GetWeatherKey(cityCode string) string {
	timeString := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s*%s", cityCode, timeString)
}

func DeleteExpireWeather() {
	for k, _ := range weatherCache {
		parts := strings.Split(k, "*")[1]
		exDate := time.Unix(time.Now().Unix()-60*60*24, 0).Format("2006-01-02")
		if parts == exDate {
			DelWeatherCache(k)
		}
	}
}
