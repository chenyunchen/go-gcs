package storageprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"go-gcs/src/imageresize"
	"go-gcs/src/logger"
	"go-gcs/src/service/googlecloud"
	"google.golang.org/api/option"
)

// GoogleCloudStoragePublicBaseUrl is public url for download
const GoogleCloudStoragePublicBaseUrl = "https://storage.googleapis.com"

// Storage is the structure for config
type Storage struct {
	Bucket             string `json:"bucket"`
	ImageResizeBucket  string `json:"imageResizeBucket"`
	ContentLengthRange string `json:"contentLengthRange"`
	AccessControl      string `json:"accessControl"`
}

// Service is the structure for service
type Service struct {
	Config  *Storage
	Client  *storage.Client
	Context context.Context
}

// New will reture a new service
func New(ctx context.Context, googleCloudConfig *googlecloud.Config, storageConfig *Storage) *Service {
	plan, err := json.Marshal(googleCloudConfig)
	if err != nil {
		logger.Warnf("error while read config file: %s", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(plan))
	if err != nil {
		logger.Warnf("error while create google cloud storage client: %s", err)
	}

	return &Service{
		Config:  storageConfig,
		Client:  client,
		Context: ctx,
	}
}

// Upload will call if need to upload to google cloud storage
func (s *Service) Upload(bucket, path, filePath string) error {
	r, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer r.Close()

	bh := s.Client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err := bh.Attrs(s.Context); err != nil {
		return err
	}

	obj := bh.Object(path)
	w := obj.NewWriter(s.Context)
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return obj.ACL().Set(s.Context, storage.AllUsers, storage.RoleReader)
}

// Delete will call if need to delete google cloud storage object
func (s *Service) Delete(bucket, path string) error {
	bh := s.Client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err := bh.Attrs(s.Context); err != nil {
		return err
	}

	obj := bh.Object(path)

	return obj.Delete(s.Context)
}

// UploadImage will upload image to google cloud storage
func (s *Service) UploadImage(img image.Image, contentType, bucket, path string) error {
	tmpPath := "/tmp/" + path
	if err := imageresize.WriteImageFile(img, contentType, tmpPath); err != nil {
		logger.Warnf("error while write image to local: %s", err)
		return err
	}

	if err := s.Upload(bucket, path, tmpPath); err != nil {
		logger.Warnf("error while upload image: %s", err)
		return err
	}

	return imageresize.DeleteImage(tmpPath)
}

// ResizeImageAndUpload will resize image and upload to google cloud storage
func (s *Service) ResizeImageAndUpload(img image.Image, width int, contentType, path string) error {
	img = imageresize.ResizeImage(img, width, 0)
	path = fmt.Sprintf("%s_%d", path, width)

	return s.UploadImage(img, contentType, s.Config.ImageResizeBucket, path)
}

// ResizeMultiImageSizeAndUpload will resize multiple image and upload to google cloud storage
func (s *Service) ResizeMultiImageSizeAndUpload(contentType, bucket, path string) (string, error) {
	// Check if already resize origin image
	imageResizeBucket := s.Config.ImageResizeBucket
	resizeUrl := fmt.Sprintf("%s/%s/%s", GoogleCloudStoragePublicBaseUrl, imageResizeBucket, path)
	_, err := imageresize.DownloadImageFromUrl(resizeUrl, contentType)
	if err == nil {
		return resizeUrl, nil
	}

	url := fmt.Sprintf("%s/%s/%s", GoogleCloudStoragePublicBaseUrl, bucket, path)
	img, err := imageresize.DownloadImageFromUrl(url, contentType)
	if err != nil {
		logger.Warnf("error while download image: %s", err)
		return "", err
	}

	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		s.ResizeImageAndUpload(img, 100, contentType, path)
	}()
	go func() {
		defer wg.Done()
		s.ResizeImageAndUpload(img, 150, contentType, path)
	}()
	go func() {
		defer wg.Done()
		s.ResizeImageAndUpload(img, 300, contentType, path)
	}()
	go func() {
		defer wg.Done()
		s.ResizeImageAndUpload(img, 640, contentType, path)
	}()
	go func() {
		defer wg.Done()
		s.ResizeImageAndUpload(img, 1080, contentType, path)
	}()
	wg.Wait()

	if err := s.UploadImage(img, contentType, imageResizeBucket, path); err != nil {
		return "", err
	}

	return resizeUrl, nil
}
