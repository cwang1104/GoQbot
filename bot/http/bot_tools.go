package http

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseUrl = "http://127.0.0.1:5700"
)

func httpPost(url string, b []byte) ([]byte, error) {
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println("Post failed:", err)
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
