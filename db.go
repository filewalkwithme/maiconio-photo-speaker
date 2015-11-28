package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var db *sql.DB

//get a list of all photos on postgres
func allPhotos() ([]Photo, error) {
	photos := []Photo{}

	//select all photos from db
	rows, err := db.Query(`SELECT id, file FROM photos order by id desc`)
	defer rows.Close()
	if err == nil {
		for rows.Next() {

			var id int
			var file string

			err = rows.Scan(&id, &file)
			if err == nil {
				idStr := strconv.Itoa(id)
				filename := "static/images/" + idStr + ".jpg"

				//check if file exists
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					//if not exists, then we need to retrieve the content from postgres...
					var content sql.NullString
					err := db.QueryRow(`SELECT content FROM photos where id = $1`, id).Scan(&content)
					if err != nil {
						fmt.Printf("error: %v\n", err)
						continue
					}
					contentStr := ""
					if content.Valid {
						contentStr = content.String
						data, err := base64.StdEncoding.DecodeString(contentStr)
						if err != nil {
							fmt.Printf("error: %v\n", err)
							continue
						}

						//and then save it to disk
						err = ioutil.WriteFile(filename, data, 0600)
						if err != nil {
							fmt.Printf("error: %v\n", err)
							continue
						}
						fmt.Printf("Saving: %v\n", filename)
					}
				}

				currentPhoto := Photo{ID: id, File: file}
				photos = append(photos, currentPhoto)
			} else {
				return photos, err
			}
		}
	}

	//now we need to retrive the photo labels
	for i, photo := range photos {
		labels := []Label{}
		rows, err := db.Query(`SELECT id, label, score FROM photo_labels where photo_id=$1`, photo.ID)
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				var idLabel int
				var nameLabel string
				var scoreLabel string

				err = rows.Scan(&idLabel, &nameLabel, &scoreLabel)
				if err == nil {
					currentLabel := Label{Name: nameLabel, Score: scoreLabel}
					labels = append(labels, currentLabel)
				}
			}
		}
		photos[i].Labels = labels
	}

	return photos, err
}

//insert a photo on postgres
func insertPhoto(file, content string, watsonResponse WatsonResponse) (int, error) {
	var photoID int
	err := db.QueryRow(`INSERT INTO photos(file, content) VALUES ($1, $2) RETURNING id`, file, content).Scan(&photoID)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Last inserted ID: %v\n", photoID)

	//also, insert the photo labels, based on Watson Visual Recognition response
	if photoID > 0 {
		if len(watsonResponse.Images) > 0 {
			for _, label := range watsonResponse.Images[0].Labels {
				var photoLabelID int
				err := db.QueryRow(`INSERT INTO photo_labels(photo_id, label, score) VALUES ($1, $2, $3) RETURNING id`, photoID, label.Name, label.Score).Scan(&photoLabelID)

				if err != nil {
					return 0, err
				}
			}
		}
	}
	return photoID, err
}

//remove a photo on postgres
func removePhoto(photoID int) (int, error) {
	//remove photo labels from postgres
	res, err := db.Exec(`delete from photo_labels where photo_id = $1`, photoID)
	if err != nil {
		return 0, err
	}

	//remove photo from postgres
	res, err = db.Exec(`delete from photos where id = $1`, photoID)
	if err != nil {
		return 0, err
	}

	//get number of rows deleted
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsDeleted), nil
}
