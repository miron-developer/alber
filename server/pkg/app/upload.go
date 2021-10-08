package app

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/transform"
)

func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
}

func resize(img image.Image) image.Image {
	return transform.Resize(img, 240, 240, transform.Linear)
}

func uploadFile(fileFormKey string, r *http.Request) (string, string, error) {
	file, fh, e := r.FormFile(fileFormKey)
	if e != nil || fh == nil {
		return "", "", errors.New("файл не найден")
	}

	if fh.Size/1024/1024 > 100 {
		return "", "", errors.New("файл больше 100мб")
	}

	fileExt := fh.Header.Get("Content-Type")
	if fileExt != "image/png" && fileExt != "image/jpeg" && fileExt != "image/jpg" && fileExt != "image/gif" {
		return "", "", errors.New("dont support this type of photo")
	}

	fileName := RandomStringFromCharsetAndLength("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
	ext := strings.Split(fh.Filename, ".")
	link := "/assets/img/" + fileName + "." + ext[len(ext)-1]

	wd, _ := os.Getwd()
	f, e := createFile(wd + link)
	if e != nil {
		return "", "", e
	}

	// resize
	var img image.Image
	if fileExt == "image/jpg" || fileExt == "image/jpeg" {
		img, e = jpeg.Decode(file)
		if e != nil {
			return "", "", e
		}
		e = jpeg.Encode(f, resize(img), nil)
	} else if fileExt == "image/png" {
		img, e = png.Decode(file)
		if e != nil {
			return "", "", e
		}
		e = png.Encode(f, resize(img))
	} else {
		img, e = gif.Decode(file)
		if e != nil {
			return "", "", e
		}
		e = gif.Encode(f, resize(img), nil)
	}
	if e != nil {
		return "", "", e
	}

	io.Copy(f, file)
	f.Close()
	file.Close()
	return link, fh.Filename, nil
}
