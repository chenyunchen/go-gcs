package imageresize

import (
	"bufio"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// DownloadImageFromUrl will download image from url
func DownloadImageFromUrl(url, contentType string) (image.Image, error) {
	// Fetch an image.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch contentType {
	case "image/jpg":
		return jpeg.Decode(resp.Body)
	case "image/jpeg":
		return jpeg.Decode(resp.Body)
	case "jpeg":
		return jpeg.Decode(resp.Body)
	case "image/png":
		return png.Decode(resp.Body)
	}
	return nil, errors.New("invalid content type")
}

// ReadImageFile will read image struct from file path
func ReadImageFile(contentType, path string) (image.Image, error) {
	imgFp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFp.Close()

	switch contentType {
	case "image/jpg":
		return jpeg.Decode(bufio.NewReader(imgFp))
	case "image/jpeg":
		return jpeg.Decode(bufio.NewReader(imgFp))
	case "jpeg":
		return jpeg.Decode(bufio.NewReader(imgFp))
	case "image/png":
		return png.Decode(bufio.NewReader(imgFp))
	}
	return nil, errors.New("invalid content type")
}

// WriteImageFile will write image struct to file path
func WriteImageFile(image image.Image, contentType, path string) error {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	switch contentType {
	case "image/jpg":
		return WriteJpegImageFile(image, path, 90)
	case "image/jpeg":
		return WriteJpegImageFile(image, path, 90)
	case "jpeg":
		return WriteJpegImageFile(image, path, 90)
	case "image/png":
		return WritePngImageFile(image, path)
	}
	return errors.New("invalid content type")
}

// WritePngImageFile will write image struct to png file path
func WritePngImageFile(image image.Image, path string) error {
	imgFp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer imgFp.Close()

	err = png.Encode(imgFp, image)
	if err != nil {
		return err
	}
	return nil
}

// WriteJpegImageFile will write image struct to jpeg file path
func WriteJpegImageFile(image image.Image, path string, quality int) error {
	imgFp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer imgFp.Close()

	err = jpeg.Encode(imgFp, image, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}
	return nil
}

// ResizeImage will resize image struct and return image struct
func ResizeImage(image image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, image, resize.Lanczos3)
}

// ThumbnailImage will fail if width or height == 0
// If maxWidth and maxHeight > image size, return the image with its original size
func ThumbnailImage(image image.Image, maxWidth, maxHeight uint) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, image, resize.Lanczos3)
}

// DeleteImage will delete image from path
func DeleteImage(path string) error {
	return os.Remove(path)
}
