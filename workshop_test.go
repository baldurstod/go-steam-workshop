package steamworkshop_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	steamworkshop "github.com/baldurstod/go-steam-workshop"
)

type Config struct {
	Database `json:"database"`
	Steam    `json:"steam"`
}
type Steam struct {
	ApiKey string `json:"api_key"`
}
type Database struct {
	DataSource string `json:"datasource"`
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
	testAppId(t, 440)
}

func TestSFM(t *testing.T) {
	testAppId(t, 1840)
}

func testAppId(t *testing.T, appId uint) {
	content, err := steamworkshop.QueryFiles(config.ApiKey, appId, "*",
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

	err = os.WriteFile("./var/published_files_"+strconv.FormatUint(uint64(appId), 10)+".json", j, 0666)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestRefreshTf2(t *testing.T) {
	db := openPostgre(config.DataSource)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := steamworkshop.RefreshPublishedFiles(ctx, db, config.ApiKey, 440, false,
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

	err = os.WriteFile("./var/published_files_all.json", j, 0666)
	if err != nil {
		t.Error(err)
		return
	}
}

var dbs = []*sql.DB{}

func openPostgre(dataSourceName string) *sql.DB {
	var err error
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// db.Open() only creates a connection pool, and doesn't actually establish
	// a connection. To ensure the connection works you need to do *something*
	// with a connection.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbs = append(dbs, db)

	return db
}

func closePostgre() {
	for _, db := range dbs {
		if db != nil {
			db.Close()
		}
	}
}
