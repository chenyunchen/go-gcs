package imageresize

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//Decode is image.Decode handling orientation in EXIF tags if exists.
//Requires io.ReadSeeker instead of io.Reader.
func Decode(reader io.ReadSeeker, contentType string) (image.Image, error) {
	var err error
	var img image.Image
	switch contentType {
	case "image/jpg":
		img, err = jpeg.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "jpeg":
		img, err = jpeg.Decode(reader)
	case "image/png":
		img, err = png.Decode(reader)
	default:
		err = errors.New("invalid content type")
	}
	if err != nil {
		return nil, err
	}

	reader.Seek(0, io.SeekStart)
	orientation := getOrientation(reader)
	switch orientation {
	case "1":
	case "2":
		img = imaging.FlipV(img)
	case "3":
		img = imaging.Rotate180(img)
	case "4":
		img = imaging.Rotate180(imaging.FlipV(img))
	case "5":
		img = imaging.Rotate270(imaging.FlipV(img))
	case "6":
		img = imaging.Rotate270(img)
	case "7":
		img = imaging.Rotate90(imaging.FlipV(img))
	case "8":
		img = imaging.Rotate90(img)
	}

	return img, err
}

func getOrientation(reader io.Reader) string {
	x, err := exif.Decode(reader)
	if err != nil {
		return "1"
	}
	if x != nil {
		orient, err := x.Get(exif.Orientation)
		if err != nil {
			return "1"
		}
		if orient != nil {
			return orient.String()
		}
	}

	return "1"
}

// DownloadImageFromUrl will download image from url
func DownloadImageFromUrl(url, contentType string) (image.Image, error) {
	// Fetch an image.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Decode(bytes.NewReader(body), contentType)
}

// ReadImageFile will read image struct from file path
func ReadImageFile(contentType, path string) (image.Image, error) {
	imgFp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFp.Close()

	return Decode(imgFp, contentType)
}

// WriteImageFile will write image struct to file path
func WriteImageFile(img image.Image, contentType, path string) error {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	switch contentType {
	case "image/jpg":
		return WriteJpegImageFile(img, path, 90)
	case "image/jpeg":
		return WriteJpegImageFile(img, path, 90)
	case "jpeg":
		return WriteJpegImageFile(img, path, 90)
	case "image/png":
		return WritePngImageFile(img, path)
	}
	return errors.New("invalid content type")
}

// WritePngImageFile will write image struct to png file path
func WritePngImageFile(img image.Image, path string) error {
	imgFp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer imgFp.Close()

	err = png.Encode(imgFp, img)
	if err != nil {
		return err
	}
	return nil
}

// WriteJpegImageFile will write image struct to jpeg file path
func WriteJpegImageFile(img image.Image, path string, quality int) error {
	imgFp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer imgFp.Close()

	err = jpeg.Encode(imgFp, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}
	return nil
}

// ResizeImage will resize image struct and return image struct
func ResizeImage(img image.Image, width, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// ThumbnailImage will fail if width or height == 0
// If maxWidth and maxHeight > image size, return the image with its original size
func ThumbnailImage(img image.Image, maxWidth, maxHeight int) *image.NRGBA {
	return imaging.Thumbnail(img, maxWidth, maxHeight, imaging.Lanczos)
}

// DeleteImage will delete image from path
func DeleteImage(path string) error {
	return os.Remove(path)
}
