package imageresize

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

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
	return nil, errors.New("invalid content type.")
}

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
	return errors.New("invalid content type.")
}

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

func ResizeImage(image image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, image, resize.Lanczos3)
}

// Thumbnail will fail if width or height == 0
// If maxWidth and maxHeight > image size, return the image with its original size
func ThumbnailImage(image image.Image, maxWidth, maxHeight uint) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, image, resize.Lanczos3)
}

func DeleteImage(path string) error {
	return os.Remove(path)
}
