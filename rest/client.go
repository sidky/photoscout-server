package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func BuildURL(baseURL string, parameters url.Values) (*url.URL, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	q := url.Query()

	for key, value := range parameters {
		for _, v := range value {
			q.Add(key, v)
		}
	}
	url.RawQuery = q.Encode()

	return url, nil
}

func Get(url *url.URL, response interface{}) error {
	fmt.Println(url.String())
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}
