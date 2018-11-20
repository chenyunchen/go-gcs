package storage

import (
	"encoding/hex"
	"encoding/json"
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
func SignURL(sp *service.Container, path string, fileName string, contentType string, method string, expires time.Time) (entity.SignedUrl, error) {
	opts := storage.SignedURLOptions{
		GoogleAccessID: sp.Config.GoogleCloud.ClientEmail,
		PrivateKey:     []byte(sp.Config.GoogleCloud.PrivateKey),
		Method:         method,
		Expires:        expires,
		ContentType:    contentType,
	}

	rand.Seed(time.Now().UnixNano())
	r := strings.NewReplacer("/", ":")
	b := make([]byte, 4) //equals 8 charachters
	rand.Read(b)
	hashFileName := fmt.Sprintf("%s_%s", hex.EncodeToString(b), r.Replace(fileName))

	su, err := storage.SignedURL(sp.Config.Storage.BucketName, path+hashFileName, &opts)
	fmt.Println(su)
	if err != nil {
		return entity.SignedUrl{}, err
	}
	u, err := url.Parse(su)
	if err != nil {
		return entity.SignedUrl{}, err
	}
	queries, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return entity.SignedUrl{}, err
	}

	return entity.SignedUrl{
		Url: strings.Split(su, "?")[0],
		UploadQueries: entity.UploadQueries{
			Expires:        queries["Expires"][0],
			GoogleAccessId: queries["GoogleAccessId"][0],
			Signature:      queries["Signature"][0],
		},
	}, nil
}

// CreateGCSGroupSignedUrl will sign group url by google cloud storage
func CreateGCSGroupSignedUrl(sp *service.Container, userId string, fileName string, contentType string, payload string) (entity.SignedUrl, error) {
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

// CreateGCSSingleSignedUrl will sign single url by google cloud storage
func CreateGCSSingleSignedUrl(sp *service.Container, userId string, fileName string, contentType string, payload string) (entity.SignedUrl, error) {
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)

	singlePayload := entity.SinglePayload{}
	json.Unmarshal([]byte(payload), &singlePayload)
	if err := sp.Validator.Struct(singlePayload); err != nil {
		return entity.SignedUrl{}, err
	}

	path := fmt.Sprintf("Single/%s/%s/", userId, singlePayload.To)

	return SignURL(sp, path, fileName, contentType, method, expires)
}
