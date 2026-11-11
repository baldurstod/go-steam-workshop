package steamworkshop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-viper/mapstructure/v2"
	_ "github.com/lib/pq"
)

// Insert a published file in database
// Return true if the file was inserted or updated, false if it was up to date or an error occured
func insertPublishedFile(db *sql.DB, publishedFile map[string]any) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var workshopItem WorkshopItem
	err := mapstructure.Decode(publishedFile, &workshopItem)
	if err != nil {
		return false, fmt.Errorf("an error occured while decoding workshop item : <%w>", err)
	}

	var tags = []string{}

	for _, t := range workshopItem.Tags {
		tags = append(tags, t.Tag)
	}

	j, err := json.Marshal(&publishedFile)
	if err != nil {
		return false, fmt.Errorf("failed to marshal item: <%w>", err)
	}

	fileId, err := strconv.ParseUint(workshopItem.Publishedfileid, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert Publishedfileid "+workshopItem.Publishedfileid+": <%w>", err)
	}

	fileSize, err := strconv.ParseUint(workshopItem.FileSize, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert file size "+workshopItem.FileSize+": <%w> for item "+workshopItem.Publishedfileid, err)
	}

	row := db.QueryRow(`SELECT time_updated from published_files where publishedfileid = $1;`, fileId)

	var timeUpdated uint64

	err = row.Scan(&timeUpdated)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to scan row in insertItem: <%w>", err)
	}

	if err != sql.ErrNoRows && timeUpdated == workshopItem.TimeUpdated {
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
		workshopItem.Title,
		workshopItem.TimeCreated,
		workshopItem.TimeUpdated,
		workshopItem.Creator,
		tags,
		fileSize,
		workshopItem.FileUrl,
		workshopItem.PreviewUrl,
		workshopItem.Subscriptions,
		workshopItem.CreatorAppid,
		workshopItem.ConsumerAppid,
		j,
	)

	if err != nil {
		return false, fmt.Errorf("failed to insert workshop item "+workshopItem.Publishedfileid+" : <%w>", err)
	}

	return true, nil
}
