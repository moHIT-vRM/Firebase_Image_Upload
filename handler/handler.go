package handler

import (
	"GoImageUpload/database/dbHelper"
	"GoImageUpload/utils"
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ImageInfo struct {
	image *multipart.File
	name  string
}

func UploadImage(w http.ResponseWriter, r *http.Request) {

	imageDetails := ImageInfo{}
	//if parseErr := utils.ParseBody(r.Body, &imageDetails); parseErr != nil {
	//	utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
	//	return
	//}

	imageUrl, err := UploadImageFirebase(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to update")
		return
	}
	imageId, saveErr := dbHelper.SaveImageInfo(imageDetails.name, imageUrl)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "Failed to store information in db")
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		MessageToUser string `json:"messageToUser"`
		ImageId       string `json:"imageId"`
	}{
		MessageToUser: "Image Stored Successfully",
		ImageId:       imageId,
	})
}

type FirebaseApp struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}

func UploadImageFirebase(request *http.Request) (string, error) {
	client := FirebaseApp{}

	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	file, fileHeader, err := request.FormFile("image")
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	var multiPartMinValue, multiPartMaxValue int64
	multiPartMinValue = 10
	multiPartMaxValue = 20
	err = request.ParseMultipartForm(multiPartMinValue << multiPartMaxValue)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	defer func(file multipart.File) {
		fileErr := file.Close()
		if fileErr != nil {
			return
		}
	}(file)
	imagePath := "images/" + fileHeader.Filename
	bucket := "imageupload35.appspot.com"
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)

	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if err1 := bucketStorage.Close(); err1 != nil {
		logrus.Error(err1)
		return "", err
	}

	hours := 100

	signedURL := &cloud.SignedURLOptions{
		Scheme:  cloud.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(hours) * time.Hour),
	}

	url, err := client.Storage.Bucket(bucket).SignedURL(imagePath, signedURL)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return url, nil
}

// File Upload On the GoLang Folder:----->       //
//
//
// func UploadImage(w http.ResponseWriter, r *http.Request) {
//	//Limit 10mb
//	r.ParseMultipartForm(10 * 1024 * 1024)
//
//	file, handler, err := r.FormFile("yourFile")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer file.Close()
//
//	//fmt.Println("File Info")
//	fmt.Println("File Name:", handler.Filename)
//	//fmt.Println("File Size", handler.Size)
//	//fmt.Println("File type", handler.Header.Get("Content-Type"))
//
//	//Upload Image
//	tempFile, error := ioutil.TempFile("uploads", "uploads-*.jpg")
//
//	if error != nil {
//		fmt.Println(error)
//	}
//
//	defer tempFile.Close()
//
//	fileBytes, error2 := ioutil.ReadAll(file)
//
//	if error2 != nil {
//		fmt.Println(error2)
//	}
//	tempFile.Write(fileBytes)
//	fmt.Println("Done!!")
//
//}
