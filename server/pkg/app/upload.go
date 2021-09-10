package app

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
)

func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
}

func uploadPhoto(file *multipart.File, filename, ext string) error {
	f, e := createFile(filename)
	if e != nil {
		return e
	}
	defer f.Close()
	var img image.Image

	if ext == "image/jpg" || ext == "image/jpeg" {
		if img, e = jpeg.Decode(*file); e == nil {
			e = jpeg.Encode(f, img, nil)
		}
	} else if ext == "image/png" {
		if img, e = png.Decode(*file); e == nil {
			e = png.Encode(f, img)
		}
	} else if ext != "image/gif" {
		if img, e = gif.Decode(*file); e == nil {
			e = gif.Encode(f, img, nil)
		}
	} else {
		os.Remove(filename)
		return errors.New("dont support this type of photo")
	}
	return e
}

func uploadFile(fileFormKey string, r *http.Request) (string, string, error) {
	file, fh, e := r.FormFile(fileFormKey)
	if e != nil || fh == nil {
		return "", "", errors.New("file did not found")
	}
	defer file.Close()

	if fh.Size/1024/1024 > 100 {
		return "", "", errors.New("this size is greater than 100mb")
	}

	fileExt := fh.Header.Get("Content-Type")
	filePreName := StringWithCharset(8)
	link := "/assets/img/" + filePreName + fh.Filename
	return link, fh.Filename, uploadPhoto(&file, link, fileExt)
}
