package steamworkshop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	_ "github.com/lib/pq"
)

// Insert a published file in database
// Return true if the file was inserted or updated, false if it was up to date or an error occured
func insertPublishedFile(db *sql.DB, publishedFile map[string]any) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var workshopFile WorkshopFile
	err := mapstructure.Decode(publishedFile, &workshopFile)
	if err != nil {
		return false, fmt.Errorf("an error occured while decoding workshop file : <%w>", err)
	}

	var tags = []string{}

	for _, t := range workshopFile.Tags {
		tags = append(tags, t.Tag)
	}

	j, err := json.Marshal(&publishedFile)
	if err != nil {
		return false, fmt.Errorf("failed to marshal file: <%w>", err)
	}

	fileId, err := strconv.ParseUint(workshopFile.Publishedfileid, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert Publishedfileid "+workshopFile.Publishedfileid+": <%w>", err)
	}

	fileSize, err := strconv.ParseUint(workshopFile.FileSize, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert file size "+workshopFile.FileSize+": <%w> for file "+workshopFile.Publishedfileid, err)
	}

	row := db.QueryRow(`SELECT time_updated from published_files where publishedfileid = $1;`, fileId)

	var timeUpdated uint64

	err = row.Scan(&timeUpdated)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to scan row in insertPublishedFile: <%w>", err)
	}

	if err != sql.ErrNoRows && timeUpdated == workshopFile.TimeUpdated {
		return false, nil
	}

	_, err = db.Exec(`INSERT INTO published_files (
	publishedfileid, title, time_created, time_updated, creator, tags, file_size, file_url, preview_url, subscriptions, creator_appid, consumer_appid, detail )
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13 )
	ON CONFLICT (publishedfileid) DO UPDATE SET
	title = $2,
	time_created = $3,
	time_updated = $4,
	creator = $5,
	tags = $6,
	file_size = $7,
	file_url = $8,
	preview_url = $9,
	subscriptions = $10,
	creator_appid = $11,
	consumer_appid = $12,
	detail = $13
	`,
		fileId,
		workshopFile.Title,
		workshopFile.TimeCreated,
		workshopFile.TimeUpdated,
		workshopFile.Creator,
		tags,
		fileSize,
		workshopFile.FileUrl,
		workshopFile.PreviewUrl,
		workshopFile.Subscriptions,
		workshopFile.CreatorAppid,
		workshopFile.ConsumerAppid,
		j,
	)

	if err != nil {
		return false, fmt.Errorf("failed to insert workshop file "+workshopFile.Publishedfileid+" : <%w>", err)
	}

	return true, nil
}

type FileFilter struct {
	Name  string
	Tags  []string
	Types []string
}

type FileSort struct {
	Field     string `json:"field" mapstructure:"field"`
	Ascending bool   `json:"ascending" mapstructure:"ascending"`
}

type FileParams struct {
	Filter FileFilter `json:"filter" mapstructure:"filter"`
	Sort   FileSort   `json:"sort" mapstructure:"sort"`
}

var sortFields = map[string]string{
	"index":         "publishedfileid",
	"name":          "title",
	"subscriptions": "subscriptions",
	"updated":       "time_updated",
	"created":       "time_created",
}

func GetPublishedFiles(db *sql.DB, appId uint, params FileParams) ([]WorkshopFileSummary, error) {
	if db == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	sortField, found := sortFields[params.Sort.Field]
	if !found {
		sortField = "time_created"
	}

	var sortDirection = "desc"
	if params.Sort.Ascending {
		sortDirection = "asc"
	}

	valueIndex := 2

	values := []any{appId}
	namePredicate := ""
	if params.Filter.Name != "" {
		// title
		namePredicate = " AND title ILIKE $" + strconv.Itoa(valueIndex)
		valueIndex++
		values = append(values, "%"+strings.ReplaceAll(strings.ReplaceAll(params.Filter.Name, "%", `\%`), "_", `\_`)+"%")
	}

	keys := []string{}
	for _, v := range params.Filter.Tags {
		keys = append(keys, " AND $"+strconv.Itoa(valueIndex)+"=ANY(tags)")
		values = append(values, v)
		valueIndex++

	}

	tagsPredicate := strings.Join(keys, "")

	types := []string{}
	for _, v := range params.Filter.Types {
		types = append(types, " OR $"+strconv.Itoa(valueIndex)+"=ANY(tags)")
		values = append(values, v)
		valueIndex++
	}
	typesPredicate := " AND (false " + strings.Join(types, "") + ") "

	query := `SELECT publishedfileid, title, preview_url, time_created, time_updated, subscriptions, array_to_json(tags)
	FROM published_files
	WHERE disabled = false AND consumer_appid = $1` + namePredicate + tagsPredicate + typesPredicate + ` ORDER BY ` + sortField + ` ` + sortDirection + ` LIMIT(50000);` // TODO: set limit as a param
	log.Println(query)
	res, err := db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in getPublishedFiles: <%w>", err)
	}
	defer res.Close()

	result := []WorkshopFileSummary{}

	for res.Next() {
		var publishedfileid uint64
		var title string
		var fileUrl string
		var timeCreated uint64
		var timeUpdated uint64
		var subscriptions uint64
		var tags string

		err = res.Scan(&publishedfileid, &title, &fileUrl, &timeCreated, &timeUpdated, &subscriptions, &tags)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in getPublishedFiles: <%w>", err)
		}
		var arr []string
		err = json.Unmarshal([]byte(tags), &arr)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: <%w>", err)
		}

		result = append(result, WorkshopFileSummary{publishedfileid, title, fileUrl, timeCreated, timeUpdated, subscriptions, arr})
	}

	// Check if we have an error
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in getPublishedFiles: <%w>", err)
	}

	return result, nil
}

func GetFileUrl(db *sql.DB, id uint64) (string, error) {
	row := db.QueryRow(`SELECT file_url FROM published_files WHERE publishedfileid = $1;`, id)

	file_url := ""

	err := row.Scan(&file_url)
	if err != nil {
		return "", fmt.Errorf("failed to scan row in GetFileUrl: <%w>", err)
	}

	return file_url, nil
}
