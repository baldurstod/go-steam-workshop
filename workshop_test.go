package steamworkshop_test

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"testing"

	steamworkshop "github.com/baldurstod/go-steam-workshop"
)

type Config struct {
	Steam `json:"steam"`
}
type Steam struct {
	ApiKey string `json:"api_key"`
}

var config = Config{}

func init() {
	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Failed to Unmarshal config: %v", err)
	}

	if config.ApiKey == "" {
		log.Fatalf("Steam api key is empty")
	}
}

func TestTf2(t *testing.T) {
	content, err := steamworkshop.QueryFiles(config.ApiKey, 440, "*",
		steamworkshop.WithReturnPreviews(),
		steamworkshop.WithReturnVoteData(),
		steamworkshop.WithReturnTags(),
		steamworkshop.WithReturnMetadata(),
		steamworkshop.WithReturnKvTags(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	response, ok := content["response"].(map[string]any)
	if !ok {
		t.Error(errors.New("missing response key in api response"))
		return
	}

	log.Println(response)

	j, _ := json.MarshalIndent(&response, "", "\t")

	err = os.WriteFile("./var/items_440.json", j, 0666)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetPublishedFileDetails(t *testing.T) {
	content, err := steamworkshop.GetPublishedFileDetails(config.ApiKey, []int64{3345507915, 606778917})
	if err != nil {
		t.Error(err)
		return
	}

	response, ok := content["response"].(map[string]any)
	if !ok {
		t.Error(errors.New("missing response key in api response"))
		return
	}

	j, _ := json.MarshalIndent(&response, "", "\t")

	err = os.WriteFile("./var/item_440.json", j, 0666)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestQueryFilesAllParams(t *testing.T) {
	content, err := steamworkshop.QueryFiles(config.ApiKey, 0, "*",
		steamworkshop.WithQueryType(steamworkshop.K_PublishedFileQueryType_RankedByPublicationDate),
		steamworkshop.WithPage(0),
		steamworkshop.WithNumPerPage(100),
		steamworkshop.WithCreatorAppId(10),
		steamworkshop.WithRequiredTags(""),
		steamworkshop.WithExcludedTags(""),
		steamworkshop.WithMatchAllTags(),
		steamworkshop.WithRequiredFlags(""),
		steamworkshop.WithExcludedFlags(""),
		steamworkshop.WithSearchText(""),
		steamworkshop.WithFileType(steamworkshop.K_PFI_MatchingFileType_Items),
		steamworkshop.WithChildPublishedFileId(0),
		steamworkshop.WithDays(0),
		steamworkshop.WithIncludeRecentVotesOnly(),
		steamworkshop.WithCacheMaxAgeSeconds(0),
		steamworkshop.WithLanguage(0),
		steamworkshop.WithTotalOnly(),
		steamworkshop.WithIdsOnly(),
		steamworkshop.WithReturnVoteData(),
		steamworkshop.WithReturnTags(),
		steamworkshop.WithReturnKvTags(),
		steamworkshop.WithReturnPreviews(),
		steamworkshop.WithReturnChildren(),
		steamworkshop.WithReturnShortDescription(),
		steamworkshop.WithReturnForSaleData(),
		steamworkshop.WithReturnMetadata(),
		steamworkshop.WithReturnPlaytimeStats(0),
	)
	if err != nil {
		t.Error(err)
		return
	}

	response, ok := content["response"].(map[string]any)
	if !ok {
		t.Error(errors.New("missing response key in api response"))
		return
	}

	log.Println(response)

	j, _ := json.MarshalIndent(&response, "", "\t")

	err = os.WriteFile("./var/items_all.json", j, 0666)
	if err != nil {
		t.Error(err)
		return
	}
}
