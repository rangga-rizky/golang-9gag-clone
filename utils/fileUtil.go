package utils

import (
	"io/ioutil"
	"mime/multipart"
)

func UploadFile(file multipart.File, handle *multipart.FileHeader, path string) (map[string]interface{}, bool) {
	allowedFormats := []string{"image/jpeg", "image/png"}

	mimeType := handle.Header.Get("Content-Type")

	for _, format := range allowedFormats {
		if format == mimeType {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				return Message(false, err.Error()), false
			}

			err = ioutil.WriteFile(path, data, 0666)
			if err != nil {
				return Message(false, err.Error()), false
			}
			return Message(true, "File Berhasil di upload"), true

		}
	}
	return Message(false, "Format file tidak diijinkan"), false

}
