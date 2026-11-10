package steamworkshop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const retries = 10

func getApi(u *url.URL) (map[string]any, error) {
	resp, err := getPage(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get api response in getApi: <%w>", err)
	}
	defer resp.Body.Close()

	response := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json in getApi: <%w>", err)
	}

	return response, nil
}

func getPage(url string) (*http.Response, error) {
	var resp *http.Response
	for range retries {

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 429 { //Too Many Requests
			time.Sleep(60 * time.Second)
			continue
		}

		if resp.StatusCode != 200 { //Everything except 429 and 200
			return nil, fmt.Errorf("http request returned HTTP status code: %d", resp.StatusCode)
		}
		break
	}

	return resp, nil
}

func postApi(u *url.URL, headers map[string]string, body string) (map[string]any, error) {
	resp, err := postPage(u.String(), headers, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get api response in postApi: <%w>", err)
	}
	defer resp.Body.Close()

	response := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json in postApi: <%w>", err)
	}

	return response, nil
}

func postPage(url string, headers map[string]string, body string) (*http.Response, error) {
	var resp *http.Response
	for range retries {
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		for k, v := range headers {
			req.Header.Add(k, v)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 429 { //Too Many Requests
			time.Sleep(60 * time.Second)
			continue
		}

		if resp.StatusCode != 200 { //Everything except 429 and 200
			return nil, fmt.Errorf("http request returned HTTP status code: %d", resp.StatusCode)
		}
		break
	}

	return resp, nil
}
