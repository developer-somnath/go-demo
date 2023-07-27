package helpers

import (
	"io"
	"mime/multipart"
	"os"
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

func UploadFileSingle(file *multipart.FileHeader, uploadLocation string) UploadResult {
	fileName := file.Filename

	// Create a new file on the server
	newFile, err := os.Create(uploadLocation + fileName)
	if err != nil {
		return UploadResult{FileName: "", Err: err}
	}
	defer newFile.Close()

	// Open the uploaded file for reading
	uploadedFile, err := file.Open()
	if err != nil {
		return UploadResult{FileName: "", Err: err}
	}
	defer uploadedFile.Close()

	// Copy the uploaded file data to the new file
	_, err = io.Copy(newFile, uploadedFile)
	if err != nil {
		return UploadResult{FileName: "", Err: err}
	}

	return UploadResult{FileName: fileName, Err: nil}
}
