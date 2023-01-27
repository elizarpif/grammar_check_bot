package grammar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Response struct {
	Fixed    string `json:"fixed"`
	Original string `json:"original"`
}

type Grammar struct {
	Url    string
	ApiKey string
	Host   string
}

func NewGrammar(url string, apiKey string, host string) *Grammar {
	return &Grammar{Url: url, ApiKey: apiKey, Host: host}
}

func (g Grammar) Check(text string) (string, error) {
	payload := strings.NewReader(fmt.Sprintf("{\n    \"text\": \"%s\"\n}", text))

	req, _ := http.NewRequest("POST", g.Url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", g.ApiKey)
	req.Header.Add("X-RapidAPI-Host", g.Host)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.Fixed, nil
}
