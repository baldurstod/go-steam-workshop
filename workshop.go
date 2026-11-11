package steamworkshop

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func RefreshPublishedFiles(ctx context.Context, db *sql.DB, apiKey string, appId uint, full bool, opts ...QueryFilesOption) error {
	log.Println("Starting items refresh")

	cursor := "*"
	page := 1

	for {
		log.Println("QueryFiles cursor "+cursor+" page", page)
		content, err := QueryFiles(apiKey, appId, cursor, opts...)
		if err != nil {
			return fmt.Errorf("error in RefreshPublishedFiles: <%w>", err)
		}

		response, ok := content["response"].(map[string]any)
		if !ok {
			return errors.New("missing response key in api response")
		}

		next, ok := response["next_cursor"].(string)
		if !ok {
			return errors.New("next_cursor is missing or don't have an expected format")
		}

		if next == cursor {
			// Pagination end
			break
		}

		files, ok := response["publishedfiledetails"].([]any)
		if !ok {
			return errors.New("publishedfiledetails don't have an expected format")
		}

		anyUpdated := false
		for _, f := range files {
			file, ok := f.(map[string]any)
			if !ok {
				return errors.New("publishedfiledetails don't have an expected format")
			}
			result, ok := file["result"].(float64)
			if !ok {
				return errors.New("publishedfiledetails don't have an expected format for key result")
			}
			if result == 0 {
				return errors.New("publishedfiledetails result is not 1")
			}
			if result != 1 {
				continue
			}

			updated, err := insertPublishedFile(db, file)
			if err != nil {
				log.Println("failed to insert item: ", err)
			}

			if updated {
				anyUpdated = true
			}
		}

		if !full && !anyUpdated {
			log.Println("no item updated on the entire page, leaving")
			break
		}

		page++
		cursor = next

		// Check context cancelation
		if ctx.Err() != nil {
			return nil
		}
	}

	return nil
}
