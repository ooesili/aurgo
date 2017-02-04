package aur

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func New(baseURL string) (API, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return API{}, err
	}

	return API{
		scheme: url.Scheme,
		host:   url.Host,
	}, nil
}

type API struct {
	scheme string
	host   string
}

func (a API) Version(packageName string) (string, error) {
	infoURL := a.buildInfoURL(packageName)

	var infoResponse infoResponseType
	err := jsonGet(infoURL, &infoResponse)
	if err != nil {
		return "", err
	}

	results := infoResponse.Results
	if len(results) == 0 {
		return "", fmt.Errorf("package not found: %s", packageName)
	}

	version := results[0].Version

	return version, nil
}

func jsonGet(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (a API) buildInfoURL(packageName string) string {
	query := url.Values{
		"v":     {"5"},
		"type":  {"info"},
		"arg[]": {packageName},
	}

	infoURL := url.URL{
		Scheme:   a.scheme,
		Host:     a.host,
		Path:     "/rpc",
		RawQuery: query.Encode(),
	}

	return infoURL.String()
}

type infoResponseType struct {
	Results []struct {
		Version string `json:"Version"`
	} `json:"results"`
}
