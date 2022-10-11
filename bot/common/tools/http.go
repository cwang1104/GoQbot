package tools

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"qbot/pkg/logger"
)

const (
	CqHttpBaseUrl = "http://127.0.0.1:5700"
)

func HttpGet(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("[http get failed][err:%v]", err)
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("[io readall failed][err:%v]", err)
		return nil, err
	}
	return content, nil
}

func HttpPost(url string, b []byte) ([]byte, error) {
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		logger.Log.Errorf("[http post failed][err:%v]", err)
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("[io read failed][err:%v]", err)
		return nil, err
	}
	return content, nil
}
