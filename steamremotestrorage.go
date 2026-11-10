package steamworkshop

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

const steamRemoteStorageGetPublishedFileDetails = "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"

func GetPublishedFileDetails(apiKey string, fileIds []int64) (map[string]any, error) {

	urlParams := url.Values{
		"key": {apiKey},
	}

	body := url.Values{
		"itemcount": {strconv.Itoa(len(fileIds))},
	}

	for i, fileId := range fileIds {
		body["publishedfileids["+strconv.Itoa(i)+"]"] = []string{strconv.FormatInt(fileId, 10)}
	}

	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	u, err := url.Parse(steamRemoteStorageGetPublishedFileDetails)
	if err != nil {
		return nil, errors.New("can't parse steamApiGetPublishedFileDetails")
	}
	u.RawQuery = urlParams.Encode()
	response, err := postApi(u, headers, body.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get api response in readWorkshopPage: <%w>", err)
	}

	return response, nil
}
