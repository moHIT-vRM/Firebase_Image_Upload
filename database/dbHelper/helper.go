package dbHelper

import "GoImageUpload/database"

func SaveImageInfo(imageName, imageUrl string) (string, error) {

	SQL := `INSERT INTO imageUpload(filename, url)
			VALUES ($1,$2)
			RETURNING id`
	var imageId string
	if err := database.ImageUploader.Get(&imageId, SQL, imageName, imageUrl); err != nil {
		return imageId, err
	}
	return imageId, nil
}
