package storage

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"go-gcs/src/entity"
	"go-gcs/src/service"
)

// SignURL will sign url to google cloud storage
func SignURL(sp *service.Container, path string, fileName string, contentType string, method string, expires time.Time) (su entity.SignedUrl, err error) {
	opts := storage.SignedURLOptions{
		GoogleAccessID: sp.Config.GoogleCloud.ClientEmail,
		PrivateKey:     []byte(sp.Config.GoogleCloud.PrivateKey),
		Method:         method,
		Expires:        expires,
		ContentType:    contentType,
		Headers:        []string{"x-goog-content-length-range:" + sp.GoogleCloudStorage.Config.ContentLengthRange},
	}

	rand.Seed(time.Now().UnixNano())
	r := strings.NewReplacer("/", ":")
	b := make([]byte, 4) //equals 8 characters
	rand.Read(b)
	hashFileName := fmt.Sprintf("%s_%s", hex.EncodeToString(b), r.Replace(fileName))

	s, err := storage.SignedURL(sp.Config.Storage.Bucket, path+hashFileName, &opts)
	if err != nil {
		return su, err
	}

	u := strings.Split(s, "?")
	queries, err := url.ParseQuery(u[1])
	if err != nil {
		return su, err
	}

	return entity.SignedUrl{
		Url: u[0],
		UploadHeaders: entity.UploadHeaders{
			ContentType:        contentType,
			ContentLengthRange: sp.GoogleCloudStorage.Config.ContentLengthRange,
		},
		UploadQueries: entity.UploadQueries{
			Expires:        queries["Expires"][0],
			GoogleAccessId: queries["GoogleAccessId"][0],
			Signature:      queries["Signature"][0],
		},
	}, nil
}

// CreateGCSSingleSignedUrl will sign single url by google cloud storage
func CreateGCSSingleSignedUrl(sp *service.Container, userId, fileName, contentType, payload string) (su entity.SignedUrl, err error) {
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)
	singlePayload := entity.SinglePayload{}
	json.Unmarshal([]byte(payload), &singlePayload)
	if err := sp.Validator.Struct(singlePayload); err != nil {
		return su, err
	}

	path := fmt.Sprintf("Single/%s/%s/", userId, singlePayload.To)

	return SignURL(sp, path, fileName, contentType, method, expires)
}

// CreateGCSGroupSignedUrl will sign group url by google cloud storage
func CreateGCSGroupSignedUrl(sp *service.Container, userId, fileName, contentType, payload string) (entity.SignedUrl, error) {
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)

	groupPayload := entity.GroupPayload{}
	json.Unmarshal([]byte(payload), &groupPayload)
	if err := sp.Validator.Struct(groupPayload); err != nil {
		return entity.SignedUrl{}, err
	}

	path := fmt.Sprintf("Group/%s/%s/", groupPayload.GroupId, userId)

	return SignURL(sp, path, fileName, contentType, method, expires)
}

// ResizeGCSImage will resize image from google cloud storage
func ResizeGCSImage(sp *service.Container, url, contentType string) (ri entity.ResizeImage, err error) {
	if contentType == "image/jpg" || contentType == "image/jpeg" || contentType == "jpeg" || contentType == "image/png" {
		u := strings.Split(url, "/")
		path := strings.Join(u[4:], "/")
		url, err = sp.GoogleCloudStorage.ResizeMultiImageSizeAndUpload(contentType, sp.GoogleCloudStorage.Config.Bucket, path)
		if err != nil {
			return ri, err
		}

		return entity.ResizeImage{
			Origin:         url,
			ThumbWidth100:  url + "_100",
			ThumbWidth150:  url + "_150",
			ThumbWidth300:  url + "_300",
			ThumbWidth640:  url + "_640",
			ThumbWidth1080: url + "_1080",
		}, nil
	}

	return ri, errors.New("invalid content type.")
}
