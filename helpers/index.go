package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

// IntToString converts an int to a string.
func IntToString(value int) string {
	return strconv.Itoa(value)
}

// ValueOrDefault handles nil values for strings.
func ValueOrDefault(value *string) string {
	if value == nil {
		return "" // Set it to an empty string if it is nil
	}
	return *value // Otherwise, return the value of the pointer
}

type UploadResult struct {
	FileName string
	Err      error
}

func UploadFileSingle(fileHeader *multipart.FileHeader, file io.Reader, uploadLocation string) UploadResult {
	fileName := fileHeader.Filename

	// Create a new file on the server
	uploadsDirectory := filepath.Join("./public", uploadLocation)
	_, err := os.Stat(uploadsDirectory)
	if os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDirectory, 0777)
		if err != nil {
			return UploadResult{FileName: "", Err: err}
		}
	} else {
		newFile, err := os.Create(uploadsDirectory + "/" + fileName)
		if err != nil {
			return UploadResult{FileName: "", Err: err}
		}
		defer newFile.Close()

		// Copy the uploaded file data to the new file
		_, err = io.Copy(newFile, file)
		if err != nil {
			return UploadResult{FileName: "", Err: err}
		}
	}

	return UploadResult{FileName: fileName, Err: nil}
}
