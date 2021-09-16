package app

import (
	"errors"
	"io"
	"net/http"
	"os"
)

func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
}

func uploadFile(fileFormKey string, r *http.Request) (string, string, error) {
	file, fh, e := r.FormFile(fileFormKey)
	if e != nil || fh == nil {
		return "", "", errors.New("file did not found")
	}

	if fh.Size/1024/1024 > 100 {
		return "", "", errors.New("this size is greater than 100mb")
	}

	filePreName := StringWithCharset(8)
	link := "/assets/img/" + filePreName + fh.Filename

	wd, _ := os.Getwd()
	f, e := createFile(wd + link)
	if e != nil {
		return "", "", e
	}

	io.Copy(f, file)
	f.Close()
	file.Close()
	return link, fh.Filename, nil
}
